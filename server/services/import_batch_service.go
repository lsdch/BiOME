package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/services/storage"
	"github.com/lsdch/biome/types"
)

type ImportBatchService struct {
	storage storage.RawFileStorage
}

func NewImportBatchService(storage storage.RawFileStorage) *ImportBatchService {
	return &ImportBatchService{storage: storage}
}

func (s *ImportBatchService) GetImportBatch(ctx context.Context, q db.Querier, id uuid.UUID) (models.ImportBatch, error) {
	ib, err := q.Queries().GetImportBatch(ctx, id)
	if err != nil {
		return models.ImportBatch{}, err
	}
	importBatch := models.ImportBatchFromDB(ib)
	return importBatch, nil
}

func (s *ImportBatchService) LoadDatasetsForImportBatch(ctx context.Context, q db.Querier, importBatchID uuid.UUID) ([]models.Dataset, error) {
	datasets, err := q.Queries().LoadDatasetsForImportBatch(ctx, importBatchID)
	if err != nil {
		return nil, err
	}
	result := make([]models.Dataset, len(datasets))
	for i, d := range datasets {
		result[i] = models.DatasetFromDB(d)
	}
	return result, nil
}

func (s *ImportBatchService) GetImportBatchWithDetails(ctx context.Context, q db.Querier, id uuid.UUID) (models.ImportBatchWithDetails, error) {
	ib, err := q.Queries().GetImportBatchWithContent(ctx, id)
	if err != nil {
		return models.ImportBatchWithDetails{}, err
	}
	datasets, err := s.LoadDatasetsForImportBatch(ctx, q, id)
	if err != nil {
		return models.ImportBatchWithDetails{}, err
	}
	importBatch := models.ImportBatchFromDB(ib.ImportBatch).
		WithContent(ib.OccurrenceCount, ib.SamplingCount, models.UserFromDB(ib.User), models.UserFromDB(ib.User_2)).
		WithDatasets(datasets)
	return importBatch, nil
}

func (s *ImportBatchService) GetImportBatchForOccurrence(ctx context.Context, q db.Querier, occurrenceID types.ULID) (models.Optional[models.ImportBatch], error) {
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

func (s *ImportBatchService) ListImportBatchesWithContent(ctx context.Context, q db.Querier) ([]models.ImportBatchListItem, error) {
	ibs, err := q.Queries().ListImportBatchesWithContent(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]models.ImportBatchListItem, len(ibs))
	for i, ib := range ibs {
		result[i] = models.ImportBatchFromDB(ib.ImportBatch).
			WithContent(ib.OccurrenceCount, ib.SamplingCount, models.UserFromDB(ib.User), models.UserFromDB(ib.User_2))
	}
	return result, nil
}

func (s *ImportBatchService) DeleteImportBatch(ctx context.Context, tx *db.Tx, id uuid.UUID) error {
	err := tx.Queries().DeleteOccurrencesFromBatch(ctx, id)
	if err != nil {
		return err
	}
	_, err = tx.Queries().DeleteImportBatch(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ImportBatchService) DeleteImportBatchWithOccurrences(ctx context.Context, tx *db.Tx, id uuid.UUID) error {
	q := tx.Queries()
	if err := q.DeleteOccurrencesFromBatch(ctx, id); err != nil {
		return err
	}
	if _, err := q.DeleteImportBatch(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *ImportBatchService) AddToDataset(ctx context.Context, q db.Querier, importBatchID uuid.UUID, datasetID types.ULID) error {
	return q.Queries().AddImportBatchToDataset(ctx, datasetID, importBatchID)
}

func (s *ImportBatchService) RemoveFromDataset(ctx context.Context, q db.Querier, importBatchID uuid.UUID, datasetID types.ULID) error {
	return q.Queries().RemoveImportBatchFromDataset(ctx, datasetID, importBatchID)
}
