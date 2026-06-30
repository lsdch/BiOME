package stores

import (
	"context"

	"github.com/lsdch/biome/db/biomedb"
)

type ImportStore struct {
	q *biomedb.Queries
}

func NewImportStore(q *biomedb.Queries) *ImportStore {
	return &ImportStore{q: q}
}

func (s *ImportStore) GetImportState(ctx context.Context, importHash string) (state biomedb.ImportWorkflow, err error) {
	state, err = s.q.GetImportState(ctx, importHash)
	return state, err
}

func (s *ImportStore) InitBatchImport(ctx context.Context, importHash, label string) (biomedb.ImportWorkflow, error) {
	return s.q.InitBatchImport(ctx, importHash, label)
}

func (s *ImportStore) Bootstrap(ctx context.Context, importHash string) error {
	if err := s.q.CleanUpStagingImport(ctx, importHash); err != nil {
		return err
	}
	if err := s.q.CleanupTaxonCandidates(ctx, importHash); err != nil {
		return err
	}
	if err := s.q.CleanUpTaxonResolution(ctx, importHash); err != nil {
		return err
	}
	return s.q.CleanUpGBIFDependencies(ctx, importHash)
}

func (s *ImportStore) InsertStaging(ctx context.Context, importHash string, rows []biomedb.CopyImportStagingParams) error {
	for i := range rows {
		rows[i].ImportHash = importHash
	}
	_, err := s.q.CopyImportStaging(ctx, rows)
	return err
}
