package people_test

import (
	"darco/proto/db"
	"darco/proto/models/people"
	"darco/proto/tests"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"
)

// Generates a fake activated user account
func FakeUserAccount(t *testing.T, role people.UserRole) *people.User {
	p := tests.FakeData[people.UserInput](t)
	person, err := FakePersonInput(t).Create(db.Client())
	require.NoError(t, err)
	user, err := p.Create(
		db.Client(),
		people.Contributor,
		person.PersonInner,
	)
	logrus.Infof("Generated test user: %+v", *user)
	require.NoError(t, err)
	return user
}

func TestFindUser(t *testing.T) {
	user := FakeUserAccount(t, people.Visitor)
	var client = db.Client()
	var assertUser = func(u people.User, err error) {
		require.NoError(t, err)
		assert.Equal(t, u.ID, user.ID)
	}

	t.Run("Find user by their UUID", func(t *testing.T) {
		u, err := people.FindID(client, user.ID)
		assertUser(u, err)
	})
	t.Run("Find user by their email", func(t *testing.T) {
		u, err := people.Find(client, user.Email)
		assertUser(u, err)
	})
	t.Run("Find user by their login", func(t *testing.T) {
		u, err := people.Find(client, user.Login)
		assertUser(u, err)
	})
	t.Run("Attempt to find non existing user", func(t *testing.T) {
		_, err := people.Find(client, "#thisuserdoesnotexist#")
		require.Error(t, err)
	})

	t.Run("Fetch current user", func(t *testing.T) {
		client := db.WithCurrentUser(user.ID)
		u, err := people.Current(client)
		assertUser(u, err)
	})

	t.Run("Fetch current user fails when global is not set", func(t *testing.T) {
		_, err := people.Current(client)
		require.Error(t, err)
	})
}

func TestDeleteUser(t *testing.T) {
	client := db.Client()
	user := FakeUserAccount(t, people.Visitor)
	deleted, err := user.Delete(client)
	require.NoError(t, err)
	assert.Equal(t, deleted, user)
	_, err = people.FindID(client, deleted.ID)
	require.Error(t, err)
}
