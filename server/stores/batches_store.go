package stores

import (
	"context"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	csvmodels "github.com/lsdch/biome/models/csv"
)

type BatchesStore struct{}

func NewBatchesStore() *BatchesStore {
	return &BatchesStore{}
}

func (s *BatchesStore) GetImportState(ctx context.Context, q db.Querier, importID uuid.UUID) (state biomedb.ImportBatch, err error) {
	state, err = q.Queries().GetImportState(ctx, importID)
	return state, err
}

func (s *BatchesStore) CreateBatch(ctx context.Context, q db.Querier, userID uuid.UUID, w models.ImportBatchWithFileInput) (models.ImportBatch, error) {
	batch, err := q.Queries().InitImportBatch(ctx, w.ToParams(userID))
	if err != nil {
		return models.ImportBatch{}, err
	}
	return models.ImportBatchFromDB(batch), nil
}

func (s *BatchesStore) ListBatches(ctx context.Context, q db.Querier) ([]models.ImportBatch, error) {
	batchsDB, err := q.Queries().ListImportBatchs(ctx)
	if err != nil {
		return nil, err
	}
	batchs := make([]models.ImportBatch, len(batchsDB))
	for i, batch := range batchsDB {
		batchs[i] = models.ImportBatchFromDB(batch)
	}

	return batchs, nil
}

func (s *BatchesStore) Bootstrap(ctx context.Context, q db.Querier, importID uuid.UUID) error {
	if err := q.Queries().CleanUpStagingImport(ctx, importID); err != nil {
		return err
	}
	if err := q.Queries().CleanupTaxonCandidates(ctx, importID); err != nil {
		return err
	}
	if err := q.Queries().CleanUpTaxonResolution(ctx, importID); err != nil {
		return err
	}
	return q.Queries().CleanUpGBIFDependencies(ctx, importID)
}

func (s *BatchesStore) InsertStaging(
	ctx context.Context, q db.Querier, importID uuid.UUID,
	rows []csvmodels.OccurrenceImportRow,
	taxonDefinitions []models.TaxonDefinition,
	mergeUndated bool,
) error {

	taxonDefinitionsMap := make(map[string]models.TaxonDefinition)
	for _, td := range taxonDefinitions {
		taxonDefinitionsMap[td.Name] = td
	}
	for i := range rows {
		if taxonDef, ok := taxonDefinitionsMap[rows[i].TaxonName]; ok {
			rows[i].WithTaxonDefinition(taxonDef)
		}
	}

	stagingRows := make([]biomedb.CopyImportStagingParams, len(rows))
	for i, row := range rows {
		stagingRows[i] = row.ToStaging(importID, mergeUndated)
	}
	_, err := q.Queries().CopyImportStaging(ctx, stagingRows)
	return err
}

func (s *BatchesStore) InsertStagingCollections(ctx context.Context, q db.Querier, importID uuid.UUID, rows []csvmodels.OccurrenceImportRow) error {

	toStage := []biomedb.StageOccurrenceCollectionsParams{}
	for _, row := range rows {
		for _, collection := range row.Collections {
			toStage = append(toStage, biomedb.StageOccurrenceCollectionsParams{
				OccurrenceID:   row.ID,
				CollectionName: collection.Name,
				Vouchers:       collection.Vouchers,
			})
		}
	}
	_, err := q.Queries().StageOccurrenceCollections(ctx, toStage)
	if err != nil {
		return err
	}
	return nil
}

// func (s *BatchesStore) MaterializeStaging(ctx context.Context, q db.Querier, importID uuid.UUID) (models.ImportBatch, error) {
// 	batch, err := q.Queries().MaterializeImportBatch(ctx, types.MakeULID(), importID)
// 	if err != nil {
// 		return models.ImportBatch{}, err
// 	}
// 	return models.ImportBatchFromDB(batch), nil
// }

func (s *BatchesStore) CheckReadyToMaterialize(ctx context.Context, q db.Querier, importID uuid.UUID) (models.MaterializationReadyCheck, error) {
	ready, err := q.Queries().CheckReadyToMaterialize(ctx, importID)
	if err != nil {
		return models.MaterializationReadyCheck{}, err
	}
	return models.MaterializationReadyCheckFromDB(ready), nil
}

func (s *BatchesStore) DeleteBatch(ctx context.Context, q db.Querier, importID uuid.UUID) error {
	if _, err := q.Queries().DeleteImportBatch(ctx, importID); err != nil {
		return err
	}
	return nil
}

func (s *BatchesStore) SetBatchStatus(ctx context.Context, q db.Querier, importID uuid.UUID, status models.ImportBatchStatus) error {
	if err := q.Queries().SetBatchStatus(ctx, status, importID); err != nil {
		return err
	}
	return nil
}

func (s *BatchesStore) SetBatchCompleted(ctx context.Context, q db.Querier, importID uuid.UUID, completedBy uuid.UUID) error {
	if err := q.Queries().SetBatchCompleted(ctx, completedBy, importID); err != nil {
		return err
	}
	return nil
}
