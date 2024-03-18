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
	if settings.Security().AuthTokenDuration().Seconds() < 1 {
		t.Fatalf("Duration should be more than 1 second")
	}
}

func TestAccountTokenDuration(t *testing.T) {
	if settings.Security().AccountTokenDuration().Hours() < 1 {
		t.Fatalf("Duration should be more than 1 hour")
	}
}

func TestSecuritySave(t *testing.T) {
	input := tests.FakeData[settings.SecuritySettingsInput](t)
	s, err := input.Save(db.Client())
	require.NoError(t, err)
	assert.Equal(t, s.SecuritySettingsInput, *input)
}
