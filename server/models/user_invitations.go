package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"golang.org/x/crypto/bcrypt"
)

type CreateInvitationParams struct {
	Email       string
	InviteeName string
	Role        UserRole
	Message     Optional[string]
	InviterID   uuid.UUID
	ExpiresAt   time.Time
}

type InvitationResult struct {
	Invitation biomedb.Invitation
	Token      string // token brut à envoyer par email
}

type RegisterFromInvitationParams struct {
	Login    string `validate:"required,min=2"`
	Password string `validate:"required"`

	FirstName    string `validate:"required,min=2"`
	LastName     string `validate:"required,min=2"`
	Organisation Optional[string]
	Contact      Optional[string]
	Bio          Optional[string]
}

// Validate checks registration fields without modifying their stored values.
// Email comes from the invitation, so password strength includes trusted user context.
func (p RegisterFromInvitationParams) Validate(email string, minimumPasswordStrength int) error {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterStructValidation(func(sl validator.StructLevel) {
		input := sl.Current().Interface().(RegisterFromInvitationParams)
		if len(input.Password) > 72 {
			sl.ReportError(input.Password, "Password", "Password", "maxbytes", "72")
			return
		}
		user := User{Login: p.Login, Email: email, FirstName: p.FirstName, LastName: p.LastName}
		if input.Password != "" && !user.ValidatePasswordStrength(input.Password, minimumPasswordStrength) {
			sl.ReportError(input.Password, "Password", "Password", "password_strength", fmt.Sprint(minimumPasswordStrength))
		}
	}, RegisterFromInvitationParams{})

	// min counts Unicode characters; whitespace must not satisfy the name/login minimum.
	input := p
	input.Login = strings.TrimSpace(input.Login)
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	return v.Struct(input)
}

func (p *RegisterFromInvitationParams) ToParams(tokenHash string) (biomedb.CreateUserFromInvitationTokenParams, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return biomedb.CreateUserFromInvitationTokenParams{}, fmt.Errorf("hash registration password: %w", err)
	}
	return biomedb.CreateUserFromInvitationTokenParams{
		TokenHash:    tokenHash,
		Login:        p.Login,
		PasswordHash: string(passwordHash),
		FirstName:    p.FirstName,
		LastName:     p.LastName,
		Organisation: p.Organisation.ToPtr(),
		Contact:      p.Contact.ToPtr(),
		Bio:          p.Bio.ToPtr(),
	}, nil
}
