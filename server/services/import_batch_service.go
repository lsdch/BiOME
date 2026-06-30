package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/oklog/ulid/v2"
)

type ImportBatchService struct {
	db *db.DB
}

func NewImportBatchService(db *db.DB) *ImportBatchService {
	return &ImportBatchService{
		db: db,
	}
}

func (s *ImportBatchService) GetImportBatch(ctx context.Context, id ulid.ULID) (*models.ImportBatch, error) {
	ib, err := s.db.Queries().GetImportBatch(ctx, id)
	if err != nil {
		return nil, err
	}
	importBatch := models.ImportBatchFromDB(ib)
	return &importBatch, nil
}

func (s *ImportBatchService) GetImportBatchForOccurrence(ctx context.Context, occurrenceID ulid.ULID) (models.Optional[models.ImportBatch], error) {
	ib, err := s.db.Queries().GetImportBatchForOccurrence(ctx, occurrenceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Optional[models.ImportBatch]{}, nil
		}
		return models.Optional[models.ImportBatch]{}, err
	}
	importBatch := models.ImportBatchFromDB(ib)
	return models.NewOptional(importBatch), nil
}

func (s *ImportBatchService) ListImportBatches(ctx context.Context) ([]models.ImportBatch, error) {
	ibs, err := s.db.Queries().ListImportBatches(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]models.ImportBatch, len(ibs))
	for i, ib := range ibs {
		result[i] = models.ImportBatchFromDB(ib)
	}
	return result, nil
}

func (s *ImportBatchService) DeleteImportBatch(ctx context.Context, id ulid.ULID) error {
	err := s.db.Queries().DeleteImportBatch(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ImportBatchService) DeleteImportBatchWithOccurrences(ctx context.Context, id ulid.ULID) error {
	return s.db.WithTx(ctx, func(q *biomedb.Queries) error {
		err := q.DeleteOccurrencesFromBatch(ctx, id)
		if err != nil {
			return err
		}
		err = q.DeleteImportBatch(ctx, id)
		if err != nil {
			return err
		}
		return nil
	})
}
