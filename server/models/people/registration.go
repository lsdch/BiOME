package people

import (
	"context"
	_ "embed"
	"encoding/json"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type UserInput struct {
	Login         string `edgedb:"login" json:"login" binding:"login,required,unique_login" faker:"username,unique"`
	Email         string `edgedb:"email" json:"email" binding:"email,required,unique_email" format:"email" faker:"email,unique"`
	PasswordInput `json:",inline"`
} // @name UserInput

type InnerPendingUserRequest struct {
	Person struct {
		PersonInner `edgedb:"$inline" json:",inline"`
		Institution string `edgedb:"institution" json:"institution,omitempty" faker:"word"`
	} `json:"identity" edgedb:"identity"`
	Motive string `json:"motive" edgedb:"motive" faker:"sentence"`
} // @name InnerPendingUserRequest

type PendingUserRequestInput struct {
	User                    UserInput `json:"user"`
	InnerPendingUserRequest `json:",inline"`
} // @name PendingUserRequestInput

//go:embed queries/register_pending_user.edgeql
var registerPendingUserQuery string

// Creates an inactive user without identity, and an account request
// with personal informations which can be validated by an admin.
func (u *PendingUserRequestInput) Register(db edgedb.Executor) (*PendingUserRequest, error) {
	args, _ := json.Marshal(u)
	var pendingUser PendingUserRequest
	err := db.QuerySingle(context.Background(), registerPendingUserQuery, &pendingUser, args)
	return &pendingUser, err
}

type PendingUserRequest struct {
	ID                      edgedb.UUID `edgedb:"id"`
	User                    User        `json:"user" edgedb:"user"`
	InnerPendingUserRequest `json:",inline" edgedb:"$inline"`
	CreatedOn               time.Time `json:"created_on" edgedb:"created_on"`
} // @name PendingUserRequest

func (p *PendingUserRequest) Delete(db edgedb.Executor) error {
	return db.Execute(context.Background(),
		`delete <people::PendingUserRequest><uuid>$0;`,
		p.ID,
	)
}

// Validate an inactive user by setting their identity to an existing person.
// PendingUserRequest is deleted in the database afterwards.
// All operations are done within a database transaction.
func (p *PendingUserRequest) Validate(db *edgedb.Client, person *Person) (*User, error) {
	var user *User
	var err error
	db.Tx(context.Background(), func(ctx context.Context, tx *edgedb.Tx) error {
		user, err = p.ValidateTx(tx, person)
		return err
	})
	return user, err
}

// Like [*PendingUserRequest.ValidateTx] but the transaction executor is provided as argument
func (p *PendingUserRequest) ValidateTx(tx *edgedb.Tx, person *Person) (*User, error) {
	if err := p.User.SetIdentity(tx, person); err != nil {
		return nil, err
	}
	if err := p.Delete(tx); err != nil {
		return nil, err
	}
	logrus.Infof(
		"Pending user request validated for user '%+v'.\nAssigned identity %+v",
		p.User.IsActive, p.User.Person,
	)
	return &p.User, nil
}

// SendConfirmationEmail sends a confirmation email to the user with a verification token.
// It generates a confirmation token, and sends an email with the confirmation link.
// The confirmation token is included as a query parameter in the URL.
func (user *User) SendConfirmationEmail(db *edgedb.Client, target TokenURL) error {
	token, err := user.CreateAccountToken(db, EmailConfirmationToken)
	if err != nil {
		return err
	}
	target.SetToken(token)

	return user.SendEmail(
		"Activation of your account",
		"email_verification.html",
		map[string]any{
			"Name": user.Person.FirstName,
			"URL":  target.Encode().String(),
		})
}

func (user *User) RequestPasswordReset(db *edgedb.Client, tokenURL TokenURL) error {
	token, err := user.CreateAccountToken(db, PasswordResetToken)
	if err != nil {
		return err
	}
	tokenURL.SetToken(token)
	return user.SendEmail(
		"Reset your account password",
		"email_password_reset.html",
		map[string]interface{}{
			"Name": user.Person.FirstName,
			"URL":  tokenURL.String(),
		})
}
