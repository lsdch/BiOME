package tokens

import (
	"context"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/config"

	"github.com/sirupsen/logrus"
)

type pwdResetToken struct {
	TokenRecord `gel:"$inline" json:",inline"`
	UserID      geltypes.UUID `gel:"user_id" json:"user_id"`
}

func (t pwdResetToken) Save(db geltypes.Executor) error {
	return db.Execute(context.Background(),
		`#edgeql
			insert tokens::PasswordReset {
			user := (select(<people::User><uuid>$0)),
			token := <str>$1,
			expires := <datetime>$2,
		}`,
		t.UserID, t.Token, t.Expires)
}

func NewPwdResetToken(userID geltypes.UUID) pwdResetToken {
	return pwdResetToken{
		UserID:      userID,
		TokenRecord: GenerateToken(config.Get().AccountTokenDuration()),
	}
}

func RetrievePwdResetToken(db geltypes.Executor, token Token) (pwdResetToken, error) {
	var db_token pwdResetToken
	err := db.QuerySingle(context.Background(),
		`#edgeql
			select tokens::PasswordReset { token, expires, user_id:= .user.id }
			filter .token = <str>$0
		`, &db_token, token,
	)
	logrus.Debugf("%+v, %v", db_token, err)
	return db_token, err
}
