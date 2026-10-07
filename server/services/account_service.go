package services

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	email_templates "github.com/lsdch/biome/templates"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	MIN_PASSWORD_STRENGTH = 3
	passwordResetLifetime = time.Hour
)

var (
	ErrInvalidCurrentPassword    = errors.New("invalid current password or inactive account")
	ErrInvalidPasswordUpdate     = errors.New("invalid password update parameters")
	ErrInvalidPasswordResetToken = errors.New("password reset token is invalid, expired or already consumed")
	ErrInvalidPasswordReset      = errors.New("invalid password reset parameters")
)

type AccountService struct {
	authService  *AuthService
	config       config.BootstrapConfig
	emailService interface {
		Available() bool
		Send(context.Context, string, string, templ.Component) error
	}
	publicBaseURL url.URL
}

func NewAccountService(authService *AuthService, cfg config.Config, emailService *EmailService) *AccountService {
	s := &AccountService{authService: authService, config: cfg.Bootstrap, publicBaseURL: cfg.AppPublicBaseURL}
	if emailService != nil {
		s.emailService = emailService
	}
	return s
}

// RequestPasswordReset persists a hashed token before sending its link by email.
// It requires a database rather than a caller-owned transaction so the token is
// committed before email delivery. Unknown and inactive accounts are ignored.
// A delivery failure leaves the token valid; the caller may submit a new request.
func (s *AccountService) RequestPasswordReset(ctx context.Context, database *db.DB, email, clientPath string) error {
	if s.emailService == nil || !s.emailService.Available() {
		return ErrMailerUnavailable
	}
	base := s.publicBaseURL
	if (base.Scheme != "https" && base.Scheme != "http") || base.Host == "" || strings.TrimSpace(clientPath) == "" {
		return fmt.Errorf("invalid password reset URL configuration")
	}
	user, err := database.Queries().GetUserByLoginOrEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("get password reset user: %w", err)
	}
	if !user.Active {
		return nil
	}
	token, err := GenerateSecureToken()
	if err != nil {
		return fmt.Errorf("generate password reset token: %w", err)
	}
	if err := database.Queries().CreatePasswordResetToken(ctx, biomedb.CreatePasswordResetTokenParams{
		UserID: user.ID, TokenHash: HashToken(token), ExpiresAt: time.Now().Add(passwordResetLifetime),
	}); err != nil {
		return fmt.Errorf("create password reset token: %w", err)
	}
	link := base.JoinPath(clientPath)
	query := link.Query()
	query.Set("token", token)
	link.RawQuery = query.Encode()
	data := email_templates.PasswordResetData{Name: user.FullName, URL: *link}
	if err := s.emailService.Send(ctx, user.Email, data.Subject(), email_templates.PasswordReset(data)); err != nil {
		return fmt.Errorf("send password reset email: %w", err)
	}
	return nil
}

// VerifyPasswordResetToken checks current validity without consuming or extending
// the token. ResetPassword must still validate it when the password is submitted.
func (s *AccountService) VerifyPasswordResetToken(ctx context.Context, q db.Querier, token string) error {
	if token == "" {
		return ErrInvalidPasswordResetToken
	}
	valid, err := q.Queries().IsPasswordResetTokenValid(ctx, HashToken(token))
	if err != nil {
		return fmt.Errorf("verify password reset token: %w", err)
	}
	if !valid {
		return ErrInvalidPasswordResetToken
	}
	return nil
}

