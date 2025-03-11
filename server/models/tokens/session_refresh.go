package tokens

import (
	"context"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models/settings"
)

type HashedSessionRefreshToken struct {
	UserID      geltypes.UUID `gel:"user_id" json:"user_id"`
	TokenRecord `gel:"$inline" json:",inline"`
}

// Unhashed session refresh token to be shared with the client after authentication
type SessionRefreshToken struct {
	UserID      geltypes.UUID `gel:"user_id" json:"user_id"`
	TokenRecord `gel:"$inline" json:",inline"`
}

// Consumes the session refresh token and issues a new one
func (t HashedSessionRefreshToken) Rotate(db *gel.Client) (newToken SessionRefreshToken, err error) {
	err = db.Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {
		if err := t.Consume(db); err != nil {
			return err
		}
		newToken, err = CreateSessionRefreshToken(db, t.UserID)
		return err
	})
	return
}

// Generate and save a session refresh token in the database.
// Returns the unhashed refresh token.
func CreateSessionRefreshToken(db geltypes.Executor, userID geltypes.UUID) (SessionRefreshToken, error) {
	token := GenerateToken(settings.Security().RefreshTokenDuration())
	err := db.Execute(context.Background(),
		`#edgeql
			insert tokens::SessionRefreshToken {
				user := (select(<people::User><uuid>$0)),
				token := <str>$1,
				expires := <datetime>$2,
			}
		`, userID, token.Token.Hash(), token.Expires)
	return SessionRefreshToken{
		UserID:      userID,
		TokenRecord: token,
	}, err
}

// Get a refresh token in the database, provided the unhashed token string.
// Returns error if no token matches.
func RetrieveSessionRefreshToken(edb geltypes.Executor, token Token) (sessionToken HashedSessionRefreshToken, err error) {
	err = edb.QuerySingle(context.Background(),
		`#edgeql
			select tokens::SessionRefreshToken { *, user_id := .user.id }
			filter .token = <str>$0 limit 1
		`,
		&sessionToken,
		token.Hash(),
	)
	return
}

func DropSessionToken(edb geltypes.Executor, token Token) error {
	return edb.Execute(context.Background(),
		`#edgeql
			delete tokens::SessionRefreshToken filter .token = <str>0 limit 1;
		`,
		string(token))
}
