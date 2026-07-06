package models

import (
	"net/netip"
	"time"

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
	UserAgent Optional[string] `header:"User-Agent"`
	IP        *netip.Addr      `header:"client-ip"`
}

type SessionTokens struct {
	SessionID           uuid.UUID `json:"session_id"`
	AuthToken           string    `json:"auth_token"`
	AuthTokenExpiration time.Time `json:"auth_token_expiration"`
	RefreshToken        string    `json:"refresh_token"`
}

type LoginResult struct {
	User    User          `json:"user"`
	Session SessionTokens `json:"session"`
}
