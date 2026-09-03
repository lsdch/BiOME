package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/imports"
	"github.com/lsdch/biome/lib/app_errors"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/middleware"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
	"github.com/sirupsen/logrus"
)

type ImportController struct {
	db      *db.DB
	manager *imports.ImportManager
}

func NewImportController(db *db.DB, manager *imports.ImportManager) *ImportController {
	return &ImportController{
		db:      db,
		manager: manager,
	}
}

func (c *ImportController) getRunner(importID uuid.UUID) (*imports.ImportRunner, error) {
	runner, ok := c.manager.GetRunner(importID)
	if !ok {
		return nil, app_errors.NotFoundError(fmt.Errorf("import runner not found for import ID %s", importID))
	}
	return runner, nil
}

type JSONField[T any] struct {
	Value T
}

func (j *JSONField[T]) UnmarshalText(text []byte) error {
	return json.Unmarshal(text, &j.Value)
}
func (j *JSONField[T]) UnmarshalJSON(text []byte) error {
	return json.Unmarshal(text, &j.Value)
}

func (j *JSONField[T]) Schema(r huma.Registry) *huma.Schema {
	return r.Schema(reflect.TypeFor[T](), true, "")
}

type OccurrenceBatchInput struct {
	RawBody huma.MultipartFormFiles[struct {
		Batch                 JSONField[models.ImportBatchInput]  `form:"batch" contentType:"application/json" required:"true"`
		TaxonDefinitions      JSONField[[]models.TaxonDefinition] `form:"taxon_definitions" contentType:"application/json" required:"false" doc:"List of taxon definitions to resolve inconsistent taxa in the import batch."`
		File                  huma.FormFile                       `form:"file" contentType:"text/tab-separated-values" required:"true"`
		Separator             string                              `form:"separator"`
		QuoteChar             string                              `form:"quotes"`
		MergeUndatedSamplings bool                                `form:"merge_undated_samplings" required:"false" doc:"If true, undated samplings with the same location and method will be merged into a single sampling."`
	}]
}

func (c *ImportController) ImportOccurrencesCSV(
	ctx context.Context,
	input *OccurrenceBatchInput,
) (*BodyTransporter[models.ImportBatch], error) {

	session, ok := auth.SessionFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to retrieve session")
	}

	formData := input.RawBody.Data()
	file := formData.File

	runner, err := c.manager.NewBatch(ctx, session.UserID,
		models.ImportBatchWithFileInput{
			ImportBatchInput: formData.Batch.Value,
			File: models.File{
				ContentType: file.ContentType,
				Reader:      file,
				Name:        file.Filename,
				Size:        file.Size,
			},
		})

	logrus.Infof("Parsing file %s with separator '%s' and quote char '%s'", file.Filename, formData.Separator, formData.QuoteChar)
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek file: %v", err)
	}
	err = runner.StartBatchCSV(file, rune(formData.Separator[0]), formData.TaxonDefinitions.Value, formData.MergeUndatedSamplings)
	if err != nil {
		if err2 := runner.Delete(context.Background()); err2 != nil {
			logrus.Errorf("failed to delete import batch: %v", err2)
		}
		return nil, err
	}

	return &BodyTransporter[models.ImportBatch]{Body: runner.Batch()}, nil
}

func (c *ImportController) ListImports(ctx context.Context, _ *struct{}) (*BodyTransporter[[]imports.BatchSnapshot], error) {
	return &BodyTransporter[[]imports.BatchSnapshot]{Body: c.manager.Snapshots()}, nil
}

func (c *ImportController) ListImportsForCurrentUser(ctx context.Context, _ *struct{}) (*BodyTransporter[[]imports.BatchSnapshot], error) {
	session, ok := auth.SessionFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to retrieve session")
	}

	return &BodyTransporter[[]imports.BatchSnapshot]{Body: c.manager.SnapshotsForUser(session.UserID)}, nil
}

func (c *ImportController) TrackImports(ctx context.Context, _ *struct{}, send sse.Sender) {
	events, unsubscribe := c.manager.Subscribe()
	defer unsubscribe()

	// initial state
	for _, runner := range c.manager.Snapshots() {
		send(sse.Message{Data: runner})
	}

	for {
		select {
		case <-ctx.Done():
			return

		case event, ok := <-events:
			if !ok {
				return
			}
			send(sse.Message{Data: event})
		}
	}
}

