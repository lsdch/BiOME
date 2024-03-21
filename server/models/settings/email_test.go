package settings_test

import (
	"darco/proto/db"
	"darco/proto/models/settings"
	"darco/proto/tests"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestEmailSettingsSave(t *testing.T) {
	input := tests.FakeData[settings.EmailSettingsInput](t)
	settings, err := input.Save(db.Client())
	if err != nil {
		t.Fatalf("%v", err)
	}
	if settings.EmailSettingsInput != *input {
		t.Fatalf(
			"Saved email settings %+v do not match input %+v",
			settings, input,
		)
	}
}

func TestEmailWriteYAML(t *testing.T) {
	input := tests.FakeData[settings.EmailSettingsInput](t)
	dir := t.TempDir()
	output := dir + "/email.yml"
	err := input.WriteYAML(output)
	require.NoError(t, err)
	content, err := os.ReadFile(output)
	require.NoError(t, err)
	var yamlSettings settings.EmailSettingsInput
	err = yaml.Unmarshal(content, &yamlSettings)
	require.NoError(t, err)
	assert.Equal(t, *input, yamlSettings)
}
