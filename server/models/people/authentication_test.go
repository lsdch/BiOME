package people_test

import (
	"darco/proto/db"
	users "darco/proto/models/people"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivationAuthentication(t *testing.T) {
	user, teardown := setupMockUser()
	defer teardown()
	db := db.Client()

	assert.Falsef(t, user.EmailConfirmed, "New user email must not be already confirmed when first created.")

	var credentials = users.UserCredentials{
		Identifier: "mock.user",
		Password:   "mockuserpassword",
	}

	_, auth_err := credentials.Authenticate(db)
	require.Errorf(t, auth_err, "Inactive user should not be able to authenticate.")
	assert.Equal(t, auth_err.Reason, users.AccountInactive)

	// if err := user.SetActive(db, true); err != nil {
	// 	t.Fatalf("Failed to activate user account: %v", err)
	// }

	require.NoError(t, user.SetEmailConfirmed(db, true))
	person, err := FakePersonInput(t).Create(db)
	require.NoError(t, err)
	require.NoError(t, user.SetIdentity(db, &person))
	auth_user, err := credentials.Authenticate(db)
	assert.Nil(t, err)
	assert.Equal(t, user, auth_user)

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
				require.Error(t, err)
				assert.Nil(t, u)
			})
	}
}
