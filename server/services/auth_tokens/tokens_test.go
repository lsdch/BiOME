package auth_tokens_test

import (
	"testing"

	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/services/auth_tokens"
	"github.com/lsdch/biome/tests"

	"github.com/edgedb/edgedb-go"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"
)

func TestJWT(t *testing.T) {
	user := tests.FakeData[people.User](t)
	token, err := auth_tokens.NewJWT(
		user.ID,
		config.Get().AuthTokenDuration(),
	)
	require.NoError(t, err)
	id, err := auth_tokens.ValidateJWT(token)
	require.NoError(t, err)
	uuid, err := edgedb.ParseUUID(id.(string))
	require.NoError(t, err)
	assert.Equal(t, user.ID, uuid)
}
