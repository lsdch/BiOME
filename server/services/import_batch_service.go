package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/oklog/ulid/v2"
)

type ImportBatchService struct{}

func NewImportBatchService() *ImportBatchService {
	return &ImportBatchService{}
}

func (s *ImportBatchService) GetImportBatch(ctx context.Context, q db.Querier, id ulid.ULID) (*models.ImportBatch, error) {
	ib, err := q.Queries().GetImportBatch(ctx, id)
	if err != nil {
		return nil, err
	}
	importBatch := models.ImportBatchFromDB(ib)
	return &importBatch, nil
}

func (s *ImportBatchService) GetImportBatchForOccurrence(ctx context.Context, q db.Querier, occurrenceID ulid.ULID) (models.Optional[models.ImportBatch], error) {
	ib, err := q.Queries().GetImportBatchForOccurrence(ctx, occurrenceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Optional[models.ImportBatch]{}, nil
		}
		return models.Optional[models.ImportBatch]{}, err
	}
	importBatch := models.ImportBatchFromDB(ib)
	return models.NewOptional(importBatch), nil
}

func (s *ImportBatchService) ListImportBatches(ctx context.Context, q db.Querier) ([]models.ImportBatch, error) {
	ibs, err := q.Queries().ListImportBatches(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]models.ImportBatch, len(ibs))
	for i, ib := range ibs {
		result[i] = models.ImportBatchFromDB(ib)
	}
	return result, nil
}

func (s *ImportBatchService) DeleteImportBatch(ctx context.Context, q db.Querier, id ulid.ULID) error {
	err := q.Queries().DeleteImportBatch(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ImportBatchService) DeleteImportBatchWithOccurrences(ctx context.Context, tx *db.Tx, id ulid.ULID) error {
	q := tx.Queries()
	if err := q.DeleteOccurrencesFromBatch(ctx, id); err != nil {
		return err
	}
	if err := q.DeleteImportBatch(ctx, id); err != nil {
		return err
	}
	return nil
}
