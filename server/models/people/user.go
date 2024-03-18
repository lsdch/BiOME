package people

import (
	"context"
	"darco/proto/models"
	"darco/proto/services/email"
	_ "embed"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type User struct {
	ID             edgedb.UUID    `edgedb:"id" json:"id" binding:"required"`
	Email          string         `edgedb:"email" json:"email" binding:"required"`
	Login          string         `edgedb:"login" json:"login" binding:"required"`
	Password       string         `edgedb:"password" json:"-"`
	Role           UserRole       `edgedb:"role" json:"role" binding:"required"`
	EmailConfirmed bool           `edgedb:"email_confirmed" json:"email_confirmed" binding:"required"`
	Person         OptionalPerson `edgedb:"identity" json:"identity" binding:"required"`
	IsActive       bool           `edgedb:"is_active" json:"is_active" binding:"required"`
	models.Meta    `edgedb:"meta" json:"meta"`
} //@name User

func (user *User) SetIdentity(db edgedb.Executor, person *Person) error {
	return db.QuerySingle(context.Background(),
		`with module people
			select (update (<User><uuid>$0) set {
				identity := assert_single((select Person filter .id = <uuid>$1))
			}) { *, identity: { * }}`,
		user, user.ID, person.ID,
	)
}

func (user *User) PasswordSensitiveInfos() PasswordSensitiveInfos {
	return PasswordSensitiveInfos{
		Email:     user.Email,
		Login:     user.Login,
		FirstName: user.Person.FirstName,
		LastName:  user.Person.LastName,
	}
}

func (user *User) SendEmail(subject string, template_file string, data map[string]any) error {
	emailData := &email.EmailData{
		To:       user.Email,
		Subject:  subject,
		Template: template_file,
		Data:     data,
	}

	return email.Send(email.AdminEmailAddress(), emailData)
}

func (user *User) SetEmailConfirmed(db *edgedb.Client, active bool) error {
	logrus.Infof("Confirm email address '%s' for user '%s'",
		user.Login, user.Email,
	)
	return db.QuerySingle(context.Background(),
		`select (update (<people::User><uuid>$0)
			set { email_confirmed := <bool>$1 }
		) { *, identity: { * }};`,
		user, user.ID, active)
}

// Custom user marshaller that obfuscates password
func (user *User) MarshalJSON() ([]byte, error) {
	u := User(*user)
	u.Password = "**********"
	return json.Marshal(u)
}

// Deletes a user account
func (user *User) Delete(db edgedb.Executor) (*User, error) {
	return DeleteUser(db, user.ID)
}

// Delete a user account using its UUID
func DeleteUser(db edgedb.Executor, uuid edgedb.UUID) (*User, error) {
	var user User
	err := db.QuerySingle(context.Background(),
		`select (delete (<people::User><uuid>$0)) { *, identity: { * }}`,
		&user, uuid,
	)
	return &user, err
}

// Convenience function to retrieve a user using a database query with arguments.
func find(db *edgedb.Client, query string, args ...interface{}) (*User, error) {
	var user User
	if err := db.QuerySingle(context.Background(), query, &user, args...); err != nil {
		return nil, err
	}
	return &user, nil
}

// Find a user by UUID
//
// Returns edgedb.NoDataError if nothing matches
func FindID(db *edgedb.Client, uuid edgedb.UUID) (*User, error) {
	return find(db,
		`select (<people::User><uuid>$0) { * , identity: { * } } limit 1`,
		edgedb.UUID(uuid),
	)
}

// Find a user by login or email
//
// Returns edgedb.NoDataError if nothing matches
func Find(db *edgedb.Client, identifier string) (*User, error) {
	return find(db,
		`select people::User { * , identity: { * } }
		filter .email = <str>$0 or .login = <str>$0 limit 1`,
		identifier)
}
