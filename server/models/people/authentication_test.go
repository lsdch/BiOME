package people_test

import (
	"fmt"
	"testing"

	"github.com/lsdch/biome/db"
	users "github.com/lsdch/biome/models/people"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthentication(t *testing.T) {
	db := db.Client()
	user := FakeUserAccount(t, users.Contributor)

	userPwd := "test account password"
	require.NoError(t, user.SetPassword(db, userPwd))
	var credentials = users.UserCredentials{
		Identifier: user.Email,
		Password:   userPwd,
	}

	t.Run("Valid credentials", func(t *testing.T) {
		auth_user, err := credentials.Authenticate(db)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, auth_user.ID)
	})

	var invalidCredentials = []users.UserCredentials{
		{"invalid.identifier", "whateverpassword"},
		{"mock.user", "invalidpassword"},
		{"mock.user@mockemail.com", "invalidpassword"},
	}
	for _, creds := range invalidCredentials {
		t.Run(
			fmt.Sprintf("Invalid creds %s:%s", creds.Identifier, creds.Password),
			func(t *testing.T) {
				_, err := creds.Authenticate(db)
				require.Error(t, err)
			})
	}
}
