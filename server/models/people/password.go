package people

import (
	"context"
	"darco/proto/db"
	"darco/proto/models/settings"
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

type PasswordSensitiveInfos struct {
	Email     string
	Login     string
	FirstName string
	LastName  string
}

func (p *PasswordSensitiveInfos) Slice() []string {
	return []string{p.Email, p.Login, p.FirstName, p.LastName}
}

// Validates that a password has a high enough entropy.
// Strength score ranges from 1 to 5.
//
// Returns error if password is too weak.
func ValidatePasswordStrength(p PasswordSensitiveInfos, pwd string, strength int) (*validatedPassword, error) {
	score := zxcvbn.PasswordStrength(pwd, p.Slice()).Score
	if score < strength {
		return nil, fmt.Errorf(
			"new password has strength score %d, at least %d is required",
			score, strength,
		)
	}
	return &validatedPassword{pwd}, nil
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
	maybeValidPwd, err := ValidatePasswordStrength(
		user.PasswordSensitiveInfos(), pwd, int(settings.Security().MinPasswordStrength),
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
	if err := db.Execute(context.Background(),
		`update (<people::User><uuid>$0) set { password := <str>$1 }`,
		user.ID, pwd.value,
	); err != nil {
		logrus.Fatalf(
			`Query failed while attempting to set password for user %v.\nError: %v`,
			user.ID, err,
		)
	}
}
