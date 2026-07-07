package stores

import (
	"context"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
)

type WorkflowStore struct{}

func NewImportStore() *WorkflowStore {
	return &WorkflowStore{}
}

func (s *WorkflowStore) GetImportState(ctx context.Context, q db.Querier, importID uuid.UUID) (state biomedb.ImportWorkflow, err error) {
	state, err = q.Queries().GetImportState(ctx, importID)
	return state, err
}

func (s *WorkflowStore) CreateWorkflow(ctx context.Context, q db.Querier, label string) (models.ImportWorkflow, error) {
	workflow, err := q.Queries().InitBatchImport(ctx, label)
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

func (s *WorkflowStore) InsertStaging(ctx context.Context, q db.Querier, importID uuid.UUID, rows []models.OccurrenceImportRow) error {
	stagingRows := make([]biomedb.CopyImportStagingParams, len(rows))
	for i, row := range rows {
		stagingRows[i] = row.ToStaging(importID)
	}
	_, err := q.Queries().CopyImportStaging(ctx, stagingRows)
	return err
}
