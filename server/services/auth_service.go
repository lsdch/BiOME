package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/app_errors"
	"github.com/lsdch/biome/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidRefreshToken = app_errors.ForbiddenError(fmt.Errorf("invalid refresh token"))
	ErrRefreshTokenReuse   = app_errors.ForbiddenError(fmt.Errorf("refresh token reuse detected"))
)

type AuthService struct {
	config config.AuthTokensConfig
}

func NewAuthService(config config.AuthTokensConfig) *AuthService {
	return &AuthService{config: config}
}

func (s *AuthService) RefreshTokenExpiresAt() time.Time {
	return time.Now().Add(s.config.RefreshTokenLifetime)
}

func (s *AuthService) Login(ctx context.Context, db db.Querier, credentials models.UserCredentials, origin models.AuthOrigin) (models.LoginResult, error) {
	user, err := s.AuthenticateCredentials(ctx, db, credentials)
	if err != nil {
		return models.LoginResult{}, err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return models.LoginResult{}, err
	}

	sessionID := uuid.New()
	session, err := db.Queries().CreateSession(ctx, biomedb.CreateSessionParams{
		SessionID:        sessionID,
		UserID:           user.ID,
		RefreshTokenHash: s.HashRefreshToken(refreshToken),
		ExpiresAt:        s.RefreshTokenExpiresAt(),
		UserAgent:        origin.UserAgent.ToPtr(),
		IPAddress:        origin.IP,
	})
	if err != nil {
		return models.LoginResult{}, err
	}

	authToken, err := s.GenerateJWT(user.ID.String(), session.ID.String())
	if err != nil {
		return models.LoginResult{}, err
	}

	return models.LoginResult{
		User: models.UserFromDB(user),
		Session: models.SessionTokens{
			SessionID:           session.ID,
			AuthToken:           authToken,
			AuthTokenExpiration: time.Now().Add(s.config.AuthTokenLifetime),
			RefreshToken:        refreshToken,
		},
	}, nil
}

func (s *AuthService) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (s *AuthService) AuthenticateCredentials(ctx context.Context, q db.Querier, credentials models.UserCredentials) (biomedb.User, error) {
	user, err := q.Queries().GetUserByLoginOrEmail(ctx, credentials.Identifier)
	if err != nil {
		return biomedb.User{}, app_errors.ForbiddenError(fmt.Errorf("invalid credentials"))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password))
	if err != nil {
		return biomedb.User{}, app_errors.ForbiddenError(fmt.Errorf("invalid credentials"))
	}

	return user, nil
}

func (s *AuthService) AuthenticateJWT(ctx context.Context, q db.Querier, token string) (*models.SessionContext, error) {
	auth_ctx, err := s.ValidateJWT(token)
	if err != nil {
		return nil, err
	}

	u, err := q.Queries().GetUserByID(ctx, auth_ctx.UserID)
	if err != nil {
		return nil, err
	}
	user := models.UserFromDB(u)
	return &models.SessionContext{
		AuthContext: *auth_ctx,
		User:        user,
	}, nil
}

func (s *AuthService) GenerateJWT(userID string, sessionID string) (string, error) {
	claims := JWTClaims{
		Sub: userID,
		Sid: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AuthTokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("generating JWT token failed : %w", err)
	}

	return signedToken, nil
}

type JWTClaims struct {
	Sub string `json:"sub"`
	Sid string `json:"sid"`
	jwt.RegisteredClaims
}

func (s *AuthService) ValidateJWT(token string) (*models.AuthContext, error) {
	claims := &JWTClaims{}

	tok, err := jwt.ParseWithClaims(
		token,
		claims,
		func(jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, app_errors.ForbiddenError(fmt.Errorf("unexpected signing method: %v", jwtToken.Header["alg"]))
			}
			return []byte(s.config.SecretKey), nil
		},
	)

	if err != nil {
		return nil, app_errors.ForbiddenError(fmt.Errorf("failed to parse token: %w", err))
	}

	if !tok.Valid {
		return nil, app_errors.ForbiddenError(fmt.Errorf("invalid token"))
	}

	userID, err := uuid.Parse(claims.Sub)
	if err != nil {
		return nil, app_errors.ForbiddenError(fmt.Errorf("invalid user ID in token: %w", err))
	}

	sessionID, err := uuid.Parse(claims.Sid)
	if err != nil {
		return nil, app_errors.ForbiddenError(fmt.Errorf("invalid session ID in token: %w", err))
	}

	return &models.AuthContext{
		UserID:    userID,
		SessionID: sessionID,
	}, nil
}

func (s *AuthService) RefreshSession(ctx context.Context, q db.Querier, sessionID uuid.UUID, refreshToken string) (models.SessionTokens, error) {

	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		return models.SessionTokens{}, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	oldHash := s.HashRefreshToken(refreshToken)

	session, err := q.Queries().RotateSessionRefreshToken(ctx, biomedb.RotateSessionRefreshTokenParams{
		SessionID:           sessionID,
		NewRefreshTokenHash: s.HashRefreshToken(newRefreshToken),
		OldRefreshTokenHash: oldHash,
	})

	// Successful rotation, return new tokens
	if err == nil {
		token, err := s.GenerateJWT(session.UserID.String(), session.ID.String())
		if err != nil {
			return models.SessionTokens{}, err
		}

		return models.SessionTokens{
			SessionID:           session.ID,
			AuthToken:           token,
			AuthTokenExpiration: time.Now().Add(s.config.AuthTokenLifetime),
			RefreshToken:        newRefreshToken,
		}, nil
	}

	// Check token reuse
	sessionCheck, err2 := q.Queries().GetSession(ctx, sessionID)
	if err2 != nil {
		return models.SessionTokens{}, ErrInvalidRefreshToken
	}

	// session still active : session compromised, revoke it and return error
	if sessionCheck.RevokedAt.Valid == false {
		_ = s.Logout(ctx, q, sessionID)
		return models.SessionTokens{}, ErrRefreshTokenReuse
	}

	return models.SessionTokens{}, ErrInvalidRefreshToken
}

func (s *AuthService) Logout(ctx context.Context, q db.Querier, sessionID uuid.UUID) error {
	_, err := q.Queries().RevokeSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}
	return nil
}

func (s *AuthService) RevokeAllSessions(ctx context.Context, q db.Querier, userID uuid.UUID) error {
	err := q.Queries().RevokeAllSessionsForUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke sessions: %w", err)
	}
	return nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32) // 256 bits
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *AuthService) HashRefreshToken(token string) string {
	h := sha256.New()
	h.Write([]byte("refresh_token:v1:"))
	h.Write([]byte(token))
	h.Write([]byte(s.config.RefreshTokenPepper))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func (s *AuthService) SessionCookie(authToken string) http.Cookie {
	return http.Cookie{
		Name:     s.config.AuthTokenCookieName,
		Value:    authToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(s.config.AuthTokenLifetime.Seconds()),
		Expires:  time.Now().Add(s.config.AuthTokenLifetime),
	}
}

func (s *AuthService) DropSessionCookie() http.Cookie {
	return http.Cookie{
		Name:     s.config.AuthTokenCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}
}