// ResetPassword atomically changes the password, invalidates outstanding reset
// tokens and revokes existing sessions. Failed validation leaves the token usable.
func (s *AccountService) ResetPassword(ctx context.Context, database *db.DB, token string, params models.PasswordResetParams) error {
	if token == "" {
		return ErrInvalidPasswordResetToken
	}
	return database.WithTx(ctx, func(tx *db.Tx) error {
		tokenHash := HashToken(token)
		user, err := tx.Queries().GetPasswordResetUserForUpdate(ctx, tokenHash)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvalidPasswordResetToken
		}
		if err != nil {
			return fmt.Errorf("get password reset user: %w", err)
		}
		if err := params.Validate(models.UserFromDB(user), MIN_PASSWORD_STRENGTH); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidPasswordReset, err)
		}
		passwordHash, err := s.HashPassword(params.NewPassword)
		if err != nil {
			return fmt.Errorf("hash reset password: %w", err)
		}
		// Recheck after obtaining the account lock and hashing: another reset
		// may have consumed this token, or it may have expired in the meantime.
		if _, err := tx.Queries().ConsumePasswordResetToken(ctx, tokenHash); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrInvalidPasswordResetToken
			}
			return fmt.Errorf("consume password reset token: %w", err)
		}
		if err := tx.Queries().UpdateUserPassword(ctx, string(passwordHash), user.ID); err != nil {
			return fmt.Errorf("update reset password: %w", err)
		}
		if err := tx.Queries().InvalidatePasswordResetTokensForUser(ctx, user.ID); err != nil {
			return fmt.Errorf("invalidate password reset tokens: %w", err)
		}
		if err := tx.Queries().RevokeAllSessionsForUser(ctx, user.ID); err != nil {
			return fmt.Errorf("revoke password reset sessions: %w", err)
		}
		return nil
	})
}

func (s *AccountService) ListUsers(ctx context.Context, q db.Querier) ([]models.User, error) {
	users, err := q.Queries().ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	userModels := make([]models.User, len(users))
	for i, u := range users {
		userModels[i] = models.UserFromDB(u)
	}
	return userModels, nil
}

// UpdatePassword atomically changes the authenticated user's password and revokes
// existing sessions and reset tokens. userID must come from the trusted session.
func (s *AccountService) UpdatePassword(ctx context.Context, database *db.DB, userID uuid.UUID, params models.PasswordUpdateParams) error {
	if userID == uuid.Nil || params.OldPassword == "" {
		return ErrInvalidCurrentPassword
	}
	return database.WithTx(ctx, func(tx *db.Tx) error {
		user, err := tx.Queries().GetActiveUserForPasswordUpdate(ctx, userID)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvalidCurrentPassword
		}
		if err != nil {
			return fmt.Errorf("get password update user: %w", err)
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.OldPassword)); err != nil {
			return ErrInvalidCurrentPassword
		}
		if err := (models.PasswordResetParams{NewPassword: params.NewPassword}).Validate(models.UserFromDB(user), MIN_PASSWORD_STRENGTH); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidPasswordUpdate, err)
		}
		hash, err := s.HashPassword(params.NewPassword)
		if err != nil {
			return fmt.Errorf("hash updated password: %w", err)
		}
		if err := tx.Queries().UpdateUserPassword(ctx, string(hash), user.ID); err != nil {
			return fmt.Errorf("update password: %w", err)
		}
		if err := tx.Queries().InvalidatePasswordResetTokensForUser(ctx, user.ID); err != nil {
			return fmt.Errorf("invalidate reset tokens: %w", err)
		}
		if err := tx.Queries().RevokeAllSessionsForUser(ctx, user.ID); err != nil {
			return fmt.Errorf("revoke sessions: %w", err)
		}
		return nil
	})
}

func (s *AccountService) SetRole(ctx context.Context, q db.Querier, userID uuid.UUID, role models.UserRole) error {
	err := q.Queries().UpdateUserRole(ctx, role, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *AccountService) UpdateUserInfo(ctx context.Context, q db.Querier, userID uuid.UUID, params models.UserInfoUpdateParams) error {
	_, err := q.Queries().UpdateUserPersonalInfo(ctx, params.ToParams(userID))
	return err
}

func (s *AccountService) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (s *AccountService) BootstrapUsers(ctx context.Context, q db.Querier) error {
	for _, user := range s.config.Users {
		_, err := q.Queries().GetUserByLoginOrEmail(ctx, user.Login)
		if err == nil {
			logrus.Infof("Bootstrap user %s already exists, skipping", user.Login)
			continue // User already exists, skip
		}
		pwdHash, err := s.HashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password for user %s: %w", user.Login, err)
		}
		logrus.Infof("Creating bootstrap user %s with role %s", user.Login, user.Role)
		if _, err := q.Queries().CreateUser(ctx, biomedb.CreateUserParams{
			Login:        user.Login,
			Email:        user.Email,
			PasswordHash: string(pwdHash),
			Role:         user.Role,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
		}); err != nil {
			return fmt.Errorf("failed to create bootstrap user %s: %w", user.Login, err)
		}
	}
	return nil
}
