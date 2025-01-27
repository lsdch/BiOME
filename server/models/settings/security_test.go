package settings_test

import (
	"testing"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/settings"
	"github.com/lsdch/biome/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecuritySave(t *testing.T) {
	input := tests.FakeData[settings.SecuritySettingsInput](t)
	s, err := input.Save(db.Client())
	require.NoError(t, err)
	assert.Equal(t, s.SecuritySettingsInput, *input)
}
