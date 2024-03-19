package settings_test

import (
	"darco/proto/db"
	"darco/proto/models/settings"
	"darco/proto/tests"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthTokenDuration(t *testing.T) {
	assert.GreaterOrEqual(t, settings.Security().AuthTokenDuration().Seconds(), 1.)
}

func TestAccountTokenDuration(t *testing.T) {
	assert.GreaterOrEqual(t, settings.Security().AccountTokenDuration().Hours(), 1.)
}

func TestSecuritySave(t *testing.T) {
	input := tests.FakeData[settings.SecuritySettingsInput](t)
	s, err := input.Save(db.Client())
	require.NoError(t, err)
	assert.Equal(t, s.SecuritySettingsInput, *input)
}
