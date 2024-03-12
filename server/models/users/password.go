package users

import (
	"context"
	"darco/proto/config"
	"darco/proto/db"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/trustelem/zxcvbn"
)

type NewPasswordInput struct {
	Password    string        `json:"password" binding:"required"` // Old password
	NewPassword PasswordInput `json:"new_password" binding:"required"`
} // @name NewPasswordInput

type PasswordInput struct {
	Password   string `json:"password" binding:"required"`
	ConfirmPwd string `json:"password_confirmation" binding:"eqfield=Password,required"`
} //@name PasswordInput

// The raw string value of a password after it has successfully passed all validation checks.
type validatedPassword struct {
	value string
}

// Validates that a password has a high enough entropy.
// Strength score ranges from 1 to 5.
//
// Returns error if password is too weak.
func (u *InnerUserInput) ValidatePasswordStrength(pwd string, strength int) (*validatedPassword, error) {
	score := zxcvbn.PasswordStrength(pwd, u.PasswordSensitiveValues()).Score
	if score < strength {
		return nil, fmt.Errorf(
			"new password has strength score %d, at least %d is required",
			score, strength,
		)
	}
	return &validatedPassword{pwd}, nil
}

// Used to compute password strength with regard to personal infos
func (u *InnerUserInput) PasswordSensitiveValues() []string {
	return []string{u.Email, u.Login, u.Person.FirstName, u.Person.LastName}
}

// Checks that a password matches the hashed password for a user.
func (user *User) PasswordMatch(db db.Executor, pwd string) bool {
	var match bool
	query := `select exists (select people::User
			filter .id = <uuid>$0
			and .password = ext::pgcrypto::crypt(<str>$1, .password)
		);`
	if err := db.QuerySingle(context.Background(), query, &match, user.ID, pwd); err != nil {
		logrus.Fatalf("Password matching query failed: %v", err)
	}
	return match
}

// Sets the password of a user.
// Returns string error if password is not strong enough.
func (user *User) SetPassword(db db.Executor, pwd string) error {
	maybeValidPwd, err := user.InnerUserInput().ValidatePasswordStrength(
		pwd, config.Get().Accounts.PasswordStrength,
	)
	if err != nil {
		return err
	}

	user.setPassword(db, *maybeValidPwd)
	return nil
}

// Sets the password for a user in DB.
// The query failing indicates a programming error, which results in a fatal error.
func (user *User) setPassword(db db.Executor, pwd validatedPassword) {
	query := `update people::User filter .id = <uuid>$0 set { password := <str>$1 }`
	if err := db.Execute(context.Background(), query, user.ID, pwd.value); err != nil {
		logrus.Fatalf(
			`Query failed while attempting to set password for user %d.\nQuery: %s\nError: %v`,
			user.ID, query, err,
		)
	}
}
