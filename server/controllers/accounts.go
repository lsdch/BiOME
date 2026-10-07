package controllers

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/app_errors"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
	"github.com/sirupsen/logrus"
)

type AccountsController struct {
	db      *db.DB
	service interface {
		UpdatePassword(context.Context, *db.DB, uuid.UUID, models.PasswordUpdateParams) error
		ListUsers(context.Context, db.Querier) ([]models.User, error)
		RequestPasswordReset(context.Context, *db.DB, string, string) error
		ResetPassword(context.Context, *db.DB, string, models.PasswordResetParams) error
		VerifyPasswordResetToken(context.Context, db.Querier, string) error
	}
}

func NewAccountsController(db *db.DB, service *services.AccountService) *AccountsController {
	return &AccountsController{db: db, service: service}
}

func (c *AccountsController) ListUsers(ctx context.Context, _ *struct{}) (*BodyTransporter[[]models.User], error) {
	users, err := c.service.ListUsers(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.User]{Body: users}, nil
}

type PasswordResetRequestInput struct {
	Body struct {
		Email      string `json:"email" minLength:"1" format:"email" doc:"Account email address"`
		ClientPath string `json:"client_path" minLength:"1" doc:"Client-side path to the password reset page, appended to the public base URL"`
	}
}

type PasswordResetInput struct {
	Body struct {
		Token       string `json:"token" minLength:"1" doc:"Password reset token received by email"`
		NewPassword string `json:"new_password" minLength:"1" maxLength:"72" writeOnly:"true" doc:"New password (maximum 72 bytes)"`
	}
}

type PasswordResetVerificationInput struct {
	Body struct {
		Token string `json:"token" minLength:"1" doc:"Password reset token received by email"`
	}
}

func (c *AccountsController) VerifyPasswordResetToken(ctx context.Context, input *PasswordResetVerificationInput) (*struct{}, error) {
	err := c.service.VerifyPasswordResetToken(ctx, c.db, input.Body.Token)
	return nil, passwordResetAPIError(err)
}

func (c *AccountsController) RequestPasswordReset(ctx context.Context, input *PasswordResetRequestInput) (*struct{}, error) {
	err := c.service.RequestPasswordReset(ctx, c.db, input.Body.Email, input.Body.ClientPath)
	return nil, passwordResetAPIError(err)
}

func (c *AccountsController) ResetPassword(ctx context.Context, input *PasswordResetInput) (*struct{}, error) {
	err := c.service.ResetPassword(ctx, c.db, input.Body.Token, models.PasswordResetParams{NewPassword: input.Body.NewPassword})
	return nil, passwordResetAPIError(err)
}

// The router recognizes AppErrorProvider; plain Huma errors would become 500s.
type passwordResetHTTPError struct {
	status int
	detail string
}

func (e passwordResetHTTPError) Error() string { return e.detail }

func (e passwordResetHTTPError) AppError() *app_errors.AppError {
	return &app_errors.AppError{ErrorModel: huma.ErrorModel{
		Status: e.status, Title: http.StatusText(e.status), Detail: e.detail,
	}}
}

func passwordResetAPIError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, services.ErrInvalidPasswordResetToken):
		return passwordResetHTTPError{http.StatusBadRequest, "Password reset token is invalid, expired or already consumed"}
	case errors.Is(err, services.ErrInvalidPasswordReset):
		return passwordResetHTTPError{http.StatusUnprocessableEntity, "New password does not meet password requirements"}
	case errors.Is(err, services.ErrMailerUnavailable):
		return services.ErrMailerUnavailable
	default:
		logrus.Error(err)
		return passwordResetHTTPError{http.StatusInternalServerError, "Unable to process password reset"}
	}
}

func (c *AccountsController) UpdatePassword(ctx context.Context, input *BodyTransporter[models.PasswordUpdateParams]) (*struct{}, error) {
	session, ok := auth.SessionFromContext(ctx)
	if !ok || session == nil || session.User.ID == uuid.Nil {
		return nil, passwordResetHTTPError{http.StatusUnauthorized, "Authentication required"}
	}
	err := c.service.UpdatePassword(ctx, c.db, session.User.ID, input.Body)
	switch {
	case err == nil:
		return nil, nil
	case errors.Is(err, services.ErrInvalidCurrentPassword):
		return nil, passwordResetHTTPError{http.StatusForbidden, "Current password is incorrect or account is inactive"}
	case errors.Is(err, services.ErrInvalidPasswordUpdate):
		return nil, passwordResetHTTPError{http.StatusUnprocessableEntity, "New password does not meet password requirements"}
	default:
		return nil, passwordResetHTTPError{http.StatusInternalServerError, "Unable to update password"}
	}
}

func (c *AccountsController) RegisterRoutes(r *router.Router) {
	accountsAPI := r.RouteGroup("/accounts").WithTags([]string{"Accounts"})

	router.NewSpec(accountsAPI, "UpdatePassword", huma.Operation{
		Method:        http.MethodPut,
		Path:          "/password",
		Summary:       "Change the current user's password",
		Description:   "Requires the current password. Revokes all sessions, including the current session, and invalidates reset tokens. Sign in again after success.",
		DefaultStatus: http.StatusNoContent,
		Errors:        []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusUnprocessableEntity},
	}, c.UpdatePassword).WithAccessPolicy(auth.Authenticated()).Register(r)

	router.NewSpec(accountsAPI, "VerifyPasswordResetToken", huma.Operation{
		Method:        http.MethodPost,
		Path:          "/password-reset/verify",
		Summary:       "Verify a password reset token",
		Description:   "Checks token validity and account activity without consuming the token or extending its expiration. Password submission still revalidates the token.",
		DefaultStatus: http.StatusNoContent,
		Errors:        []int{http.StatusBadRequest, http.StatusUnprocessableEntity},
	}, c.VerifyPasswordResetToken).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(accountsAPI, "RequestPasswordReset", huma.Operation{
		Method:        http.MethodPost,
		Path:          "/password-reset/request",
		Summary:       "Request a password reset email",
		Description:   "Returns the same empty response for existing, unknown and inactive accounts.",
		DefaultStatus: http.StatusNoContent,
		Errors:        []int{http.StatusUnprocessableEntity, http.StatusServiceUnavailable},
	}, c.RequestPasswordReset).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(accountsAPI, "ResetPassword", huma.Operation{
		Method:        http.MethodPost,
		Path:          "/password-reset",
		Summary:       "Set a new password using a reset token",
		Description:   "Consumes the reset token, invalidates other reset tokens and revokes existing sessions.",
		DefaultStatus: http.StatusNoContent,
		Errors:        []int{http.StatusBadRequest, http.StatusUnprocessableEntity},
	}, c.ResetPassword).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		accountsAPI,
		"ListUsers",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/",
			Summary: "List all users",
		},
		c.ListUsers,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleAdmin)).Register(r)
}
