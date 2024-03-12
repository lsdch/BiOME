package users_test

import (
	"darco/proto/db"
	"darco/proto/models/users"
	"errors"
	"fmt"
	"testing"
)

func TestActivationAuthentication(t *testing.T) {
	user, teardown := setupMockUser()
	defer teardown()
	db := db.Client()

	if user.Verified {
		t.Fatalf("Mock user is initially marked as active")
	}

	var credentials = users.UserCredentials{
		Identifier: "mock.user",
		Password:   "mockuserpassword",
	}

	_, auth_err := credentials.Authenticate(db)
	if auth_err == nil {
		t.Fatalf("Inactive user should not be able to authenticate.")
	} else {
		var inactive = new(users.LoginFailedError)
		if !errors.As(auth_err, &inactive) || inactive.Reason != users.AccountInactive {
			t.Fatalf("Authentication denied to inactive user but has unexpected error reason, got %v", inactive)
		}

	}

	if err := user.SetActive(db, true); err != nil {
		t.Fatalf("Failed to activate user account: %v", err)
	}

	auth_user, err := credentials.Authenticate(db)
	if err != nil {
		t.Fatalf("Failed authentication with valid credentials: %v", err)
	}
	if auth_user.ID != user.ID {
		t.Fatalf("Authentication succeeded, but does not return the expected user.")
	}

	var invalidCredentials = []users.UserCredentials{
		{"invalid.identifier", "whateverpassword"},
		{"mock.user", "invalidpassword"},
		{"mock.user@mockemail.com", "invalidpassword"},
	}
	for _, creds := range invalidCredentials {
		t.Run(
			fmt.Sprintf("Invalid creds %s:%s", creds.Identifier, creds.Password),
			func(t *testing.T) {
				u, err := creds.Authenticate(db)
				if err == nil || u != nil {
					t.Fatalf("Authentication succeeded with invalid credentials")
				}
			})
	}
}
