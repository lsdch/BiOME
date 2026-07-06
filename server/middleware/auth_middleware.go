package middleware

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/services"
	"github.com/sirupsen/logrus"
)

const ACCESS_POLICY_EXTENSION = "policy"

func WithAccessPolicy(op huma.Operation, policy auth.Policy) huma.Operation {
	if op.Extensions == nil {
		op.Extensions = map[string]any{}
	}
	op.Extensions[ACCESS_POLICY_EXTENSION] = policy
	return op
}

func GetAccessPolicy(op huma.Operation) (auth.Policy, bool) {
	rawPolicy, ok := op.Extensions[ACCESS_POLICY_EXTENSION]
	if !ok {
		return nil, false
	}
	policy, ok := rawPolicy.(auth.Policy)
	return policy, ok
}

type AuthMiddleware struct {
	api         huma.API
	db          *db.DB
	authService *services.AuthService
	config      config.AuthTokensConfig
}

func NewAuthMiddleware(api huma.API, db *db.DB, cfg config.AuthTokensConfig, as *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: as, api: api, db: db, config: cfg}
}

// Authentication middleware that validates the JWT token from the Authorization header or session cookie.
func (m *AuthMiddleware) AuthN(ctx huma.Context, next func(huma.Context)) {

	authHeader := ctx.Header("Authorization")
	token := extractBearer(authHeader)
	if token == "" {
		logrus.Infof("Reading cookie for authentication token: %s", m.config.AuthTokenCookieName)
		cookie, err := huma.ReadCookie(ctx, m.config.AuthTokenCookieName)
		if err == nil {
			token = cookie.Value
		} else {
			logrus.Warnf("Failed to read cookie for authentication token: %v", err)
		}
	}

	if token == "" {
		next(ctx)
		return
	}

	session, err := m.authService.AuthenticateJWT(ctx.Context(), m.db, token)
	if err != nil {
		logrus.Warnf("Failed to authenticate JWT token: %v", err)
		next(ctx)
		return
	}

	logrus.Debugf("Authenticated user: %s", session.User.Login)
	ctx = auth.WithSession(ctx, session)

	next(ctx)
}

func extractBearer(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

// Authorization middleware that checks if the authenticated user has the required permissions to access the endpoint.
func (m *AuthMiddleware) AuthZ(ctx huma.Context, next func(huma.Context)) {

	op := ctx.Operation()
	rawPolicy := op.Extensions["policy"]
	if rawPolicy == nil {
		next(ctx)
		return
	}

	policy, ok := rawPolicy.(auth.Policy)
	if !ok {
		huma.WriteErr(m.api, ctx, 500, "invalid policy type")
		return
	}

	if policy.IsAuthRequired() {
		session, ok := auth.SessionFromContext(ctx.Context())
		if !ok || session == nil {
			huma.WriteErr(m.api, ctx, 401, "authentication required")
			return
		}
		if !policy.Eval(&session.User) {
			huma.WriteErr(m.api, ctx, 403, "access denied: insufficient permissions")
			return
		}
	}

	next(ctx)
}
