package tests

import (
	"context"
	"darco/proto/db"
	"errors"
	"testing"

	"github.com/edgedb/edgedb-go"
	"github.com/go-faker/faker/v4"
	"github.com/sirupsen/logrus"
)

func FakeData[T any](t *testing.T) *T {
	item := new(T)
	if err := faker.FakeData(item); err != nil {
		t.Fatalf("Failed to generate mock data: %v", err)
	}
	logrus.Infof("Generated item %+v", item)
	return item
}

func WrapTransaction(t *testing.T, f func(tx *edgedb.Tx) error) func(t *testing.T) {
	return func(t *testing.T) {
		db.Client().Tx(context.Background(), func(ctx context.Context, tx *edgedb.Tx) error {
			err := f(tx)
			if err != nil {
				t.Errorf("%v", err)
			}
			logrus.Infof("Transaction rollback")
			return errors.New("Rollback passing test")
		})
	}

}
