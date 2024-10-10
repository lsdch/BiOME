package auth_tokens_test

import (
	"darco/proto/config"
	"darco/proto/models/people"
	"darco/proto/services/auth_tokens"
	"darco/proto/tests"
	"testing"

	"github.com/edgedb/edgedb-go"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"
)

func TestJWT(t *testing.T) {
	user := tests.FakeData[people.User](t)
	token, err := auth_tokens.GenerateToken(
		user.ID,
		config.Get().AuthTokenDuration(),
	)
	require.NoError(t, err)
	id, err := auth_tokens.ValidateToken(token)
	require.NoError(t, err)
	uuid, err := edgedb.ParseUUID(id.(string))
	require.NoError(t, err)
	assert.Equal(t, user.ID, uuid)
}
