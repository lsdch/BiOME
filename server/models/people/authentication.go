package people

import (
	"context"
	"darco/proto/config"
	"darco/proto/models/tokens"
	"darco/proto/services/auth_tokens"
	"net/http"
	"time"

	"github.com/edgedb/edgedb-go"
)

type UserCredentials struct {
	Identifier string `json:"identifier" binding:"required"`
	// Unhashed, password hash comparison is done within EdgeDB
	Password string `json:"password" binding:"required"`
}

// Attempts to authenticate a user given the provided credentials.
// Return edgedb.NoData error if credentials are invalid.
func (creds *UserCredentials) Authenticate(db *edgedb.Client) (user User, err error) {
	err = db.QuerySingle(context.Background(),
		`select people::User { *, identity: { * } }
			filter (.email = <str>$0 or .login = <str>$0)
			and .password = ext::pgcrypto::crypt(<str>$1, .password)
			limit 1`,
		&user,
		creds.Identifier, creds.Password,
	)
	return
}

// Returns currently authenticated user or edgedb.NoDataError if not authenticated
func Current(db *edgedb.Client) (user User, err error) {
	err = db.QuerySingle(context.Background(),
		`select (global current_user) { * , identity: { * } } limit 1`,
		&user,
	)
	return
}

func (user *User) GenerateJWT() (tokens.TokenRecord, error) {
	lifetime :=
		config.Get().AuthTokenDuration()
	token, err := auth_tokens.NewJWT(user.ID, lifetime)
	if err != nil {
		return tokens.TokenRecord{}, err
	}
	return tokens.TokenRecord{
		Token:   tokens.Token(token),
		Expires: time.Now().Add(lifetime),
	}, nil
}

func (user *User) JWTCookie(jwt string, domain string) http.Cookie {
	return http.Cookie{
		Name:     auth_tokens.AUTH_TOKEN_COOKIE,
		Value:    jwt,
		Path:     "/",
		Domain:   domain,
		MaxAge:   int(config.Get().AuthTokenDuration().Seconds()),
		Secure:   true,
		HttpOnly: true,
	}
}
