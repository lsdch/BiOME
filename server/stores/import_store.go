package stores

import (
	"context"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
)

type ImportStore struct {
}

func NewImportStore() *ImportStore {
	return &ImportStore{}
}

func (s *ImportStore) GetImportState(ctx context.Context, q db.Querier, importHash string) (state biomedb.ImportWorkflow, err error) {
	state, err = q.Queries().GetImportState(ctx, importHash)
	return state, err
}

func (s *ImportStore) InitBatchImport(ctx context.Context, q db.Querier, importHash, label string) (biomedb.ImportWorkflow, error) {
	return q.Queries().InitBatchImport(ctx, importHash, label)
}

func (s *ImportStore) Bootstrap(ctx context.Context, q db.Querier, importHash string) error {
	if err := q.Queries().CleanUpStagingImport(ctx, importHash); err != nil {
		return err
	}
	if err := q.Queries().CleanupTaxonCandidates(ctx, importHash); err != nil {
		return err
	}
	if err := q.Queries().CleanUpTaxonResolution(ctx, importHash); err != nil {
		return err
	}
	return q.Queries().CleanUpGBIFDependencies(ctx, importHash)
}

func (s *ImportStore) InsertStaging(ctx context.Context, q db.Querier, importHash string, rows []biomedb.CopyImportStagingParams) error {
	for i := range rows {
		rows[i].ImportHash = importHash
	}
	_, err := q.Queries().CopyImportStaging(ctx, rows)
	return err
}
