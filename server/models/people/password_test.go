package people_test

import (
	"context"
	"darco/proto/db"
	"errors"
	"testing"

	"github.com/edgedb/edgedb-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordManipulation(t *testing.T) {

	user := SetupUser(t)
	db := db.Client()

	t.Run("Correct password matches", func(t *testing.T) {
		input := FakePendingUserInput(t)
		u, err := input.Register(db)
		require.NoError(t, err)
		assert.True(t, u.User.PasswordMatch(db, input.User.Password))
	})
	t.Run("Invalid password does not match", func(t *testing.T) {
		assert.False(t, user.PasswordMatch(db, "invalidpassword"))
	})

	t.Run("Password update succeeds", func(t *testing.T) {
		newPwd := "somenewsufficientlystrongpassword"
		db.WithTxOptions(
			edgedb.NewTxOptions().WithIsolation(edgedb.Serializable),
		).WithRetryOptions(
			edgedb.NewRetryOptions().WithDefault(
				edgedb.NewRetryRule().WithAttempts(1),
			),
		).Tx(context.Background(), func(ctx context.Context, tx *edgedb.Tx) error {
			err := user.SetPassword(tx, newPwd)
			require.NoError(t, err)
			assert.Truef(t, user.PasswordMatch(tx, newPwd),
				"New password does not match in database")
			return errors.New("Rollback")
		})
	})

	t.Run("Weak password update fails", func(t *testing.T) {
		newPwd := "weak"
		err := user.SetPassword(db, newPwd)
		require.Error(t, err)
	})
}
