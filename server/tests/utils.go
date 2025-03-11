package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func FakeData[T any](t *testing.T) *T {
	item := new(T)
	require.NoError(t, gofakeit.Struct(item))
	logrus.Debugf("Generated item %+v", item)
	return item
}

func WrapTransaction(t *testing.T, f func(tx geltypes.Tx) error) func(t *testing.T) {
	return func(t *testing.T) {
		db.Client().Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {
			err := f(tx)
			require.NoError(t, err)
			logrus.Infof("Transaction rollback")
			return errors.New("Rollback passing test")
		})
	}
}
