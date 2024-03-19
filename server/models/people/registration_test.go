package people_test

import (
	"darco/proto/db"
	"darco/proto/models/people"
	"darco/proto/tests"
	"fmt"
	"testing"

	"github.com/edgedb/edgedb-go"
	"github.com/stretchr/testify/require"
)

func FakePendingUserInput(t *testing.T) *people.PendingUserRequestInput {
	p := tests.FakeData[people.PendingUserRequestInput](t)
	p.User.ConfirmPwd = p.User.Password
	return p
}

func TestPendingUser(t *testing.T) {

	t.Run("Register user",
		tests.WrapTransaction(t, func(tx *edgedb.Tx) error {
			_, err := FakePendingUserInput(t).Register(tx)
			return err
		}))

	t.Run("New pending user is initially inactive",
		tests.WrapTransaction(t, func(tx *edgedb.Tx) error {
			pendingUser, err := FakePendingUserInput(t).Register(tx)
			if err != nil {
				return err
			}
			if pendingUser.User.IsActive {
				return fmt.Errorf("User is active")
			}
			return nil
		}))

	t.Run("Validate pending account request", func(t *testing.T) {
		client := db.Client()
		pendingUser, err := FakePendingUserInput(t).Register(client)
		require.NoError(t, err)
		person, err := FakePersonInput(t).Create(client)
		require.NoError(t, err)
		_, err = pendingUser.Validate(client, &person)
		require.NoError(t, err)
	})

	t.Run("Delete pending user request",
		tests.WrapTransaction(t, func(tx *edgedb.Tx) error {
			pendingUser, err := FakePendingUserInput(t).Register(tx)
			if err != nil {
				return err
			}
			return pendingUser.Delete(tx)
		}))
}
