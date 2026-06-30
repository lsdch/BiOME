package models

import (
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/trustelem/zxcvbn"
)

type UserRole = biomedb.UserRole

type User struct {
	ID           uuid.UUID        `json:"id"`
	Email        string           `json:"email"`
	Login        string           `json:"login"`
	Role         UserRole         `json:"role"`
	FirstName    string           `json:"first_name"`
	LastName     string           `json:"last_name"`
	FullName     string           `json:"full_name"`
	Organisation Optional[string] `json:"organization,omitempty"`
	Contact      Optional[string] `json:"contact,omitempty"`
	Bio          Optional[string] `json:"bio,omitempty"`
	Active       bool             `json:"active"`
}

func UserFromDB(u biomedb.User) User {
	return User{
		ID:           u.ID,
		Email:        u.Email,
		Login:        u.Login,
		Role:         u.Role,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		FullName:     u.FullName,
		Organisation: NewOptionalFromPtr(u.Organisation),
		Contact:      NewOptionalFromPtr(u.Contact),
		Bio:          NewOptionalFromPtr(u.Bio),
		Active:       u.Active,
	}
}

func (u *User) IsGranted(role biomedb.UserRole) bool {
	return u.Role.IsGreaterEqual(role)
}

func (u *User) PasswordSensitiveInfos() []string {
	return []string{u.Email, u.Login, u.FirstName, u.LastName}
}

func (u *User) ValidatePasswordStrength(pwd string, strength int) bool {
	score := zxcvbn.PasswordStrength(pwd, u.PasswordSensitiveInfos()).Score
	return score > strength
}

type UserInfoUpdateParams struct {
	FirstName    Optional[string] `json:"first_name"`
	LastName     Optional[string] `json:"last_name"`
	Contact      Optional[string] `json:"contact"`
	Organisation Optional[string] `json:"organisation"`
}

func (p *UserInfoUpdateParams) ToParams(userID uuid.UUID) biomedb.UpdateUserPersonalInfoParams {
	return biomedb.UpdateUserPersonalInfoParams{
		UserID:          userID,
		FirstName:       p.FirstName.ToPtr(),
		LastName:        p.LastName.ToPtr(),
		FirstNameSet:    p.FirstName.IsSet,
		LastNameSet:     p.LastName.IsSet,
		ContactSet:      p.Contact.IsSet,
		Contact:         p.Contact.ToPtr(),
		OrganisationSet: p.Organisation.IsSet,
		Organisation:    p.Organisation.ToPtr(),
	}
}

type PasswordUpdateParams struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
