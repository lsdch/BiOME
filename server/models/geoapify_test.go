package models

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGeoapifyQuotaAndSecret(t *testing.T) {
	status := GeoapifyStatus{APIKey: "secret-key", HasApiKey: true, TodayRequests: 99, Limit: 100}
	require.NoError(t, status.AllowRequests(1))
	require.ErrorIs(t, status.AllowRequests(2), ErrLimitExceeded)
	status.HasApiKey = false
	require.ErrorIs(t, status.AllowRequests(1), ErrNoAPIKey)
	data, err := json.Marshal(status)
	require.NoError(t, err)
	require.NotContains(t, string(data), "secret-key")
}
