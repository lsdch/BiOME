package people_test

import (
	"darco/proto/db"
	users "darco/proto/models/people"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"
)

func TestFindUser(t *testing.T) {
	user := SetupUser(t)
	var client = db.Client()
	var assertUser = func(u *users.User, err error) {
		require.NoError(t, err)
		assert.Equal(t, u.ID, user.ID)
	}

	t.Run("Find user by their UUID", func(t *testing.T) {
		u, err := users.FindID(client, user.ID)
		assertUser(u, err)
	})
	t.Run("Find user by their email", func(t *testing.T) {
		u, err := users.Find(client, user.Email)
		assertUser(u, err)
	})
	t.Run("Find user by their login", func(t *testing.T) {
		u, err := users.Find(client, user.Login)
		assertUser(u, err)
	})
	t.Run("Attempt to find non existing user", func(t *testing.T) {
		_, err := users.Find(client, "#thisuserdoesnotexist#")
		require.Error(t, err)
	})

	t.Run("Fetch current user", func(t *testing.T) {
		client := db.WithCurrentUser(user.ID)
		u, err := users.Current(client)
		assertUser(u, err)
	})

	t.Run("Fetch current user fails when global is not set", func(t *testing.T) {
		u, err := users.Current(client)
		require.Error(t, err)
		assert.Nil(t, u)
	})
}

func TestDeleteUser(t *testing.T) {
	client := db.Client()
	user := SetupUser(t)
	deleted, err := user.Delete(client)
	require.NoError(t, err)
	assert.Equal(t, deleted, user)
	_, err = users.FindID(client, deleted.ID)
	require.Error(t, err)
}
