package auth

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/services"
)

type ctxKey int

const sessionCtxKey ctxKey = iota

func WithSession(ctx huma.Context, session *services.SessionContext) huma.Context {
	return huma.WithValue(ctx, sessionCtxKey, session)
}

func SessionFromContext(ctx context.Context) (*services.SessionContext, bool) {
	user, ok := ctx.Value(sessionCtxKey).(*services.SessionContext)
	return user, ok
}
