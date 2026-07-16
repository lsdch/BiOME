package services

import (
	"context"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/types"
)

type DatasetsService struct {
}

func NewDatasetsService() *DatasetsService {
	return &DatasetsService{}
}

func (s *DatasetsService) LoadDatasetsForOccurrence(ctx context.Context, q db.Querier, occurrenceID types.ULID) ([]models.Dataset, error) {
	datasets, err := q.Queries().GetDatasetsForOccurrence(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}
	result := make([]models.Dataset, len(datasets))
	for i, d := range datasets {
		result[i] = models.DatasetFromDB(d)
	}
	return result, nil
}

func (s *DatasetsService) GetDatasetByID(ctx context.Context, q db.Querier, datasetID types.ULID) (*models.Dataset, error) {
	d, err := q.Queries().GetDatasetByID(ctx, datasetID)
	if err != nil {
		return nil, err
	}
	dataset := models.DatasetFromDB(d)
	return &dataset, nil
}

func (s *DatasetsService) ListDatasets(ctx context.Context, q db.Querier) ([]models.Dataset, error) {
	datasets, err := q.Queries().ListDatasets(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]models.Dataset, len(datasets))
	for i, d := range datasets {
		result[i] = models.DatasetFromDB(d)
	}
	return result, nil
}

func (s *DatasetsService) LoadOccurrencesForDataset(ctx context.Context, q db.Querier, datasetID types.ULID) ([]models.Occurrence, error) {
	occurrences, err := q.Queries().ListOccurrencesForDataset(ctx, datasetID)
	if err != nil {
		return nil, err
	}
	result := make([]models.Occurrence, len(occurrences))
	for i, o := range occurrences {
		result[i] = models.OccurrenceFromDB(o.Occurrence, o.Taxon, o.Sampling, o.Country)
	}
	return result, nil
}

func (s *DatasetsService) AddOccurrenceToDataset(ctx context.Context, q db.Querier, datasetID types.ULID, occurrenceID types.ULID) error {
	return q.Queries().AddOccurrenceToDataset(ctx, occurrenceID, datasetID)
}

func (s *DatasetsService) RemoveOccurrenceFromDataset(ctx context.Context, q db.Querier, datasetID types.ULID, occurrenceID types.ULID) error {
	return q.Queries().RemoveOccurrenceFromDataset(ctx, occurrenceID, datasetID)
}
