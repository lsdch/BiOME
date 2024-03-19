package people_test

import (
	"darco/proto/db"
	"darco/proto/models/people"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailConfirmationToken(t *testing.T) {
	client := db.Client()
	p, err := FakePendingUserInput(t).Register(client)
	require.NoError(t, err)
	token, err := p.User.CreateAccountToken(client, people.EmailConfirmationToken)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	user, ok := people.ValidateAccountToken(client, token, people.EmailConfirmationToken)
	assert.True(t, ok)
	if err := user.SetEmailConfirmed(client, ok); err != nil {
		require.NoError(t, err)
	}
	assert.True(t, user.EmailConfirmed)
}
