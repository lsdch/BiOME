package db

import (
	"context"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
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

// Wraps a handler that requires a DB connection to provide it as an argument.
func WithDB(handler func(*gin.Context, *edgedb.Client)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		client, ok := ctx.Get("db")
		if !ok {
			client = db
		}
		handler(ctx, client.(*edgedb.Client))
	}
}

func Client() *edgedb.Client {
	return db
}
