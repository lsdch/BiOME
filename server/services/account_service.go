package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
)

const (
	MIN_PASSWORD_STRENGTH = 3
)

type AccountService struct {
	db          *db.DB
	authService *AuthService
}

func NewAccountService(db *db.DB, authService *AuthService) *AccountService {
	return &AccountService{db: db, authService: authService}
}

func (s *AccountService) UpdatePassword(ctx context.Context, loginOrEmail string, params models.PasswordUpdateParams) error {
	user, err := s.authService.AuthenticateCredentials(ctx, models.UserCredentials{Identifier: loginOrEmail, Password: params.OldPassword})
	if err != nil {
		return err
	}
	userModel := models.UserFromDB(user)
	if !userModel.ValidatePasswordStrength(params.NewPassword, MIN_PASSWORD_STRENGTH) {
		return fmt.Errorf("new password is too weak")
	}
	err = s.db.Queries().UpdateUserPassword(ctx, params.NewPassword, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *AccountService) SetRole(ctx context.Context, userID uuid.UUID, role models.UserRole) error {
	err := s.db.Queries().UpdateUserRole(ctx, role, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *AccountService) UpdateUserInfo(ctx context.Context, userID uuid.UUID, params models.UserInfoUpdateParams) error {
	_, err := s.db.Queries().UpdateUserPersonalInfo(ctx, params.ToParams(userID))
	return err
}