func (c *ImportController) GetImportStatus(ctx context.Context, input *UUIDInput) (*BodyTransporter[imports.BatchSnapshot], error) {
	runner, err := c.getRunner(input.ID)
	if err != nil {
		return nil, err
	}

	return &BodyTransporter[imports.BatchSnapshot]{Body: runner.Snapshot()}, nil
}

func (c *ImportController) TrackImportStatus(ctx context.Context, input *UUIDInput, send sse.Sender) {
	runner, ok := c.manager.GetRunner(input.ID)
	if !ok {
		return
	}

	events, unsubscribe := c.manager.Subscribe()
	defer unsubscribe()

	// initial state
	send(sse.Message{Data: runner.Snapshot()})
	for {
		select {
		case <-ctx.Done():
			return

		case event, ok := <-events:
			if !ok {
				return
			}
			if event.Batch.ID == input.ID {
				send(sse.Message{Data: event})
			}
		}
	}
}

func (c *ImportController) GetTaxonResolutionState(ctx context.Context, input *UUIDInput) (*BodyTransporter[[]models.TaxonResolutionWithCandidates], error) {
	runner, err := c.getRunner(input.ID)
	if err != nil {
		return nil, err
	}

	state, err := runner.TaxonResolver().GetTaxonResolutions(ctx, c.db, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get taxon resolutions for import ID %s: %v", input.ID, err)
	}

	return &BodyTransporter[[]models.TaxonResolutionWithCandidates]{Body: state}, nil
}

func (c *ImportController) GetBibliographyResolution(ctx context.Context, input *UUIDInput) (*BodyTransporter[[]models.PublicationResolutionWithCandidates], error) {
	runner, err := c.getRunner(input.ID)
	if err != nil {
		return nil, err
	}

	state, err := runner.BibliographyResolver().GetBibliographyResolution(ctx, c.db, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bibliography resolutions for import ID %s: %v", input.ID, err)
	}

	return &BodyTransporter[[]models.PublicationResolutionWithCandidates]{Body: state}, nil
}

type ResolveInput struct {
	UUIDInput
	Body models.ResolveInput
}

func (c *ImportController) ResolveTaxon(ctx context.Context, input *ResolveInput) (*struct{}, error) {
	runner, err := c.getRunner(input.ID)
	if err != nil {
		return nil, err
	}

	err = runner.ResolveTaxon(input.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve taxon for import ID %s: %v", input.ID, err)
	}

	return nil, nil
}

func (c *ImportController) ResolvePublication(ctx context.Context, input *ResolveInput) (*struct{}, error) {
	runner, err := c.getRunner(input.ID)
	if err != nil {
		return nil, err
	}

	err = runner.ResolvePublication(input.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve publication for import ID %s: %v", input.ID, err)
	}

	return nil, nil
}

