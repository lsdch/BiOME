package people_test

import (
	"testing"

	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/tests"

	"github.com/edgedb/edgedb-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FakePendingUserInput(t *testing.T) *people.PendingUserRequestInput {
	p := tests.FakeData[people.PendingUserRequestInput](t)
	return p
}

func TestPendingUser(t *testing.T) {

	t.Run("Register user",
		tests.WrapTransaction(t, func(tx *edgedb.Tx) error {
			_, err := FakePendingUserInput(t).Register(tx)
			return err
		}))

	t.Run("Account request email is initially not verified",
		tests.WrapTransaction(t, func(tx *edgedb.Tx) error {
			pendingUser, err := FakePendingUserInput(t).Register(tx)
			require.NoError(t, err)
			assert.False(t, pendingUser.EmailVerified)
			return nil
		}))

	// t.Run("Validate pending account request", func(t *testing.T) {
	// 	client := db.Client()
	// 	pendingUser, err := FakePendingUserInput(t).Register(client)
	// 	require.NoError(t, err)
	// 	person, err := FakePersonInput(t).Save(client)
	// 	require.NoError(t, err)
	// 	role := people.Contributor
	// 	u, err := pendingUser.Validate(client, &person.PersonInner, role)
	// 	require.NoError(t, err)
	// 	assert.Equal(t, role, u.Role)
	// })

	t.Run("Delete pending user request",
		tests.WrapTransaction(t, func(tx *edgedb.Tx) error {
			pendingUser, err := FakePendingUserInput(t).Register(tx)
			require.NoError(t, err)
			return pendingUser.Delete(tx)
		}))
}
