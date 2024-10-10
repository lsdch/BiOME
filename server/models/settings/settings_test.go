package settings_test

import (
	"context"
	"darco/proto/db"
	"darco/proto/models/settings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSettings(t *testing.T) {
	settings.Get()
}

func TestInitSettings(t *testing.T) {
	err := db.Client().Execute(context.Background(), `delete admin::Settings;`)
	require.NoError(t, err)
	settings.Get()
}
