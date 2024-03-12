package users_test

import (
	"context"
	"darco/proto/db"
	"errors"
	"testing"

	"github.com/edgedb/edgedb-go"
)

func TestPasswordManipulation(t *testing.T) {

	user, teardown := setupMockUser()
	db := db.Client()
	defer teardown()

	t.Run("Correct password matches", func(t *testing.T) {
		if !user.PasswordMatch(db, "mockuserpassword") {
			t.Fatalf("Password %s should have matched", "mockuserpassword")
		}
	})
	t.Run("Invalid password does not match", func(t *testing.T) {
		if user.PasswordMatch(db, "invalidpassword") {
			t.Fatalf("Password %s should not match", "mockuserpassword")
		}
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
			if err != nil {
				t.Fatalf("Failed to set password: %v", err)
			}
			if !user.PasswordMatch(tx, newPwd) {
				t.Fatalf("New password does not match in database")
			}
			return errors.New("Rollback")
		})
	})

	t.Run("Weak password update fails", func(t *testing.T) {
		newPwd := "weak"
		err := user.SetPassword(db, newPwd)
		if err == nil {
			t.Fatalf("Weak password should not have been accepted")
		}
	})
}
