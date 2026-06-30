package services

import (
	"context"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/oklog/ulid/v2"
)

type DatasetsService struct {
	db db.Querier
}

func NewDatasetsService(db db.Querier) *DatasetsService {
	return &DatasetsService{
		db: db,
	}
}

func (s *DatasetsService) LoadDatasetsForOccurrence(ctx context.Context, occurrenceID ulid.ULID) ([]models.Dataset, error) {
	datasets, err := s.db.Queries().GetDatasetsForOccurrence(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}
	result := make([]models.Dataset, len(datasets))
	for i, d := range datasets {
		result[i] = models.DatasetFromDB(d)
	}
	return result, nil
}

func (s *DatasetsService) GetDatasetByID(ctx context.Context, datasetID ulid.ULID) (*models.Dataset, error) {
	d, err := s.db.Queries().GetDatasetByID(ctx, datasetID)
	if err != nil {
		return nil, err
	}
	dataset := models.DatasetFromDB(d)
	return &dataset, nil
}

func (s *DatasetsService) ListDatasets(ctx context.Context) ([]models.Dataset, error) {
	datasets, err := s.db.Queries().ListDatasets(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]models.Dataset, len(datasets))
	for i, d := range datasets {
		result[i] = models.DatasetFromDB(d)
	}
	return result, nil
}

func (s *DatasetsService) LoadOccurrencesForDataset(ctx context.Context, datasetID ulid.ULID) ([]models.Occurrence, error) {
	occurrences, err := s.db.Queries().ListOccurrencesForDataset(ctx, datasetID)
	if err != nil {
		return nil, err
	}
	result := make([]models.Occurrence, len(occurrences))
	for i, o := range occurrences {
		result[i] = models.OccurrenceFromDB(o.Occurrence, o.Taxon, o.Sampling, o.Country)
	}
	return result, nil
}
