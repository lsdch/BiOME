package tests

import (
	"context"
	"darco/proto/db"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func FakeData[T any](t *testing.T) *T {
	item := new(T)
	require.NoError(t, gofakeit.Struct(item))
	logrus.Debugf("Generated item %+v", item)
	return item
}

func WrapTransaction(t *testing.T, f func(tx *edgedb.Tx) error) func(t *testing.T) {
	return func(t *testing.T) {
		db.Client().Tx(context.Background(), func(ctx context.Context, tx *edgedb.Tx) error {
			err := f(tx)
			require.NoError(t, err)
			logrus.Infof("Transaction rollback")
			return errors.New("Rollback passing test")
		})
	}
}
