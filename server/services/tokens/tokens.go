package tokens

import (
	"darco/proto/models/settings"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	AUTH_TOKEN_COOKIE    = "auth_token"
	REFRESH_TOKEN_COOKIE = "refresh_token"
)

func secretKey() []byte {
	return []byte(settings.Get().Security.SecretKey)
}

// Generates a session token for an authenticated user
func GenerateToken(payload interface{}, lifetime time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": payload,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(lifetime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey())
	if err != nil {
		return "", fmt.Errorf("generating JWT token failed : %w", err)
	}

	return signedToken, nil
}

func ValidateToken(token string) (interface{}, error) {
	tok, err := jwt.Parse(token,
		func(jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
			}
			return secretKey(), nil
		})
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims["sub"], nil
}