func (c *ImportController) CreateManualTaxonCandidate(ctx context.Context, input *struct {
	UUIDInput
	Body models.TaxonStagingParams
}) (*struct{}, error) {
	runner, err := c.getRunner(input.ID)
	if err != nil {
		return nil, err
	}
	err = c.db.WithTx(ctx, func(tx *db.Tx) error {
		return runner.TaxonResolver().CreateManualCandidate(ctx, tx, input.ID, input.Body)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create manual taxon candidate for import ID %s: %v", input.ID, err)
	}

	return nil, nil
}

func (c *ImportController) GetMethodsResolution(ctx context.Context, input *UUIDInput) (*BodyTransporter[[]models.SamplingMethodResolution], error) {
	runner, err := c.getRunner(input.ID)
	if err != nil {
		return nil, err
	}

	state, err := runner.GetMethodsResolution()
	if err != nil {
		return nil, fmt.Errorf("failed to get method resolutions for import ID %s: %v", input.ID, err)
	}

	return &BodyTransporter[[]models.SamplingMethodResolution]{Body: state}, nil
}

func (c *ImportController) ResolveMethod(ctx context.Context,
	input *struct {
		UUIDInput
		BodyTransporter[models.SamplingMethodResolutionInput]
	}) (*BodyTransporter[models.SamplingMethodResolution], error) {
	runner, err := c.getRunner(input.ID)
	if err != nil {
		return nil, err
	}

	resolution, err := runner.ResolveMethod(input.ID, input.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve method for import ID %s: %v", input.ID, err)
	}

	return &BodyTransporter[models.SamplingMethodResolution]{Body: resolution}, nil
}

func (c *ImportController) GetFixativesResolution(ctx context.Context, input *UUIDInput) (*BodyTransporter[[]models.SamplingFixativeResolution], error) {
	runner, err := c.getRunner(input.ID)
	if err != nil {
		return nil, err
	}

	state, err := runner.GetFixativesResolution()
	if err != nil {
		return nil, fmt.Errorf("failed to get fixative resolutions for import ID %s: %v", input.ID, err)
	}

	return &BodyTransporter[[]models.SamplingFixativeResolution]{Body: state}, nil
}

func (c *ImportController) ResolveFixative(ctx context.Context,
	input *struct {
		UUIDInput
		BodyTransporter[models.SamplingFixativeResolutionInput]
	}) (*BodyTransporter[models.SamplingFixativeResolution], error) {
	runner, err := c.getRunner(input.ID)
	if err != nil {
		return nil, err
	}

	resolution, err := runner.ResolveFixative(input.ID, input.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve fixative for import ID %s: %v", input.ID, err)
	}

	return &BodyTransporter[models.SamplingFixativeResolution]{Body: resolution}, nil
}

func (c *ImportController) Materialize(ctx context.Context, input *UUIDInput) (*BodyTransporter[models.ImportBatch], error) {
	runner, err := c.getRunner(input.ID)
	if err != nil {
		return nil, err
	}

	session, ok := auth.SessionFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to retrieve session")
	}

	batch, err := runner.Materialize(session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to materialize import ID %s: %v", input.ID, err)
	}

	c.manager.RemoveRunner(input.ID)

	return &BodyTransporter[models.ImportBatch]{Body: *batch}, nil
}

func (c *ImportController) DeleteBatch(ctx context.Context, input *UUIDInput) (*struct{}, error) {
	runner, err := c.getRunner(input.ID)
	if err != nil {
		return nil, err
	}

	if err := runner.Delete(ctx); err != nil {
		return nil, fmt.Errorf("failed to delete batch for import ID %s: %v", input.ID, err)
	}

	return &struct{}{}, nil
}

type BibliographyInput struct {
	UUIDInput
	RawBody huma.MultipartFormFiles[struct {
		File      huma.FormFile `form:"file" contentType:"text/tab-separated-values" required:"true"`
		Separator string        `form:"separator"`
		QuoteChar string        `form:"quotes"`
	}]
}

func (c *ImportController) AddBibliographyCSV(ctx context.Context, input *BibliographyInput) (*BodyTransporter[models.ImportBatch], error) {
	runner, err := c.getRunner(input.ID)
	if err != nil {
		return nil, err
	}

	formData := input.RawBody.Data()
	file := formData.File

	err = runner.AddBibliographyCSV(file, rune(input.RawBody.Data().Separator[0]))
	if err != nil {
		return nil, fmt.Errorf("failed to add bibliography CSV for import ID %s: %v", input.ID, err)
	}

	return &BodyTransporter[models.ImportBatch]{Body: runner.Batch()}, nil
}

