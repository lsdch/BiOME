package models

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type PasswordResetParams struct {
	NewPassword string `validate:"required"`
}

// Validate uses account details loaded from the database as password context.
func (p PasswordResetParams) Validate(user User, minimumStrength int) error {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterStructValidation(func(sl validator.StructLevel) {
		if len(p.NewPassword) > 72 {
			sl.ReportError(p.NewPassword, "NewPassword", "NewPassword", "maxbytes", "72")
			return
		}
		if p.NewPassword != "" && !user.ValidatePasswordStrength(p.NewPassword, minimumStrength) {
			sl.ReportError(p.NewPassword, "NewPassword", "NewPassword", "password_strength", fmt.Sprint(minimumStrength))
		}
	}, PasswordResetParams{})
	return v.Struct(p)
}
