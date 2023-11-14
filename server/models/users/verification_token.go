package users

import (
	"context"
	"darco/proto/config"
	"darco/proto/models"
	"time"

	_ "embed"

	"github.com/edgedb/edgedb-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/thanhpk/randstr"
)

type Token string

func (t *Token) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(*t), nil
}

func (t *Token) UnmarshalEdgeDBStr(data []byte) error {
	*t = Token(string(data))
	return nil
}

type AccountEmailToken struct {
	ID   edgedb.UUID `edgedb:"id"`
	User struct {
		ID edgedb.UUID `edgedb:"id"`
	} `edgedb:"user"`
	Token   Token     `edgedb:"token"`
	Expires time.Time `edgedb:"expires"`
}

func (token *AccountEmailToken) IsValid() bool {
	return token.Expires.After(time.Now())
}

func (token *AccountEmailToken) Consume() (err error) {
	deleteQuery := `delete people::AccountEmailToken filter .id = <uuid>$0`
	if err = models.DB().Execute(context.Background(), deleteQuery, token.ID); err != nil {
		logrus.Errorf("Database error %v (query: %s)", err, deleteQuery)
		return
	}
	return
}

func ValidatePasswordResetToken(token Token) (uuid.UUID, bool) {
	query := `select people::PasswordReset { id, user: {id}, token, expires }
		filter .token = <str>$0`
	var pwdReset AccountEmailToken
	if err := models.DB().QuerySingle(context.Background(), query, pwdReset, token); err != nil {
		return uuid.Nil, false
	}
	if pwdReset.IsValid() {
		pwdReset.Consume()
	}
	return uuid.UUID(pwdReset.User.ID), pwdReset.IsValid()
}

func ValidateEmailConfirmationToken(token Token) (uuid.UUID, bool) {
	query := `select people::EmailConfirmation { id, user: {id}, token, expires} filter .token = <str>$0`
	var emailConfirmation AccountEmailToken
	if err := models.DB().QuerySingle(context.Background(), query, &emailConfirmation, string(token)); err != nil {
		logrus.Infof("Failed to validate email confirmation token with error %+v", err)
		return uuid.Nil, false
	}
	logrus.Infof("Found confirmation token %+v", emailConfirmation)
	if emailConfirmation.IsValid() {
		emailConfirmation.Consume()
	}
	return uuid.UUID(emailConfirmation.User.ID), emailConfirmation.IsValid()
}

//go:embed queries/upsert_confirmation_token.edgeql
var queryCreateConfirmationToken string

func (user *User) CreateConfirmationToken() (Token, error) {
	config := config.Get()
	token := randstr.String(20)
	expires := time.Now().Add(config.Emailer.TokenLifetime)
	logrus.Infof("Creating confirmation token %s for user ID %v", token, user.ID)
	return Token(token),
		models.DB().Execute(
			context.Background(),
			queryCreateConfirmationToken,
			user.ID, token, expires,
		)
}

//go:embed queries/create_pwd_reset.edgeql
var queryCreatePasswordReset string

func (user *User) CreatePasswordResetToken() (Token, error) {
	config := config.Get()
	token := randstr.String(20)
	expiration := time.Now().Add(config.Emailer.TokenLifetime)
	return Token(token),
		models.DB().Execute(
			context.Background(),
			queryCreatePasswordReset,
			user.ID, token, expiration,
		)
}
