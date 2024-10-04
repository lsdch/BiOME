package people

import (
	"context"
	"darco/proto/models/settings"
	"darco/proto/models/tokens"
	"net/url"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
	"github.com/trustelem/zxcvbn"
)

type UpdatePasswordInput struct {
	Password    string        `json:"password" binding:"required"` // Old password
	NewPassword PasswordInput `json:"new_password" binding:"required"`
}

type PasswordInput struct {
	Password   string `json:"password"`
	ConfirmPwd string `json:"password_confirmation"`
}

func (p PasswordInput) ValidateEqual() bool {
	return p.Password == p.ConfirmPwd
}

type PasswordSensitiveInfos struct {
	Email     string
	Login     string
	FirstName string
	LastName  string
}

func (p *PasswordSensitiveInfos) ToSlice() []string {
	return []string{p.Email, p.Login, p.FirstName, p.LastName}
}

// Validates that a password has a high enough entropy.
// Strength score ranges from 1 to 5.
//
// Returns error if password is too weak.
func validatePasswordStrength(p PasswordSensitiveInfos, pwd string, strength int) bool {
	score := zxcvbn.PasswordStrength(pwd, p.ToSlice()).Score
	return score > strength
}

// Checks that a password matches the hashed password for a user.
func (user *User) PasswordMatch(db edgedb.Executor, pwd string) bool {
	var match bool
	query := `select exists (select people::User
			filter .id = <uuid>$0
			and .password = ext::pgcrypto::crypt(<str>$1, .password)
		);`
	if err := db.QuerySingle(context.Background(), query, &match, user.ID, pwd); err != nil {
		logrus.Errorf("Password matching query failed: %v", err)
	}
	return match
}

// Sets the password of a user.
// Returns string error if password is not strong enough.
func (user *User) ValidatePasswordStrength(db edgedb.Executor, pwd string) bool {
	strongEnough := validatePasswordStrength(
		user.PasswordSensitiveInfos(),
		string(pwd),
		int(settings.Security().MinPasswordStrength),
	)
	return strongEnough

}

// Sets the password for a user in DB.
func (user *User) SetPassword(db edgedb.Executor, pwd string) error {
	return db.Execute(context.Background(),
		`update (<people::User><uuid>$0) set { password := <str>$1 }`,
		user.ID, pwd,
	)
}

// RequestPasswordReset creates a password reset token in the DB and sends it
// to the e-mail registered for the user account.
// It can then be used to set a new password for the account.
func (user *User) RequestPasswordReset(db *edgedb.Client, target url.URL) error {

	token := tokens.NewPwdResetToken(user.ID)
	if err := token.Save(db); err != nil {
		return err
	}

	target.RawQuery = url.QueryEscape(string(token.Token))
	return user.SendEmail(
		"Reset your account password",
		"email_password_reset.html",
		map[string]interface{}{
			"Name": user.Person.FirstName,
			"URL":  target.String(),
		})
}
