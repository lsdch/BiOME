package people

import (
	"darco/proto/models/settings"
	"darco/proto/services/tokens"
	"errors"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type UserCredentials struct {
	Identifier string `json:"identifier" binding:"required"`
	// Unhashed, password hash comparison is done within EdgeDB
	Password string `json:"password" binding:"required"`
} // @name UserCredentials

type loginFailedReason string // @name LoginFailedReason

const (
	AccountInactive    loginFailedReason = "Inactive"
	InvalidCredentials loginFailedReason = "InvalidCredentials"
	ServerError        loginFailedReason = "ServerError"
)

type LoginFailedError struct {
	Reason loginFailedReason `json:"reason" binding:"required"`
} // @name LoginFailedError

func (err *LoginFailedError) Error() string {
	return string(err.Reason)
}

// Attempts to authenticate a user given the provided credentials.
func (creds *UserCredentials) Authenticate(db *edgedb.Client) (*User, *LoginFailedError) {
	query := `select people::User { *, identity: { * } }
			filter (.email = <str>$0 or .login = <str>$0)
			and .password = ext::pgcrypto::crypt(<str>$1, .password)
			limit 1`
	user, err := find(db, query, creds.Identifier, creds.Password)

	if err != nil {
		var dbErr edgedb.Error
		if errors.As(err, &dbErr) && dbErr.Category(edgedb.NoDataError) {
			return nil, &LoginFailedError{InvalidCredentials}
		} else {
			logrus.Errorf("%v", err)
			return nil, &LoginFailedError{ServerError}
		}
	}

	if !user.IsActive {
		return user, &LoginFailedError{AccountInactive}
	}

	return user, nil
}

// Returns currently authenticated user or edgedb.NoDataError if not authenticated
func Current(db *edgedb.Client) (*User, error) {
	return find(db, `select (global current_user) { * , identity: { * } } limit 1`)
}

func (user *User) GenerateJWT() (string, error) {
	return tokens.GenerateToken(user.ID,
		settings.Security().AuthTokenDuration())
}
