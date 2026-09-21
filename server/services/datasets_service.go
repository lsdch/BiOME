package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

func (s *DatasetsService) loadCurators(ctx context.Context, q db.Querier) (map[types.ULID][]models.User, error) {
	users, err := q.Queries().LoadDatasetCurators(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[types.ULID][]models.User)
	for _, u := range users {
		result[u.DatasetID] = append(result[u.DatasetID], models.UserFromDB(u.User))
	}
	return result, nil
}

func (s *DatasetsService) loadCuratorsForDataset(ctx context.Context, q db.Querier, datasetID types.ULID) ([]models.User, error) {
	curators, err := s.loadCurators(ctx, q)
	if err != nil {
		return nil, err
	}
	return curators[datasetID], nil
}

func (s *DatasetsService) ListImportBatchesInDataset(ctx context.Context, q db.Querier, datasetID types.ULID) ([]models.ImportBatchListItem, error) {
	ibs, err := q.Queries().ListImportBatchesInDataset(ctx, datasetID)
	if err != nil {
		return nil, err
	}
	result := make([]models.ImportBatchListItem, len(ibs))
	for i, ib := range ibs {
		result[i] = models.ImportBatchFromDB(ib.ImportBatch).WithContent(
			ib.OccurrenceCount, ib.SamplingCount,
			models.UserFromDB(ib.User), models.UserFromDB(ib.User_2),
		)
	}
	return result, nil
}

func (s *DatasetsService) GetDatasetByID(ctx context.Context, q db.Querier, datasetID types.ULID) (*models.DatasetWithMaintainers, error) {
	d, err := q.Queries().GetDatasetByID(ctx, datasetID)
	if err != nil {
		return nil, err
	}
	curators, err := s.loadCuratorsForDataset(ctx, q, datasetID)
	if err != nil {
		return nil, err
	}
	dataset := models.DatasetFromDB(d.Dataset).WithMaintainers(curators)
	return &dataset, nil
}

func (s *DatasetsService) ListDatasets(ctx context.Context, q db.Querier) ([]models.DatasetWithSummary, error) {
	return s.ListDatasetsForUser(ctx, q, nil)
}

func (s *DatasetsService) ListDatasetsForUser(ctx context.Context, q db.Querier, userID *uuid.UUID) ([]models.DatasetWithSummary, error) {
	uID := pgtype.UUID{Valid: false}
	if userID != nil {
		uID = pgtype.UUID{Bytes: *userID, Valid: true}
	}
	datasets, err := q.Queries().ListDatasets(ctx, uID)
	if err != nil {
		return nil, err
	}
	curators, err := s.loadCurators(ctx, q)
	if err != nil {
		return nil, err
	}
	result := make([]models.DatasetWithSummary, len(datasets))
	for i, d := range datasets {
		result[i] = models.DatasetFromDB(d.Dataset).WithMaintainers(curators[d.Dataset.ID]).WithSummary(
			d.OccurrenceCount, d.SamplingCount, d.ImportBatchCount,
		)
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
		result[i] = models.OccurrenceFromDB(o.Occurrence, o.Taxon, o.SamplingsWithCountry)
	}
	return result, nil
}

func (s *DatasetsService) AddOccurrenceToDataset(ctx context.Context, q db.Querier, datasetID types.ULID, occurrenceID types.ULID) error {
	return q.Queries().AddOccurrenceToDataset(ctx, occurrenceID, datasetID)
}

func (s *DatasetsService) RemoveOccurrenceFromDataset(ctx context.Context, q db.Querier, datasetID types.ULID, occurrenceID types.ULID) error {
	return q.Queries().RemoveOccurrenceFromDataset(ctx, occurrenceID, datasetID)
}

func (s *DatasetsService) CreateDataset(ctx context.Context, tx *db.Tx, ownerID uuid.UUID, input models.DatasetInput) (*models.DatasetWithMaintainers, error) {
	params := input.ToParams(ownerID)
	d, err := tx.Queries().CreateDataset(ctx, params)
	if err != nil {
		return nil, err
	}
	for _, curatorID := range input.Maintainers {
		if err := s.AddCuratorToDataset(ctx, tx, d.ID, curatorID); err != nil {
			return nil, err
		}
	}
	curators, err := s.loadCuratorsForDataset(ctx, tx, d.ID)
	if err != nil {
		return nil, err
	}
	dataset := models.DatasetFromDB(d).WithMaintainers(curators)
	return &dataset, nil
}

func (s *DatasetsService) AddCuratorToDataset(ctx context.Context, q db.Querier, datasetID types.ULID, userID uuid.UUID) error {
	return q.Queries().DatasetAddCurator(ctx, datasetID, userID)
}

func (s *DatasetsService) RemoveCuratorFromDataset(ctx context.Context, q db.Querier, datasetID types.ULID, userID uuid.UUID) error {
	return q.Queries().DatasetRemoveCurator(ctx, datasetID, userID)
}
