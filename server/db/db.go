package db

import (
	"errors"
	"testing"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/gelerr"
	"github.com/geldata/gel-go/geltypes"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// Opens a new connection to Gel instance
func Connect(options gelcfg.Options) (db *gel.Client) {
	if options.Branch == "" {
		options.Branch = "main"
	}
	if testing.Testing() {
		options.Branch = "testing"
	}
	logrus.Infof("Attempting connection to database branch'%s'", options.Branch)
	db, err := gel.CreateClient(options)

	if err != nil {
		log.Fatalf("Failed to connect to database branch: %+v", err)
	}

	return
}

var db *gel.Client = Connect(gelcfg.Options{})

type DatabaseConnection string

// Gets a connection to Gel instance
func Client() *gel.Client {
	return db
}

// Get a connection to Gel instance with an authenticated user identified by an UUID
func WithCurrentUser(userID geltypes.UUID) *gel.Client {
	return db.WithGlobals(map[string]interface{}{"current_user_id": userID})
}

type NoDataError struct {
	msg string
}

func (e NoDataError) Error() string {
	return e.msg
}

func NewNoDataError(msg string) NoDataError {
	return NoDataError{msg: msg}
}

// IsNoData returns true if error is gelerr.NoDataError or custom db.NoDataError
func IsNoData(err error) bool {
	var edbErr gelerr.Error
	var customErr NoDataError
	return err != nil && (errors.As(err, &customErr) ||
		(errors.As(err, &edbErr) && edbErr.Category(gelerr.NoDataError)))
}

func IsCardinalityViolation(err error) bool {
	var edbErr gelerr.Error
	return err != nil && errors.As(err, &edbErr) && edbErr.Category(gelerr.CardinalityViolationError)
}

func IsConstraintViolation(err error) (ok bool, edbErr gelerr.Error) {
	ok = err != nil && errors.As(err, &edbErr) && edbErr.Category(gelerr.ConstraintViolationError)
	return ok, edbErr
}
