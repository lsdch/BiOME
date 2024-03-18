package people_test

import (
	"darco/proto/models/people"
	"darco/proto/tests"
	"fmt"
	"testing"

	"github.com/edgedb/edgedb-go"
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

	t.Run("Validate pending account request",
		tests.WrapTransaction(t, func(tx *edgedb.Tx) error {
			pendingUser, err := FakePendingUserInput(t).Register(tx)
			if err != nil {
				return err
			}
			person, err := FakePersonInput(t).Create(tx)
			if err != nil {
				return err
			}
			_, err = pendingUser.ValidateTx(tx, &person)
			if err != nil {
				return err
			}
			return nil
		}))

	t.Run("Delete pending user request",
		tests.WrapTransaction(t, func(tx *edgedb.Tx) error {
			pendingUser, err := FakePendingUserInput(t).Register(tx)
			if err != nil {
				return err
			}
			return pendingUser.Delete(tx)
		}))
}
