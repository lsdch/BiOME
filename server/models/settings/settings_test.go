package settings_test

import (
	"context"
	"darco/proto/db"
	"darco/proto/models/settings"
	"testing"
)

func TestGetSettings(t *testing.T) {
	settings.Get()
}

func TestInitSettings(t *testing.T) {
	db.Client().Execute(context.Background(), `delete admin::Settings;`)
	settings.Get()
}
