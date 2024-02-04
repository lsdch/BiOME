package users

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/person"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/edgedb/edgedb-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type InnerUserInput struct {
	Login string `edgedb:"login" json:"login" binding:"login,required,unique_login"`
	Email string `edgedb:"email" json:"email" binding:"email,required,unique_email" format:"email"`
	// EmailPublic bool        `edbedb:"email_public" json:"email_public"`
	Person person.PersonInput `edgedb:"identity" json:"identity" binding:"required"`
} // @name InnerUserInput
type UserInput struct {
	InnerUserInput `json:",inline"`
	PasswordInput  `json:",inline"`
} // @name UserInput

func (input *UserInput) ProcessPassword() (*UserInsert, error) {
	hashed_password, err := hashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	return &UserInsert{
		InnerUserInput: input.InnerUserInput,
		PasswordHash:   hashed_password,
	}, nil
}

type UserInsert struct {
	InnerUserInput `edgedb:"$inline" json:",inline"`
	PasswordHash   string `edgedb:"password" json:"password"`
}

type UserPartial struct {
	Role     UserRole      `edgedb:"role" json:"role" binding:"required"`
	Verified bool          `edgedb:"verified" json:"verified" binding:"required"`
	Person   person.Person `edgedb:"identity" json:"identity" binding:"required"`
} // @name UserPartial

type User struct {
	UserPartial `edgedb:"$inline" json:",inline"`
	ID          edgedb.UUID `edgedb:"id" json:"-" binding:"required"`
	Email       string      `edgedb:"email" json:"email" binding:"required"`
	Login       string      `edgedb:"login" json:"login" binding:"required"`
	Password    string      `edgedb:"password" json:"-"`
} //@name User

func (user *User) Partial() *UserPartial {
	return &user.UserPartial
}

func (user *User) InnerUserInput() *InnerUserInput {
	return &InnerUserInput{
		Login:  user.Login,
		Email:  user.Email,
		Person: user.Person.PersonInput,
	}
}

func (user *User) SetActive(active bool) error {
	logrus.Infof("Activating account of %s %s (%s)",
		user.Person.FirstName, user.Person.LastName, user.Email,
	)
	query := `update people::User
	filter .email = <str>$1
	set {
		verified := <bool>$0
	} `
	return models.DB().Execute(context.Background(), query, active, user.Email)
}

func (user *User) MarshalJSON() ([]byte, error) {
	u := User(*user)
	u.Password = "**********"
	return json.Marshal(u)
}

type UserCredentials struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Remember   bool   `json:"remember" binding:"required"`
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

func (creds *UserCredentials) Authenticate(db *edgedb.Client) (*User, *LoginFailedError) {
	user, err := Find(db, creds.Identifier)

	if err != nil {
		var dbErr edgedb.Error
		if errors.As(err, &dbErr) && dbErr.Category(edgedb.NoDataError) {
			return nil, &LoginFailedError{InvalidCredentials}
		} else {
			logrus.Errorf("%v", err)
			return nil, &LoginFailedError{ServerError}
		}
	}

	if err := VerifyPassword(user.Password, creds.Password); err != nil {
		return user, &LoginFailedError{InvalidCredentials}
	}

	if !user.Verified {
		return user, &LoginFailedError{AccountInactive}
	}

	return user, nil
}

var userSelect = `select people::User {
	id, login, email, verified, role, password, identity: { * }
}`

func find(db *edgedb.Client, query string, args ...interface{}) (*User, error) {
	var user User
	if err := db.QuerySingle(context.Background(), query, &user, args...); err != nil {
		return nil, err
	}
	return &user, nil
}

func Current(db *edgedb.Client) (*User, error) {
	query := fmt.Sprintf(`%s filter .id = global current_user_id limit 1`, userSelect)
	return find(db, query)
}

// Find a user by UUID
//
// Returns edgedb.NoDataError if nothing matches
func FindID(db *edgedb.Client, uuid uuid.UUID) (*User, error) {
	query := fmt.Sprintf(`%s filter .id = <uuid>$0 limit 1`, userSelect)
	return find(db, query, edgedb.UUID(uuid))
}

// Find a user by login or email
//
// Returns edgedb.NoDataError if nothing matches
func Find(db *edgedb.Client, identifier string) (*User, error) {
	query := fmt.Sprintf(`%s filter .email = <str>$0 or .login = <str>$0 limit 1`, userSelect)
	return find(db, query, identifier)
}
