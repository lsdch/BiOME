package middleware

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/services"
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
	authService services.AuthService
}

func NewAuthMiddleware(api huma.API, as services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: as, api: api}
}

// Authentication middleware that validates the JWT token from the Authorization header or session cookie.
func (m *AuthMiddleware) AuthN(ctx huma.Context, next func(huma.Context)) {

	authHeader := ctx.Header("Authorization")
	token := extractBearer(authHeader)
	if token == "" {
		cookie, err := huma.ReadCookie(ctx, "auth_token")
		if err == nil {
			token = cookie.Value
		}
	}

	if token == "" {
		huma.WriteErr(m.api, ctx, 401, "Invalid authentication token")
		return
	}

	user, err := m.authService.AuthenticateJWT(ctx.Context(), token)
	if err != nil {
		huma.WriteErr(m.api, ctx, 401, "Invalid authentication token")
		return
	}

	ctx = auth.WithSession(ctx, user)

	next(ctx)
}

func extractBearer(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

// Authentication middleware that validates the JWT token from the Authorization header or session cookie.
func (m *AuthMiddleware) AuthZ(ctx huma.Context, next func(huma.Context)) {

	session, ok := auth.SessionFromContext(ctx.Context())
	if !ok || session == nil {
		huma.WriteErr(m.api, ctx, 401, "authentication required")
		return
	}

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

	if !policy.Eval(&session.User) {
		huma.WriteErr(m.api, ctx, 403, "access denied: insufficient permissions")
		return
	}

	next(ctx)
}
