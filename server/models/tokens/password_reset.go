package tokens

import (
	"context"

	"github.com/edgedb/edgedb-go"
)

type pwdResetToken struct {
	TokenRecord `edgedb:"$inline" json:",inline"`
	UserID      edgedb.UUID `edgedb:"user_id" json:"user_id"`
}

func (t pwdResetToken) Save(db edgedb.Executor) error {
	return db.Execute(context.Background(),
		`insert tokens::PasswordReset {
			user := (select(<people::User><uuid>$0)),
			token := <str>$1,
			expires := <str>$2,
		}`, t.UserID, t.Token, t.Expires)
}

func NewPwdResetToken(userID edgedb.UUID) pwdResetToken {
	return pwdResetToken{
		UserID:      userID,
		TokenRecord: GenerateToken(),
	}
}

func RetrievePwdResetToken(db edgedb.Executor, token Token) (pwdResetToken, error) {
	var db_token pwdResetToken
	err := db.QuerySingle(context.Background(),
		`select tokens::PasswordReset { token, expires, user_id: user.id} } filter .token = <str>$0`,
		&db_token, token,
	)
	return db_token, err
}
