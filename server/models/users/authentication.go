package users

import (
	"errors"
	"fmt"

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
	query := fmt.Sprintf(
		`%s filter (.email = <str>$0 or .login = <str>$0)
			and .password = ext::pgcrypto::crypt(<str>$1, .password)
			limit 1`,
		userSelect,
	)
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

	if !user.Verified {
		return user, &LoginFailedError{AccountInactive}
	}

	return user, nil
}
