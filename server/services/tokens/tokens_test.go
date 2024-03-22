package tokens_test

import (
	"darco/proto/models/people"
	"darco/proto/models/settings"
	"darco/proto/services/tokens"
	"darco/proto/tests"
	"testing"

	"github.com/edgedb/edgedb-go"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"
)

func TestJWT(t *testing.T) {
	user := tests.FakeData[people.User](t)
	token, err := tokens.GenerateToken(
		user.ID,
		settings.Security().AuthTokenDuration(),
	)
	require.NoError(t, err)
	id, err := tokens.ValidateToken(token)
	require.NoError(t, err)
	uuid, err := edgedb.ParseUUID(id.(string))
	require.NoError(t, err)
	assert.Equal(t, user.ID, uuid)
}
