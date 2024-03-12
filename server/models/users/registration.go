package users

import (
	"context"
	"darco/proto/config"
	"darco/proto/services/email"
	_ "embed"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
)

//go:embed queries/register_user.edgeql
var queryRegister string

// Registers new account and sends an email with an activation link
func (newUser *UserInput) Register(db *edgedb.Client) (createdUser *User, err error) {
	args, _ := json.Marshal(newUser)
	if err := db.QuerySingle(context.Background(), queryRegister, createdUser, args); err != nil {
		return nil, err
	}
	return
}

// SendConfirmationEmail sends a confirmation email to the user with a verification token.
// It generates a confirmation token, and sends an email with the confirmation link.
// The confirmation token is included as a query parameter in the URL.
func (user *User) SendConfirmationEmail(db *edgedb.Client, tokenURL TokenURL) (err error) {
	token, err := user.CreateAccountToken(db, EmailConfirmationToken)
	if err != nil {
		return err
	}
	tokenURL.SetToken(token)
	emailData := &email.EmailData{
		To:       user.Email,
		Subject:  "Your account email verification",
		Template: "email_verification.html",
		Data: map[string]interface{}{
			"Name": user.Person.FirstName,
			"URL":  tokenURL.String(),
		},
	}

	return email.Send(&config.Get().Emailer, emailData)
}

func (user *User) RequestPasswordReset(db *edgedb.Client, tokenURL TokenURL) (err error) {
	token, err := user.CreateAccountToken(db, PasswordResetToken)
	if err != nil {
		return err
	}
	tokenURL.SetToken(token)
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
