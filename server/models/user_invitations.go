package models

import (
	"time"

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
	Login    string
	Password string

	FirstName    string
	LastName     string
	Organisation Optional[string]
	Contact      Optional[string]
	Bio          Optional[string]
}

func (p *RegisterFromInvitationParams) ToParams(tokenHash string) biomedb.CreateUserFromInvitationTokenParams {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	return biomedb.CreateUserFromInvitationTokenParams{
		TokenHash:    tokenHash,
		Login:        p.Login,
		PasswordHash: string(passwordHash),
		FirstName:    p.FirstName,
		LastName:     p.LastName,
		Organisation: p.Organisation.ToPtr(),
		Contact:      p.Contact.ToPtr(),
		Bio:          p.Bio.ToPtr(),
	}
}
