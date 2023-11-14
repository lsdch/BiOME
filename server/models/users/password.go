package users

import (
	"darco/proto/models/validations"
	"fmt"

	"github.com/trustelem/zxcvbn"
	"golang.org/x/crypto/bcrypt"
)

type PasswordInput struct {
	Password   string `json:"password" binding:"required"`
	ConfirmPwd string `json:"password_confirmation" binding:"eqfield=Password,required"`
} //@name PasswordInput

func (pwd *PasswordInput) Hash() (string, error) {
	return hashPassword(pwd.Password)
}

func (user_info *InnerUserInput) PasswordSensitiveValues() []string {
	return []string{user_info.Email, user_info.Login, user_info.Person.FirstName, user_info.Person.LastName}
}

// Validates that a password has a high enough entropy.
// Strength score ranges from 1 to 5.
func (pwd *PasswordInput) ValidateStrength(strength int, user_info *InnerUserInput) *validations.InputValidationError {
	score := zxcvbn.PasswordStrength(pwd.Password, user_info.PasswordSensitiveValues()).Score
	if score < strength {
		return &validations.InputValidationError{
			Field:     "password",
			Message:   "Password is too weak",
			ErrString: fmt.Sprintf("Field password has strength score %d, at least %d is required", score, strength),
		}
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// Checks that a password matches a hash from the database
func VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}

type NewPasswordInput struct {
	Password    string        `json:"password" binding:"required"`
	NewPassword PasswordInput `json:"new_password" binding:"required"`
} // @name NewPasswordInput
