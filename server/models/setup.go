package models

import (
	"context"

	"github.com/edgedb/edgedb-go"
)

func ConnectDB() (db *edgedb.Client) {
	ctx := context.Background()
	db, err := edgedb.CreateClient(ctx, edgedb.Options{})

	if err != nil {
		panic("Failed to connect to database.")
	}

	return
}

var DB *edgedb.Client = ConnectDB()
