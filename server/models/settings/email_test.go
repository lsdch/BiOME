package settings_test

import (
	"darco/proto/db"
	"darco/proto/models/settings"
	"darco/proto/tests"
	"testing"
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
