package db

import (
	"context"

	"github.com/edgedb/edgedb-go"
	log "github.com/sirupsen/logrus"
)

// A common interface between *edgedb.Client and *edgedb.Tx
type Executor interface {
	Execute(context.Context, string, ...any) error
	Query(context.Context, string, any, ...any) error
	QueryJSON(context.Context, string, *[]byte, ...any) error
	QuerySingle(context.Context, string, any, ...any) error
	QuerySingleJSON(context.Context, string, any, ...any) error
}

// Opens a new connection to EdgeDB
func connectDB() (db *edgedb.Client) {
	ctx := context.Background()
	db, err := edgedb.CreateClient(ctx, edgedb.Options{})

	if err != nil {
		log.Fatalf("Failed to connect to the database: %+v", err)
	}

	return
}

var db *edgedb.Client = connectDB()

// Gets a connection to EdgeDB
func Client() *edgedb.Client {
	return db
}

// Get a connection to EdgeDB with an authenticated user identified by an UUID
func WithCurrentUser(userID edgedb.UUID) *edgedb.Client {
	return db.WithGlobals(map[string]interface{}{"current_user_id": userID})
}

type Optional[T any] interface {
	Get() (T, bool)
}

func OptionalAsPointer[T any](opt Optional[T]) *T {
	val, ok := opt.Get()
	if ok {
		return &val
	}
	return nil
}
