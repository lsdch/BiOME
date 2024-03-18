package settings_test

import (
	"darco/proto/db"
	"darco/proto/models/settings"
	"darco/proto/tests"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceSave(t *testing.T) {
	input := tests.FakeData[settings.InstanceSettingsInput](t)
	settings, err := input.Save(db.Client())
	if err != nil {
		t.Fatalf("%v", err)
	}
	if settings.Name != input.Name {
		t.Fatalf(
			"Mismatch between input settings %+v and saved instance settings %+v",
			input, settings,
		)
	}
}

func TestInstanceGet(t *testing.T) {
	assert.NotEmpty(t, settings.Instance().Name)
}
