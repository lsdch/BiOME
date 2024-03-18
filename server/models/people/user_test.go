package people_test

import (
	"context"
	"darco/proto/db"
	users "darco/proto/models/people"
	"testing"

	"github.com/sirupsen/logrus"
)

func setupMockUser() (*users.User, func()) {
	query := `
		with module people select (
			insert User {
				login := "mock.user",
				email := "mock.user@mockemail.com",
				password := "mockuserpassword",
				role := UserRole.Visitor,
				identity := (insert Person {
					first_name := "Mock",
					last_name := "User",
					contact := "mock.user@mockemail.com"
				})
			} #unless conflict on .email else (select User)
		) { *, identity: { * }}`

	var user users.User
	err := db.Client().QuerySingle(context.Background(), query, &user)
	if err != nil {
		logrus.Fatalf("Failed to setup test: %v", err)
	}
	var teardown = func() {
		db.Client().Execute(context.Background(),
			`delete people::User filter .id = <uuid>$0 limit 1`, user.ID,
		)
	}
	return &user, teardown
}

func TestFindUser(t *testing.T) {
	user, teardown := setupMockUser()
	defer teardown()
	var client = db.Client()
	var assertUser = func(u *users.User, err error) {
		if err != nil {
			t.Fatalf("Failed retrieving mock user")
		}
		if u.ID != user.ID {
			t.Fatalf("Retrieved user is not the expected one")
		}
	}

	t.Run("Find user by their UUID", func(t *testing.T) {
		u, err := users.FindID(client, user.ID)
		assertUser(u, err)
	})
	t.Run("Find user by their email", func(t *testing.T) {
		u, err := users.Find(client, "mock.user@mockemail.com")
		assertUser(u, err)
	})
	t.Run("Find user by their login", func(t *testing.T) {
		u, err := users.Find(client, "mock.user")
		assertUser(u, err)
	})
	t.Run("Attempt to find non existing user", func(t *testing.T) {
		_, err := users.Find(client, "#thisuserdoesnotexist#")
		if err == nil {
			t.Fatalf("Looking up for non existing user did not return an error")
		}
	})

	t.Run("Fetch current user", func(t *testing.T) {
		client := db.WithCurrentUser(user.ID)
		u, err := users.Current(client)
		assertUser(u, err)
	})

	t.Run("Fetch current user fails when global is not set", func(t *testing.T) {
		u, err := users.Current(client)
		if err == nil || u != nil {
			t.Fatalf("Should return error when EdgeDB global current_user_id is not set")
		}
	})
}
