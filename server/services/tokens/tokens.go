package tokens

import (
	"darco/proto/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// Generates a session token for an authenticated user
func GenerateToken(payload interface{}, lifetime time.Duration) (string, error) {
	config := config.Get()

	claims := jwt.MapClaims{
		"sub": payload,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(lifetime).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("generating JWT token failed : %w", err)
	}

	return signedToken, nil
}

func ValidateToken(config *config.Config, token string) (interface{}, error) {
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		return []byte(config.SecretKey), nil
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
