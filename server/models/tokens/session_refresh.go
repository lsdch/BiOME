package tokens

import (
	"context"
	"darco/proto/models/settings"

	"github.com/edgedb/edgedb-go"
)

type HashedSessionRefreshToken struct {
	UserID      edgedb.UUID `edgedb:"user_id" json:"user_id"`
	TokenRecord `edgedb:"$inline" json:",inline"`
}

// Unhashed session refresh token to be shared with the client after authentication
type SessionRefreshToken struct {
	UserID      edgedb.UUID `edgedb:"user_id" json:"user_id"`
	TokenRecord `edgedb:"$inline" json:",inline"`
}

// Consumes the session refresh token and issues a new one
func (t HashedSessionRefreshToken) Rotate(db *edgedb.Client) (newToken SessionRefreshToken, err error) {
	err = db.Tx(context.Background(), func(ctx context.Context, tx *edgedb.Tx) error {
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
func CreateSessionRefreshToken(db edgedb.Executor, userID edgedb.UUID) (SessionRefreshToken, error) {
	token := GenerateToken(settings.Security().RefreshTokenDuration())
	err := db.Execute(context.Background(),
		`insert tokens::SessionRefreshToken {
			user := (select(<people::User><uuid>$0)),
			token := <str>$1,
			expires := <datetime>$2,
		}`,
		userID, token.Token.Hash(), token.Expires)
	return SessionRefreshToken{
		UserID:      userID,
		TokenRecord: token,
	}, err
}

// Get a refresh token in the database, provided the unhashed token string.
// Returns error if no token matches.
func RetrieveSessionRefreshToken(edb edgedb.Executor, token Token) (sessionToken HashedSessionRefreshToken, err error) {
	err = edb.QuerySingle(context.Background(),
		`select tokens::SessionRefreshToken { *, user_id := .user.id }
			filter .token = <str>$0 limit 1`,
		&sessionToken,
		token.Hash(),
	)
	return
}

func DropSessionToken(edb edgedb.Executor, token Token) error {
	return edb.Execute(context.Background(),
		`delete tokens::SessionRefreshToken filter .token = <str>0 limit 1;`,
		string(token))
}
