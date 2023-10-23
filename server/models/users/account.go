package users

import (
	"context"
	"darco/proto/config"
	"darco/proto/models"
	"darco/proto/router"
	"darco/proto/services/email"
	"darco/proto/services/tokens"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
)

//go:embed register_user.edgeql
var queryRegister string

func (newUser *UserInput) Register(config *config.Config) error {
	userInsert, err := newUser.ProcessPassword()
	if err != nil {
		return err
	}

	var createdUser User
	args, _ := json.Marshal(userInsert)
	err = models.DB.QuerySingle(context.Background(), queryRegister, &createdUser, args)
	if err != nil {
		return err
	}

	return createdUser.SendConfirmationEmail()
}

// sendConfirmationEmail sends a confirmation email to a newly registered user.
//
// This function generates a confirmation token using the user's email address and a configured token lifetime.
// As opposite to authentication tokens which are generated using the user's UUID, the generate token is generated using their email address, so that they can not be interchanged.
//
// Returns:
//   - err: An error if any issues occur during token generation or email sending; otherwise, it returns nil.
func (user *User) SendConfirmationEmail() (err error) {
	config := config.Get()
	confirmation_token, err := tokens.GenerateToken(user.Email, config.Emailer.TokenLifetime)
	if err != nil {
		return err
	}

	URL := url.URL{
		Scheme:   "https",
		Host:     fmt.Sprintf("%s:%d", config.DomainName, config.Port),
		RawQuery: fmt.Sprintf("token=%s", confirmation_token),
		Path:     path.Join(router.Config.BasePath, "/users/confirm"),
	}

	emailData := &email.EmailData{
		To:       user.Email,
		Subject:  "Your account email verification",
		Template: "email_verification.html",
		Data: map[string]interface{}{
			"Name": user.Person.FirstName,
			"URL":  URL.String(),
		},
	}

	err = email.Send(&config.Emailer, emailData)
	return
}

type PasswordReset struct {
	User struct {
		ID edgedb.UUID `edgedb:"id"`
	} `edgedb:"user"`
	Token   string    `edgedb:"token"`
	Expires time.Time `edgedb:"expires"`
}

func (pwd *PasswordReset) IsValid() bool {
	return pwd.Expires.After(time.Now())
}

//go:embed register_user.edgeql
var queryCreatePasswordReset string

func (user *User) RequestPasswordReset() (err error) {
	config := config.Get()
	token := randstr.String(20)
	expiration := time.Now().Add(config.Emailer.TokenLifetime)
	if err = models.DB.Execute(context.Background(), queryCreatePasswordReset, user.ID, token, expiration); err != nil {
		return
	}
	emailData := &email.EmailData{
		To:       user.Email,
		Subject:  "Reset your account password",
		Template: "email_password_reset.html",
		Data:     map[string]interface{}{},
	}
	err = email.Send(&config.Emailer, emailData)
	return
}

func ValidatePasswordResetToken(token string) (uuid.UUID, bool) {
	query := `select people::PasswordReset { user: {id}, token, expires }
		filter .token = <str>$0`
	var pwdReset PasswordReset
	err := models.DB.QuerySingle(context.Background(), query, pwdReset, token)
	if err != nil {
		return uuid.Nil, false
	}
	return uuid.UUID(pwdReset.User.ID), pwdReset.IsValid()
}

func SetPassword(userID uuid.UUID, pwd *PasswordInput) error {
	hashed_password, err := hashPassword(pwd.Password)
	if err != nil {
		return err
	}
	query := `with module people
		update User filter .id = <uuid>$0
		set {
			password = <str>$1
		}`
	return models.DB.Execute(context.Background(), query, userID, hashed_password)
}
