package people

import (
	"time"

	"github.com/thanhpk/randstr"
)

// A confirmation token generated as a random string.
// Used to validate email addresses or password reset requests.
type Token string

func (t Token) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(t), nil
}

func (t *Token) UnmarshalEdgeDBStr(data []byte) error {
	*t = Token(string(data))
	return nil
}

func GenerateToken(length int) Token {
	return Token(randstr.String(length))
}

type TokenWrapper struct {
	Token   Token     `edgedb:"token"`
	Expires time.Time `edgedb:"expires"`
}

func (token *TokenWrapper) IsValid() bool {
	return token.Expires.After(time.Now())
}
