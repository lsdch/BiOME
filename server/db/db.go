package db

import (
	"context"
	"errors"
	"testing"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// Opens a new connection to EdgeDB
func Connect(options edgedb.Options) (db *edgedb.Client) {
	ctx := context.Background()
	if options.Database == "" {
		options.Database = "main"
	}
	if testing.Testing() {
		options.Database = "testing"
	}
	logrus.Infof("Attempting connection to database '%s'", options.Database)
	db, err := edgedb.CreateClient(ctx, options)

	if err != nil {
		log.Fatalf("Failed to connect to the database: %+v", err)
	}

	return
}

var db *edgedb.Client = Connect(edgedb.Options{})

type DatabaseConnection string

// Gets a connection to EdgeDB
func Client() *edgedb.Client {
	return db
}

// Get a connection to EdgeDB with an authenticated user identified by an UUID
func WithCurrentUser(userID edgedb.UUID) *edgedb.Client {
	return db.WithGlobals(map[string]interface{}{"current_user_id": userID})
}

// IsNoData returns true if error is edgedb.NoDataError
func IsNoData(err error) bool {
	var edbErr edgedb.Error
	return errors.As(err, &edbErr) && edbErr.Category(edgedb.NoDataError)
}
