package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	MIN_PASSWORD_STRENGTH = 3
)

type AccountService struct {
	authService *AuthService
	config      config.BootstrapConfig
}

func NewAccountService(authService *AuthService, config config.BootstrapConfig) *AccountService {
	return &AccountService{authService: authService, config: config}
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

func (s *AccountService) UpdatePassword(ctx context.Context, q db.Querier, loginOrEmail string, params models.PasswordUpdateParams) error {
	user, err := s.authService.AuthenticateCredentials(ctx, q, models.UserCredentials{Identifier: loginOrEmail, Password: params.OldPassword})
	if err != nil {
		return err
	}
	userModel := models.UserFromDB(user)
	if !userModel.ValidatePasswordStrength(params.NewPassword, MIN_PASSWORD_STRENGTH) {
		return fmt.Errorf("new password is too weak")
	}
	err = q.Queries().UpdateUserPassword(ctx, params.NewPassword, user.ID)
	if err != nil {
		return err
	}
	return nil
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
