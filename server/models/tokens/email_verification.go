package tokens

import (
	"context"
	"darco/proto/config"

	"github.com/edgedb/edgedb-go"
)

type emailVerificationToken struct {
	TokenRecord `edgedb:"$inline" json:",inline"`
	Email       string `edgedb:"email" json:"email"`
}

func (t emailVerificationToken) Save(db edgedb.Executor) error {
	return db.Execute(context.Background(),
		`insert tokens::EmailVerification {
			user_request := (
				select people::PendingUserRequest filter .email = <str>$0
			),
			token := <str>$1,
			expires := <datetime>$2,
		}`, t.Email, t.Token, t.Expires)
}

func NewEmailVerificationToken(email string) emailVerificationToken {
	return emailVerificationToken{
		Email:       email,
		TokenRecord: GenerateToken(config.Get().AccountTokenDuration()),
	}
}

func RetrieveEmailToken(db edgedb.Executor, token Token) (emailVerificationToken, error) {
	var db_token emailVerificationToken
	err := db.QuerySingle(context.Background(),
		`select tokens::EmailVerification { *, email := .user_request.email } filter .token = <str>$0`,
		&db_token, token,
	)
	return db_token, err
}
