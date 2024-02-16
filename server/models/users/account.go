package users

import (
	"context"
	"darco/proto/config"
	"darco/proto/db"
	"darco/proto/services/email"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

//go:embed queries/register_user.edgeql
var queryRegister string

// Registers new account and sends an email with an activation link
func (newUser *UserInput) Register(config *config.Config, originURL *url.URL) error {
	userInsert, err := newUser.ProcessPassword()
	if err != nil {
		return err
	}

	var createdUser User
	args, _ := json.Marshal(userInsert)
	if err = db.Client().QuerySingle(context.Background(), queryRegister, &createdUser, args); err != nil {
		return err
	}
	return createdUser.SendConfirmationEmail(originURL)
}

// SendConfirmationEmail sends a confirmation email to the user with a verification token.
// It generates a confirmation token, constructs a confirmation URL,
// and sends an email with the confirmation link.
//
// Note:
//
// The confirmation URL is constructed based on a request origin URL.
// If the origin URL matches the configured client's host,
// the confirmation URL will use the client's base URL.
// Otherwise, it defaults to an API endpoint URL ("/users/confirm").
// The confirmation token is included as a query parameter in the URL.
func (user *User) SendConfirmationEmail(originUrl *url.URL) (err error) {
	token, err := user.CreateConfirmationToken()
	if err != nil {
		return err
	}

	tokenURL := makeTokenURL(token, "/users/confirm", "/email-confirmation", originUrl)

	emailData := &email.EmailData{
		To:       user.Email,
		Subject:  "Your account email verification",
		Template: "email_verification.html",
		Data: map[string]interface{}{
			"Name": user.Person.FirstName,
			"URL":  tokenURL.String(),
		},
	}

	err = email.Send(&config.Get().Emailer, emailData)
	return
}

func makeTokenURL(token Token, api_path string, client_path string, originURL *url.URL) url.URL {
	var (
		config    = config.Get()
		tokenURL  url.URL
		clientURL = config.MakeClientURL("/email-confirmation")
	)
	if (originURL != nil) && originURL.Host == clientURL.Host {
		tokenURL = clientURL
	} else {
		tokenURL = config.MakeURL("/users/confirm")
	}
	tokenURL.RawQuery = fmt.Sprintf("token=%s", token)
	return tokenURL
}

func (user *User) RequestPasswordReset(originURL *url.URL) (err error) {
	token, err := user.CreatePasswordResetToken()
	if err != nil {
		return err
	}

	tokenURL := makeTokenURL(token, "/users/password-reset", "/password-reset", originURL)
	emailData := &email.EmailData{
		To:       user.Email,
		Subject:  "Reset your account password",
		Template: "email_password_reset.html",
		Data: map[string]interface{}{
			"Name": user.Person.FirstName,
			"URL":  tokenURL.String(),
		},
	}

	return email.Send(&config.Get().Emailer, emailData)
}

func SetPassword(userID uuid.UUID, pwd *PasswordInput) error {
	hashed_password, err := hashPassword(pwd.Password)
	if err != nil {
		return err
	}
	query := `with module people
		update User filter .id = <uuid>$0
		set { password = <str>$1 }`
	return db.Client().Execute(context.Background(), query, userID, hashed_password)
}
