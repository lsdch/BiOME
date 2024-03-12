package users

import (
	"bytes"
	"context"
	"darco/proto/config"
	"fmt"
	"net/url"
	"text/template"
	"time"

	_ "embed"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
	"github.com/thanhpk/randstr"
)

// A confirmation token generated as a random string.
// Used to validate email addresses or password reset requests.
type Token string

func (t *Token) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(*t), nil
}

func (t *Token) UnmarshalEdgeDBStr(data []byte) error {
	*t = Token(string(data))
	return nil
}

// A URL that encodes a token in its query parameters as '?token=<token-string>'
type TokenURL struct {
	url.URL
}

func (url TokenURL) SetToken(token Token) {
	url.RawQuery = fmt.Sprintf("token=%s", token)
}

type AccountTokenKind string

const (
	EmailConfirmationToken AccountTokenKind = "EmailConfirmation"
	PasswordResetToken     AccountTokenKind = "PasswordReset"
)

// A confirmation token stored in the database
type AccountToken struct {
	ID      edgedb.UUID `edgedb:"id"`
	User    *User       `edgedb:"user"`
	Token   Token       `edgedb:"token"`
	Expires time.Time   `edgedb:"expires"`
}

func (token *AccountToken) IsValid() bool {
	return token.Expires.After(time.Now())
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
	query := fmt.Sprintf(
		`select people::%s { id, user: { *, identity: { * }}, token, expires }
		filter .token = <str>$0`, kind,
	)
	var db_token AccountToken
	if err := db.QuerySingle(context.Background(), query, &db_token, token); err != nil {
		logrus.Infof("Failed to validate password reset token: %+v", err)
		return nil, false
	}
	if db_token.IsValid() {
		db_token.Consume(db)
	}
	return db_token.User, db_token.IsValid()
}

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

// Stores an AccountToken of the specified kind in the database, ready to be consumed.
// Returns the token string which can be shared, e.g. in a link to be emailed.
func (user *User) CreateAccountToken(db *edgedb.Client, kind AccountTokenKind) (Token, error) {
	config := config.Get()
	tok_str := Token(randstr.String(20))
	expires := time.Now().Add(config.Emailer.TokenLifetime)
	logrus.Infof("Creating account token %s for user ID %v", tok_str, user.ID)
	token := AccountToken{
		User:    user,
		Token:   Token(tok_str),
		Expires: expires,
	}
	query, err := kind.generateCreateQuery()
	if err != nil {
		logrus.Errorf("Failed to generated query to create token kind '%s': %v", kind, err)
		return "", err
	}
	err = db.Execute(context.Background(), query, user.ID, token.Token, token.Expires)
	if err != nil {
		logrus.Errorf("Failed to create an account token: %+v; Err:%v", token, err)
		return "", err
	}
	return token.Token, nil
}
