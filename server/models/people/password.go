package people

import (
	"net/url"

	"github.com/geldata/gel-go"
	"github.com/lsdch/biome/models/tokens"
	email_templates "github.com/lsdch/biome/templates"

	"github.com/sirupsen/logrus"
)

// RequestPasswordReset creates a password reset token in the DB and sends it
// to the e-mail registered for the user account.
// It can then be used to set a new password for the account.
func (user *User) RequestPasswordReset(db *gel.Client, target url.URL) error {

	token := tokens.NewPwdResetToken(user.ID)
	if err := token.Save(db); err != nil {
		return err
	}
	query := target.Query()
	query.Add("token", string(token.Token))
	target.RawQuery = query.Encode()

	logrus.Debugf("Sending password reset email to %s", user.Person.FullName)
	template := email_templates.PasswordReset(email_templates.PasswordResetData{
		Name: user.Person.FirstName,
		URL:  target,
	})
	return user.SendEmail("Reset your account password", template)
}
