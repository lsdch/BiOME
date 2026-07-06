package db

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/gelerr"
	"github.com/geldata/gel-go/geltypes"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

const CTX_USER_KEY string = "current_user_id"

type Querier interface {
	Queries() *biomedb.Queries
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	EnsureTx(ctx context.Context) (txdb *Tx, err error)
	WithTx(ctx context.Context, fn func(tx *Tx) error) error
}

type DB struct {
	pool    *pgxpool.Pool
	queries *biomedb.Queries
}

func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		pool:    pool,
		queries: biomedb.New(pool),
	}
}

func (d *DB) Close() {
	d.pool.Close()
}

func (d *DB) Queries() *biomedb.Queries {
	return d.queries
}

func (d *DB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return d.pool.Query(ctx, sql, args...)
}

func (d *DB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return d.pool.QueryRow(ctx, sql, args...)
}

func (d *DB) EnsureTx(ctx context.Context) (txdb *Tx, err error) {
	tx, err := d.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	q := biomedb.New(tx)
	userID, ok := ctx.Value(CTX_USER_KEY).(uuid.UUID)

	if ok {
		err := q.SetCurrentUser(ctx, userID.String())
		if err != nil {
			return nil, fmt.Errorf("failed to set current user with ID %s", userID)
		}
	}

	return &Tx{tx: tx, queries: q, UserID: userID}, nil
}

func (d *DB) WithTx(ctx context.Context, fn func(tx *Tx) error) error {
	tx, err := d.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

type Tx struct {
	tx      pgx.Tx
	queries *biomedb.Queries
	UserID  uuid.UUID
}

func (t *Tx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *Tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func (t *Tx) Queries() *biomedb.Queries {
	return t.queries
}

func (t *Tx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return t.tx.Query(ctx, sql, args...)
}

func (t *Tx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return t.tx.QueryRow(ctx, sql, args...)
}

func (t *Tx) EnsureTx(ctx context.Context) (txdb *Tx, err error) {
	return t, nil
}

func (t *Tx) WithTx(ctx context.Context, fn func(q *Tx) error) error {
	err := fn(t)
	if err != nil {
		return err
	}

	return nil
}

// Opens a new connection to Gel instance
func Connect(options gelcfg.Options) (db *gel.Client) {
	if options.Branch == "" {
		options.Branch = "main"
	}
	if testing.Testing() {
		options.Branch = "testing"
	}
	logrus.Debugf("Attempting connection to database branch '%s'", options.Branch)
	db, err := gel.CreateClient(options)

	if err != nil {
		logrus.Fatalf("Failed to connect to database branch: %+v", err)
	}
	logrus.Debugf("Connected to database branch '%s'", options.Branch)

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

func WithBatchMode(db *gel.Client, batchULID ulid.ULID) *gel.Client {
	return db.WithGlobals(map[string]interface{}{"batch_import_id": batchULID.String()})
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
