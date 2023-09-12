package models

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/edgedb/edgedb-go"
)

func ConnectDB() (db *edgedb.Client) {
	ctx := context.Background()
	db, err := edgedb.CreateClient(ctx, edgedb.Options{})

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	return
}

var DB *edgedb.Client = ConnectDB()
