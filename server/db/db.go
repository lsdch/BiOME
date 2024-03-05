package db

import (
	"context"

	"github.com/edgedb/edgedb-go"
	log "github.com/sirupsen/logrus"
)

func ConnectDB() (db *edgedb.Client) {
	ctx := context.Background()
	db, err := edgedb.CreateClient(ctx, edgedb.Options{})

	if err != nil {
		log.Fatalf("Failed to connect to the database: %+v", err)
	}

	return
}

var db *edgedb.Client = ConnectDB()

func Client() *edgedb.Client {
	return db
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
