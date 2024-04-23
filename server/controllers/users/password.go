package accounts

import (
	"context"
	"darco/proto/db"
	users "darco/proto/models/people"
	"darco/proto/resolvers"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
)

type ValidatePasswordTokenInput struct {
	Token string `path:"token"`
}

func ValidatePasswordToken(ctx context.Context, input *ValidatePasswordTokenInput) (*struct{}, error) {
	_, tokenValid := users.ValidateAccountToken(
		db.Client(),
		users.Token(input.Token),
		users.PasswordResetToken,
	)
	if !tokenValid {
		return nil, huma.Error400BadRequest("Invalid token")
	}
	return nil, nil
}

type UpdatePasswordInput struct {
	resolvers.AuthRequired
	Body users.UpdatePasswordInput `required:"true"`
}

func (i *UpdatePasswordInput) Resolve(ctx huma.Context) []error {
	if errs := i.AuthRequired.Resolve(ctx); errs != nil {
		return errs
	}
	if !i.User.PasswordMatch(i.DB(), i.Body.Password) {
		return []error{&huma.ErrorDetail{
			Message:  "Invalid password",
			Location: "password",
			Value:    i.Body.Password,
		}}
	}
	return nil
}

func UpdatePassword(ctx context.Context, input *UpdatePasswordInput) (*struct{}, error) {
	if err := input.User.SetPassword(
		input.DB(),
		input.Body.NewPassword.Password,
	); err != nil {
		return nil, huma.Error400BadRequest("New password rejected",
			&huma.ErrorDetail{
				Message:  "Password is too weak",
				Location: "new_password.password",
			})
	}
	return nil, nil
}

type PasswordResetInput struct {
	Token users.Token         `path:"token"`
	Body  users.PasswordInput `required:"true"`
	*users.User
}

func (i *PasswordResetInput) Resolve(ctx huma.Context) []error {
	user, tokenValid := users.ValidateAccountToken(
		db.Client(),
		i.Token,
		users.PasswordResetToken)
	if !tokenValid {
		return []error{huma.Error400BadRequest("Invalid token")}
	}
	i.User = user
	return nil
}

func PasswordReset(ctx context.Context, input *PasswordResetInput) (*struct{}, error) {
	if err := input.User.SetPassword(db.Client(), input.Body.Password); err != nil {
		return nil, &huma.ErrorDetail{
			Message:  "Password is too weak",
			Location: "password",
		}
	}
	return nil, nil
}

type RequestPasswordResetInput struct {
	resolvers.HostResolver
	Body struct {
		Email   string   `json:"email" binding:"required,email" format:"email"`
		Handler *url.URL `json:"handler,omitempty" doc:"A URL where a form to set the new password is available"`
	}
}

type RequestPasswordResetHandler func(context.Context, *RequestPasswordResetInput) (*struct{}, error)

func RequestPasswordReset(defaultHandlerPath string) RequestPasswordResetHandler {
	return func(ctx context.Context, input *RequestPasswordResetInput) (*struct{}, error) {
		user, err := users.Find(db.Client(), input.Body.Email)
		if err != nil {
			return nil, huma.Error400BadRequest("Unknown email address", err)
		}

		var targetURL = input.URL
		targetURL.Path = defaultHandlerPath
		if input.Body.Handler != nil {
			targetURL = *input.Body.Handler
		}

		if err = user.RequestPasswordReset(db.Client(), targetURL); err != nil {
			return nil, huma.Error500InternalServerError(
				"Failed to issue password reset token",
				err,
			)
		}
		return nil, nil
	}
}
