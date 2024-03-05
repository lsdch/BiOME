package router

import (
	"darco/proto/db"
	"darco/proto/middlewares"
	"darco/proto/models/users"
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

//go:generate go run ./generate.go

type WrappedRouterGroup[T any] struct {
	*gin.RouterGroup
	wrap func(handler T) gin.HandlerFunc
}

func Wrap[T any](group *gin.RouterGroup, wrapper func(handler T) gin.HandlerFunc) *WrappedRouterGroup[T] {
	return &WrappedRouterGroup[T]{group, wrapper}
}

func (w *WrappedRouterGroup[T]) wrapHandlers(handlers []T) []gin.HandlerFunc {
	var wrapped_handlers = make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		wrapped_handlers[i] = w.wrap(handler)
	}
	return wrapped_handlers
}

// Wraps a handler that requires a DB connection to provide it as an argument.
func WithDB(handler func(*gin.Context, *edgedb.Client)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		client, ok := ctx.Get(middlewares.CTX_DATABASE_KEY)
		if !ok {
			client = db.Client()
		}
		handler(ctx, client.(*edgedb.Client))
	}
}

// Wraps a handler that requires an authenticated user to provide it as an argument.
func WithUser(handler func(*gin.Context, *edgedb.Client, *users.User)) gin.HandlerFunc {
	return WithDB(func(ctx *gin.Context, db *edgedb.Client) {
		user, ok := ctx.Get(middlewares.CTX_CURRENT_USER_KEY)
		if !ok || user == nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, "Not authenticated")
			return
		}
		handler(ctx, db, user.(*users.User))
	})
}
