package people

import (
	"bytes"
	"context"
	"darco/proto/models/settings"
	"text/template"
	"time"

	_ "embed"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type AccountTokenKind string

const (
	EmailConfirmationToken AccountTokenKind = "EmailConfirmation"
	PasswordResetToken     AccountTokenKind = "PasswordReset"
)

// Template query to create an account token
//
//go:embed queries/upsert_account_token.edgeql
var createAccountTokenTemplate string

func (k AccountTokenKind) generateCreateQuery() (string, error) {
	tmpl, err := template.New("token_query").Parse(createAccountTokenTemplate)
	query := new(bytes.Buffer)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(query, k)
	if err != nil {
		return "", err
	}
	return query.String(), nil
}

// A confirmation token stored in the database
type AccountToken struct {
	ID           edgedb.UUID `edgedb:"id"`
	User         User        `edgedb:"user"`
	TokenWrapper `edgedb:"$inline" json:",inline"`
}

func (t *AccountToken) Save(db edgedb.Executor, kind AccountTokenKind) error {
	query, err := kind.generateCreateQuery()
	if err != nil {
		logrus.Errorf(
			"Failed to generated query to create token kind '%s': %v",
			kind, err)
		return err
	}
	err = db.Execute(context.Background(), query, t.User.ID, t.Token, t.Expires)
	if err != nil {
		logrus.Errorf(
			"Failed to create an account token: %+v; Err:%v",
			*t, err)
		return err
	}
	return nil
}

// Deletes token from the database
func (token *AccountToken) Consume(db *edgedb.Client) (err error) {
	deleteQuery := `delete people::AccountToken filter .id = <uuid>$0`
	if err = db.Execute(context.Background(), deleteQuery, token.ID); err != nil {
		logrus.Errorf("Database error %v (query: %s)", err, deleteQuery)
		return
	}
	return
}

func ValidateAccountToken(db *edgedb.Client, token Token, kind AccountTokenKind) (*User, bool) {
	var db_token AccountToken
	if err := db.QuerySingle(context.Background(),
		`select people::AccountToken { *, user: { *, identity: { * }} }
				  filter .token = <str>$0`,
		&db_token, token,
	); err != nil {
		logrus.Infof("Failed to validate password reset token: %+v", err)
		return nil, false
	}
	logrus.Infof("Retrieved token %+v", db_token)
	if db_token.IsValid() {
		db_token.Consume(db)
	}
	return &db_token.User, db_token.IsValid()
}

// Stores an AccountToken of the specified kind in the database, ready to be consumed.
// Returns the token string which can be shared, e.g. in a link to be emailed.
func (user *User) CreateAccountToken(db *edgedb.Client, kind AccountTokenKind) (Token, error) {
	tok_str := GenerateToken(20)
	expires := time.Now().Add(settings.Security().AccountTokenDuration())
	logrus.Infof(
		"Creating account token %s for user ID %v",
		tok_str, user.ID)
	token := AccountToken{
		User: *user,
		TokenWrapper: TokenWrapper{
			Token:   Token(tok_str),
			Expires: expires,
		},
	}
	if err := token.Save(db, kind); err != nil {
		return "", nil
	}
	return token.Token, nil
}
