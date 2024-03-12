package users

import (
	"context"
	"darco/proto/db"
	"darco/proto/models/people"
	"darco/proto/models/users/user_role"
	_ "embed"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type PersonUserInput struct {
	people.PersonInner `edgedb:"$inline"`
	MiddleNames        *string `json:"middle_names,omitempty" edgedb:"middle_names" binding:"omitnil,alphaunicode,max=32"`
	Contact            *string `json:"contact,omitempty" binding:"omitnil,nullemail"`
} // @name PersonUserInput

type InnerUserInput struct {
	Login  string          `edgedb:"login" json:"login" binding:"login,required,unique_login"`
	Email  string          `edgedb:"email" json:"email" binding:"email,required,unique_email" format:"email"`
	Person PersonUserInput `edgedb:"identity" json:"identity" binding:"required"`
} // @name InnerUserInput

type UserInput struct {
	InnerUserInput `json:",inline"`
	PasswordInput  `json:",inline"`
} // @name UserInput

type User struct {
	ID       edgedb.UUID        `edgedb:"id" json:"-" binding:"required"`
	Email    string             `edgedb:"email" json:"email" binding:"required"`
	Login    string             `edgedb:"login" json:"login" binding:"required"`
	Password string             `edgedb:"password" json:"-"`
	Role     user_role.UserRole `edgedb:"role" json:"role" binding:"required"`
	Verified bool               `edgedb:"verified" json:"verified" binding:"required"`
	Person   people.Person      `edgedb:"identity" json:"identity" binding:"required"`
} //@name User

func (user *User) InnerUserInput() *InnerUserInput {
	return &InnerUserInput{
		Login: user.Login,
		Email: user.Email,
		Person: PersonUserInput{
			PersonInner: user.Person.PersonInner,
			MiddleNames: db.OptionalAsPointer(user.Person.MiddleNames),
			Contact:     db.OptionalAsPointer(user.Person.Contact),
		},
	}
}

func (user *User) SetActive(db *edgedb.Client, active bool) error {
	logrus.Infof("Activating account of %s %s (%s)",
		user.Person.FirstName, user.Person.LastName, user.Email,
	)
	query := `update people::User filter .email = <str>$1 set { verified := <bool>$0 };`
	return db.Execute(context.Background(), query, active, user.Email)
}

// Custom user marshaller that obfuscates password
func (user *User) MarshalJSON() ([]byte, error) {
	u := User(*user)
	u.Password = "**********"
	return json.Marshal(u)
}

func find(db *edgedb.Client, query string, args ...interface{}) (*User, error) {
	var user User
	if err := db.QuerySingle(context.Background(), query, &user, args...); err != nil {
		return nil, err
	}
	return &user, nil
}

// Returns currently authenticated user or edgedb.NoDataError if not authenticated
func Current(db *edgedb.Client) (*User, error) {
	query := `select people::User { * , identity: { * } }
		filter .id = global current_user_id limit 1`
	return find(db, query)
}

// Find a user by UUID
//
// Returns edgedb.NoDataError if nothing matches
func FindID(db *edgedb.Client, uuid edgedb.UUID) (*User, error) {
	query := `
	select people::User { * , identity: { * } }
		filter .id = <uuid>$0 limit 1`
	return find(db, query, edgedb.UUID(uuid))
}

// Find a user by login or email
//
// Returns edgedb.NoDataError if nothing matches
func Find(db *edgedb.Client, identifier string) (*User, error) {
	query := `select people::User { * , identity: { * } }
		filter .email = <str>$0 or .login = <str>$0 limit 1`
	return find(db, query, identifier)
}
