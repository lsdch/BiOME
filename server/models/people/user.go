package people

import (
	"context"
	"darco/proto/models/settings"
	"darco/proto/services/email"
	_ "embed"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
)

type UserInner struct {
	ID             edgedb.UUID `edgedb:"id" json:"id" binding:"required" format:"uuid"`
	Email          string      `edgedb:"email" json:"email" binding:"required" format:"email"`
	Login          string      `edgedb:"login" json:"login" binding:"required"`
	Password       string      `edgedb:"password" json:"-"`
	Role           UserRole    `edgedb:"role" json:"role" binding:"required"`
	EmailConfirmed bool        `edgedb:"email_confirmed" json:"email_confirmed" binding:"required"`
}

type OptionalUserInner struct {
	edgedb.Optional
	UserInner `edgedb:"$inline" json:",inline"`
}

type User struct {
	UserInner `edgedb:"$inline" json:",inline"`
	Person    OptionalPerson `edgedb:"identity" json:"identity" binding:"required"`
}

type OptionalUser struct {
	edgedb.Optional
	User `edgedb:"$inline" json:",inline"`
}

func (user *User) SetIdentity(db edgedb.Executor, person *PersonInner) error {
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

	return emailData.Send(settings.Email().FromHeader())
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

// Find a user by UUID
//
// Returns edgedb.NoDataError if nothing matches
func FindID(db edgedb.Executor, uuid edgedb.UUID) (user User, err error) {
	err = db.QuerySingle(context.Background(),
		`select (<people::User><uuid>$0) { * , identity: { * } } limit 1`,
		&user,
		edgedb.UUID(uuid),
	)
	return
}

// Find a user by login or email
//
// Returns edgedb.NoDataError if nothing matches
func Find(db *edgedb.Client, identifier string) (user User, err error) {
	err = db.QuerySingle(context.Background(),
		`select people::User { * , identity: { * } }
		filter .email = <str>$0 or .login = <str>$0 limit 1`,
		&user,
		identifier,
	)
	return
}
