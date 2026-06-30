package models

import (
	"net/netip"

	"github.com/google/uuid"
)

type AuthContext struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
}

type SessionContext struct {
	AuthContext
	User User
}

type UserCredentials struct {
	Identifier string `json:"identifier" binding:"required"`
	// Unhashed password, will be hashed before checking against the database
	Password string `json:"password" binding:"required"`
}

type AuthOrigin struct {
	UserAgent Optional[string]
	IP        *netip.Addr
}

type SessionTokens struct {
	AuthToken    string
	RefreshToken string
}

type LoginResult struct {
	User      User
	SessionID uuid.UUID
	SessionTokens
}
