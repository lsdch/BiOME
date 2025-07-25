package people

import (
	"context"
	"net/http"
	"time"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/models/tokens"
	"github.com/lsdch/biome/services/auth_tokens"
)

type UserCredentials struct {
	Identifier string `json:"identifier" binding:"required"`
	// Unhashed, password hash comparison is done within Gel
	Password string `json:"password" binding:"required"`
}

// Attempts to authenticate a user given the provided credentials.
// Return geltypes.NoDataError if credentials are invalid.
func (creds *UserCredentials) Authenticate(db *gel.Client) (user User, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select people::User { *, identity: { * } }
			filter (.email = <str>$0 or .login = <str>$0)
			and .password = ext::pgcrypto::crypt(<str>$1, .password)
			limit 1
		`,
		&user,
		creds.Identifier, creds.Password,
	)
	return
}

// Returns currently authenticated user or geltypes.NoDataError if not authenticated
func Current(db geltypes.Executor) (user User, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (global current_user) { * , identity: { * } } limit 1
		`,
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
		SameSite: http.SameSiteLaxMode,
	}
}
