package stores

import (
	"context"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	csvmodels "github.com/lsdch/biome/models/csv"
	"github.com/lsdch/biome/types"
)

type WorkflowStore struct{}

func NewWorkflowStore() *WorkflowStore {
	return &WorkflowStore{}
}

func (s *WorkflowStore) GetImportState(ctx context.Context, q db.Querier, importID uuid.UUID) (state biomedb.ImportWorkflow, err error) {
	state, err = q.Queries().GetImportState(ctx, importID)
	return state, err
}

func (s *WorkflowStore) CreateWorkflow(ctx context.Context, q db.Querier, userID uuid.UUID, w models.ImportWorkflowInput) (models.ImportWorkflow, error) {
	workflow, err := q.Queries().InitImportWorkflow(ctx, w.ToParams(userID))
	if err != nil {
		return models.ImportWorkflow{}, err
	}
	return models.ImportWorkflowFromDB(workflow), nil
}

func (s *WorkflowStore) ListWorkflows(ctx context.Context, q db.Querier) ([]models.ImportWorkflow, error) {
	workflowsDB, err := q.Queries().ListImportWorkflows(ctx)
	if err != nil {
		return nil, err
	}
	workflows := make([]models.ImportWorkflow, len(workflowsDB))
	for i, workflow := range workflowsDB {
		workflows[i] = models.ImportWorkflowFromDB(workflow)
	}

	return workflows, nil
}

func (s *WorkflowStore) Bootstrap(ctx context.Context, q db.Querier, importID uuid.UUID) error {
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

func (s *WorkflowStore) InsertStaging(ctx context.Context, q db.Querier, importID uuid.UUID, rows []csvmodels.OccurrenceImportRow) error {
	stagingRows := make([]biomedb.CopyImportStagingParams, len(rows))
	for i, row := range rows {
		stagingRows[i] = row.ToStaging(importID)
	}
	_, err := q.Queries().CopyImportStaging(ctx, stagingRows)
	return err
}

func (s *WorkflowStore) MaterializeStaging(ctx context.Context, q db.Querier, importID uuid.UUID) (models.ImportBatch, error) {
	batch, err := q.Queries().MaterializeImportWorkflow(ctx, types.MakeULID(), importID)
	if err != nil {
		return models.ImportBatch{}, err
	}
	return models.ImportBatchFromDB(batch), nil
}

func (s *WorkflowStore) CheckReadyToMaterialize(ctx context.Context, q db.Querier, importID uuid.UUID) (models.MaterializationReadyCheck, error) {
	ready, err := q.Queries().CheckReadyToMaterialize(ctx, importID)
	if err != nil {
		return models.MaterializationReadyCheck{}, err
	}
	return models.MaterializationReadyCheckFromDB(ready), nil
}

func (s *WorkflowStore) DeleteWorkflow(ctx context.Context, q db.Querier, importID uuid.UUID) error {
	if err := q.Queries().DeleteImportWorkflow(ctx, importID); err != nil {
		return err
	}
	return nil
}
