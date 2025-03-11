package tokens

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/config"

	"github.com/sirupsen/logrus"
	"github.com/thanhpk/randstr"
)

// A confirmation token generated as a random string.
// Used to validate email addresses or password reset requests.
type Token string

func (t Token) Hash() string {
	h := sha256.New()
	return hex.EncodeToString(h.Sum([]byte(t)))
}

func (t Token) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(t), nil
}

func (t *Token) UnmarshalEdgeDBStr(data []byte) error {
	*t = Token(string(data))
	return nil
}

type TokenRecord struct {
	ID      geltypes.UUID `gel:"id"`
	Token   Token         `gel:"token"`
	Expires time.Time     `gel:"expires"`
}

func (token TokenRecord) IsValid() bool {
	return token.Expires.After(time.Now())
}

// Deletes token from the database
func (token TokenRecord) Consume(db geltypes.Executor) (err error) {
	deleteQuery := `delete tokens::Token filter .id = <uuid>$0`
	if err = db.Execute(context.Background(),
		`#edgeql
			delete tokens::Token filter .id = <uuid>$0
		`,
		token.ID); err != nil {
		logrus.Errorf("Database error %v (query: %s)", err, deleteQuery)
		return
	}
	return
}

// Generates a token string with the given length
func generateTokenStr(length int) Token {
	return Token(randstr.String(length))
}

// Generates a token with expiration date attached.
// It is NOT saved in the database yet.
func GenerateToken(duration time.Duration) TokenRecord {
	return TokenRecord{
		Token:   generateTokenStr(int(config.Get().GeneratedTokenLength)),
		Expires: time.Now().Add(duration),
	}
}
