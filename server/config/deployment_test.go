package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func deploymentConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	base, err := os.ReadFile("../config.yml")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "config.yml"), base, 0600))
	for key, value := range map[string]string{
		"BIOME_ENV":                              "prod",
		"BIOME_APP_PUBLIC_BASE_URL":              "https://biome.example.org",
		"BIOME_INSTANCE_ADMIN_EMAIL":             "admin@example.org",
		"BIOME_DB_DSN":                           "postgres://user:password@database/biome",
		"BIOME_API_API_HOST":                     "0.0.0.0",
		"BIOME_API_API_PORT":                     "8080",
		"BIOME_API_API_BASE_PATH":                "/api/v1",
		"BIOME_API_VERSION":                      "1.0.0",
		"BIOME_API_API_CONTACT_NAME":             "BiOME",
		"BIOME_API_API_CONTACT_EMAIL":            "contact@example.org",
		"BIOME_AUTH_TOKENS_JWT_SECRET_KEY":       "01234567890123456789012345678901",
		"BIOME_AUTH_TOKENS_REFRESH_TOKEN_PEPPER": "test-pepper",
		"BIOME_GBIF_GBIF_USER_AGENT":             "BiOME tests",
	} {
		t.Setenv(key, value)
	}
	return dir
}

func TestDeploymentConfigFromMountedFileAndEnv(t *testing.T) {
	dir := deploymentConfig(t)
	cfg, err := LoadConfig(dir, "config")
	require.NoError(t, err)
	require.Equal(t, "postgres://user:password@database/biome", cfg.DB.DSN)
	require.Equal(t, "01234567890123456789012345678901", cfg.AuthTokens.SecretKey)
	require.Equal(t, "0.0.0.0", cfg.API.Host)
	require.Equal(t, "https://biome.example.org", cfg.AppPublicBaseURL.String())
}

func TestDeploymentConfigOverridesUseMountedDirectory(t *testing.T) {
	dir := deploymentConfig(t)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "config.prod.yml"), []byte("instance:\n  APP_NAME: Production\n  IS_PUBLIC: true\n"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "config.local.yml"), []byte("instance:\n  APP_NAME: Local override\n"), 0600))
	t.Setenv("BIOME_INSTANCE_IS_PUBLIC", "false")
	cfg, err := LoadConfig(dir, "config")
	require.NoError(t, err)
	require.Equal(t, "Local override", cfg.Instance.AppName)
	require.False(t, cfg.Instance.IsPublic)

	require.NoError(t, os.WriteFile(filepath.Join(dir, "config.local.yml"), []byte("instance: [invalid"), 0600))
	_, err = LoadConfig(dir, "config")
	require.ErrorContains(t, err, "local config")
}
