package people_test

import (
	"testing"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/people"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordManipulation(t *testing.T) {

	user := FakeUserAccount(t, people.Visitor)
	db := db.Client()

	t.Run("Correct password matches", func(t *testing.T) {
		pwd := "the user password to test"
		user.SetPassword(db, pwd)
		assert.True(t, user.PasswordMatch(db, pwd))
	})
	t.Run("Invalid password does not match", func(t *testing.T) {
		assert.False(t, user.PasswordMatch(db, "invalidpassword"))
	})

	t.Run("Password update succeeds", func(t *testing.T) {
		newPwd := "somenewsufficientlystrongpassword"
		err := user.SetPassword(db, newPwd)
		require.NoError(t, err)
		assert.Truef(t, user.PasswordMatch(db, newPwd),
			"New password does not match in database")
	})
}