func (c *ImportController) RegisterRoutes(r *router.Router) {
	importsAPI := r.RouteGroup("/imports/batches").WithTags([]string{"Batch Imports"})

	router.NewSpec(
		importsAPI,
		"ListImports",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List import batchs",
		},
		c.ListImports,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).
		Register(r)

	router.NewSpec(
		importsAPI,
		"ImportOccurrencesCSV",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Import occurrence data from CSV",
		},
		c.ImportOccurrencesCSV,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).
		Register(r)

	router.NewSpec(
		importsAPI,
		"GetBibliographyResolutions",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/imports/batch/{id}/bibliography",
			Summary: "Get bibliography resolution state and candidates",
		},
		c.GetBibliographyResolution,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).
		Register(r)

	router.NewSpec(
		importsAPI,
		"GetTaxonResolutions",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/imports/batch/{id}/taxonomy",
			Summary: "Get taxonomy resolution state and candidates",
		},
		c.GetTaxonResolutionState,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).
		Register(r)

	router.NewSpec(
		importsAPI,
		"ResolveTaxon",
		huma.Operation{
			Method:  http.MethodPatch,
			Path:    "/imports/batch/{id}/taxonomy",
			Summary: "Resolve taxon for import ID",
		},
		c.ResolveTaxon,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).
		Register(r)

	router.NewSpec(
		importsAPI,
		"ResolvePublication",
		huma.Operation{
			Method:  http.MethodPatch,
			Path:    "/imports/batch/{id}/bibliography",
			Summary: "Resolve publication for import ID",
		},
		c.ResolvePublication,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).
		Register(r)

	router.NewSpec(
		importsAPI,
		"CreateManualTaxonCandidate",
		huma.Operation{
			Method:  http.MethodPost,
			Path:    "/imports/batch/{id}/taxonomy/candidates",
			Summary: "Create a manual taxon candidate for import ID",
		},
		c.CreateManualTaxonCandidate,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).
		Register(r)

	router.NewSpec(
		importsAPI,
		"GetMethodsResolution",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/imports/batch/{id}/sampling-methods",
			Summary: "Get sampling methods resolution state",
		},
		c.GetMethodsResolution,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).
		Register(r)

	router.NewSpec(
		importsAPI,
		"ResolveMethod",
		huma.Operation{
			Method:  http.MethodPatch,
			Path:    "/imports/batch/{id}/sampling-methods",
			Summary: "Resolve sampling method for import ID",
		},
		c.ResolveMethod,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).
		Register(r)

	router.NewSpec(
		importsAPI,
		"GetFixativesResolution",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/imports/batch/{id}/fixatives",
			Summary: "Get sampling fixatives resolution state",
		},
		c.GetFixativesResolution,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).
		Register(r)

	router.NewSpec(
		importsAPI,
		"ResolveFixative",
		huma.Operation{
			Method:  http.MethodPatch,
			Path:    "/imports/batch/{id}/fixatives",
			Summary: "Resolve sampling fixative for import ID",
		},
		c.ResolveFixative,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).
		Register(r)

	router.NewSpec(
		importsAPI,
		"ListImportsForCurrentUser",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/by-user",
			Summary: "List import batchs for the current user",
		},
		c.ListImportsForCurrentUser,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).
		Register(r)

	router.NewSpec(
		importsAPI,
		"MaterializeBatch",
		huma.Operation{
			Method:      http.MethodPost,
			Path:        "{id}/materialize",
			Summary:     "Materialize the import batch into the database",
			Description: "Materialize the import batch into the database, creating samplings, occurrences, and resolving taxa. This operation is irreversible.",
		},
		c.Materialize,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleMaintainer)).
		Register(r)

	router.NewSpec(
		importsAPI,
		"DeleteBatchWorkflow",
		huma.Operation{
			Method:      http.MethodDelete,
			Path:        "/{id}",
			Summary:     "Delete the import batch",
			Description: "Delete the import batch and all associated data, including staging tables and taxon resolutions. This operation is irreversible.",
		},
		c.DeleteBatch,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleAdmin)).
		Register(r)

	router.NewSpec(
		importsAPI,
		"AddBibliographyCSV",
		huma.Operation{
			Method:      http.MethodPost,
			Path:        "{id}/bibliography",
			Summary:     "Add bibliography to occurrences batch",
			Description: "Add bibliography to occurrences batch. This operation will stage the publications and generate candidates for resolution.",
		},
		c.AddBibliographyCSV,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).
		Register(r)

	router.NewSpec(
		importsAPI,
		"GetImportStatus",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "{id}/status",
			Summary: "Get the current status of the import batch",
		},
		c.GetImportStatus,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).
		Register(r)

	sse.Register(r.API,
		middleware.WithAccessPolicy(
			huma.Operation{
				OperationID: "TrackImportStatus",
				Method:      http.MethodGet,
				Path:        "/imports/batch/{id}/status/track",
				Summary:     "Get import status updates via Server-Sent Events (SSE)",
				Tags:        []string{"Batch Imports"},
				Responses: map[string]*huma.Response{
					"200": {
						Content: map[string]*huma.MediaType{
							"text/event-stream": {
								Schema: r.API.OpenAPI().Components.Schemas.Schema(reflect.TypeFor[imports.BatchSnapshot](), true, ""),
							},
						},
					},
				},
			},
			auth.Role(biomedb.UserRoleContributor),
		),
		map[string]any{
			"status": imports.BatchSnapshot{},
		},
		c.TrackImportStatus)

}
