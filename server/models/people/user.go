package people

import (
	"context"
	_ "embed"
	"encoding/json"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models/settings"
	"github.com/lsdch/biome/services/email"

	"github.com/a-h/templ"
)

type UserInner struct {
	ID             geltypes.UUID `gel:"id" json:"id" binding:"required" format:"uuid"`
	Email          string        `gel:"email" json:"email" binding:"required" format:"email"`
	Login          string        `gel:"login" json:"login" binding:"required"`
	Password       string        `gel:"password" json:"-"`
	Role           UserRole      `gel:"role" json:"role" binding:"required"`
	EmailConfirmed bool          `gel:"email_confirmed" json:"email_confirmed" binding:"required"`
}

type OptionalUserInner struct {
	geltypes.Optional
	UserInner `gel:"$inline" json:",inline"`
}

type User struct {
	UserInner `gel:"$inline" json:",inline"`
	Person    PersonInner `gel:"identity" json:"identity" binding:"required"`
}

type OptionalUser struct {
	geltypes.Optional
	User `gel:"$inline" json:",inline"`
}

func (user *User) SetIdentity(db geltypes.Executor, person *PersonInner) error {
	return db.QuerySingle(context.Background(),
		`#edgeql
			with module people
			select (update (<User><uuid>$0) set {
				identity := assert_single((select Person filter .id = <uuid>$1))
			}) { *, identity: { * }}
		`, user, user.ID, person.ID,
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

func (user *User) SendEmail(subject string, template templ.Component) error {
	emailData := &email.EmailData{
		To:       user.Email,
		From:     settings.Email().FromHeader(),
		Subject:  subject,
		Template: template,
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
func (user *User) Delete(db geltypes.Executor) (*User, error) {
	return DeleteUser(db, user.ID)
}

// Delete a user account using its UUID
func DeleteUser(db geltypes.Executor, uuid geltypes.UUID) (*User, error) {
	var user User
	err := db.QuerySingle(context.Background(),
		`select (delete (<people::User><uuid>$0)) { *, identity: { * }}`,
		&user, uuid,
	)
	return &user, err
}

// Find a user by UUID
//
// Returns geltypes.NoDataError if nothing matches
func FindID(db geltypes.Executor, uuid geltypes.UUID) (user User, err error) {
	err = db.QuerySingle(context.Background(),
		`select (<people::User><uuid>$0) { * , identity: { * } } limit 1`,
		&user,
		geltypes.UUID(uuid),
	)
	return
}

// Find a user by login or email
//
// Returns geltypes.NoDataError if nothing matches
func Find(db *gel.Client, identifier string) (user User, err error) {
	err = db.QuerySingle(context.Background(),
		`select people::User { * , identity: { * } }
		filter .email = <str>$0 or .login = <str>$0 limit 1`,
		&user,
		identifier,
	)
	return
}
