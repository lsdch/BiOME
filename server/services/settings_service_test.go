package services

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/stretchr/testify/require"
)

type serviceSettingsTestDB struct {
	db.Querier
	biomedb.DBTX
	key *string
	err error
}

func (d *serviceSettingsTestDB) Queries() *biomedb.Queries { return biomedb.New(d) }

func (d *serviceSettingsTestDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return d.DBTX.Query(ctx, sql, args...)
}

func (d *serviceSettingsTestDB) QueryRow(_ context.Context, _ string, args ...interface{}) pgx.Row {
	if d.err == nil && len(args) == 1 {
		d.key = args[0].(*string)
	}
	return serviceSettingsTestRow{key: d.key, err: d.err}
}

type serviceSettingsTestRow struct {
	key *string
	err error
}

func (r serviceSettingsTestRow) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	*dest[0].(*int32) = 1
	*dest[1].(**string) = r.key
	return nil
}

func TestServiceSettingsCache(t *testing.T) {
	ctx := context.Background()
	s := &SettingsService{}
	d := &serviceSettingsTestDB{}
	require.Nil(t, s.GetServiceSettings().GeoapifyAPIKey)
	oldKey, newKey := "old-key", "new-key"
	require.NoError(t, s.UpdateServiceSettings(ctx, d, models.ServiceSettings{GeoapifyAPIKey: &oldKey}))
	require.NoError(t, s.SaveServiceSettings(ctx, d, models.ServiceSettings{GeoapifyAPIKey: &newKey}))
	require.Equal(t, oldKey, *s.GetServiceSettings().GeoapifyAPIKey)
	require.NoError(t, s.ReloadServiceSettings(ctx, d))
	copy := s.GetServiceSettings()
	*copy.GeoapifyAPIKey = "modified"
	require.Equal(t, newKey, *s.GetServiceSettings().GeoapifyAPIKey)

	d.err = errors.New("database unavailable")
	require.ErrorIs(t, s.UpdateServiceSettings(ctx, d, models.ServiceSettings{}), d.err)
	require.ErrorIs(t, s.ReloadServiceSettings(ctx, d), d.err)
	require.Equal(t, newKey, *s.GetServiceSettings().GeoapifyAPIKey)
	d.err = nil
	require.NoError(t, s.UpdateServiceSettings(ctx, d, models.ServiceSettings{}))
	require.Nil(t, d.key)
	require.Nil(t, s.GetServiceSettings().GeoapifyAPIKey)
}
