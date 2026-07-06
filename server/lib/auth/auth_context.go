package auth

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/models"
)

type ctxKey int

const sessionCtxKey ctxKey = iota

func WithSession(ctx huma.Context, session *models.SessionContext) huma.Context {
	return huma.WithValue(ctx, sessionCtxKey, session)
}

func SessionFromContext(ctx context.Context) (*models.SessionContext, bool) {
	session, ok := ctx.Value(sessionCtxKey).(*models.SessionContext)
	return session, ok
}
