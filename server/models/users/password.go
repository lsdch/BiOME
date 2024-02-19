package users

import (
	"darco/proto/models/validations"
	"fmt"

	"github.com/trustelem/zxcvbn"
)

type PasswordInput struct {
	Password   string `json:"password" binding:"required"`
	ConfirmPwd string `json:"password_confirmation" binding:"eqfield=Password,required"`
} //@name PasswordInput

// Used to compute password strength with regard to personal infos
func (user_info *InnerUserInput) PasswordSensitiveValues() []string {
	return []string{user_info.Email, user_info.Login, user_info.Person.FirstName, user_info.Person.LastName}
}

// Validates that a password has a high enough entropy.
// Strength score ranges from 1 to 5.
func (pwd *PasswordInput) ValidateStrength(strength int, user_info *InnerUserInput) error {
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

type NewPasswordInput struct {
	Password    string        `json:"password" binding:"required"`
	NewPassword PasswordInput `json:"new_password" binding:"required"`
} // @name NewPasswordInput
