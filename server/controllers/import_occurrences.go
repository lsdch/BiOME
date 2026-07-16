package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/imports"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/middleware"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
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
		Workflow  JSONField[models.ImportWorkflowInput] `form:"workflow" contentType:"application/json" required:"true"`
		File      huma.FormFile                         `form:"file" contentType:"text/tab-separated-values" required:"true"`
		Separator string                                `form:"separator"`
		QuoteChar string                                `form:"quotes"`
	}]
}

func (c *ImportController) ImportOccurrencesCSV(
	ctx context.Context,
	input *OccurrenceBatchInput,
) (*BodyTransporter[models.ImportWorkflow], error) {
	formData := input.RawBody.Data()
	file := formData.File

	session, ok := auth.SessionFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to retrieve session")
	}
	runner, err := c.manager.NewWorkflow(ctx, session.UserID, models.ImportWorkflowInput{
		Label:       formData.Workflow.Value.Label,
		Description: formData.Workflow.Value.Description,
		AssembledBy: formData.Workflow.Value.AssembledBy,
	})
	if err != nil {
		return nil, err
	}

	err = runner.StartWorkflowCSV(file, rune(input.RawBody.Data().Separator[0]))
	if err != nil {
		if err2 := runner.Delete(); err2 != nil {
			return nil, huma.Error500InternalServerError(
				"failed to import CSV",
				err,
				fmt.Errorf("failed to delete import runner after error: %v", err2),
			)
		}
		return nil, err
	}

	return &BodyTransporter[models.ImportWorkflow]{Body: runner.Workflow()}, nil
}

func (c *ImportController) ListImports(ctx context.Context, _ *struct{}) (*BodyTransporter[[]imports.ImportEvent], error) {
	return &BodyTransporter[[]imports.ImportEvent]{Body: c.manager.Snapshots()}, nil
}

func (c *ImportController) ListImportsForCurrentUser(ctx context.Context, _ *struct{}) (*BodyTransporter[[]imports.ImportEvent], error) {
	session, ok := auth.SessionFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to retrieve session")
	}

	return &BodyTransporter[[]imports.ImportEvent]{Body: c.manager.SnapshotsForUser(session.UserID)}, nil
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
			if event.Workflow.ImportID == input.ID {
				send(sse.Message{Data: event})
			}
		}
	}
}

func (c *ImportController) GetTaxonResolutionState(ctx context.Context, input *UUIDInput) (*BodyTransporter[[]models.TaxonResolutionWithCandidates], error) {
	runner, ok := c.manager.GetRunner(input.ID)
	if !ok {
		return nil, fmt.Errorf("import runner not found for import ID %s", input.ID)
	}

	state, err := runner.TaxonResolver().GetTaxonResolutions(ctx, c.db, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get taxon resolutions for import ID %s: %v", input.ID, err)
	}

	return &BodyTransporter[[]models.TaxonResolutionWithCandidates]{Body: state}, nil
}

type ResolveTaxonInput struct {
	UUIDInput
	Body models.ResolveTaxonInput
}

func (c *ImportController) ResolveTaxon(ctx context.Context, input *ResolveTaxonInput) (*struct{}, error) {
	runner, ok := c.manager.GetRunner(input.ID)
	if !ok {
		return nil, fmt.Errorf("import runner not found for import ID %s", input.ID)
	}

	err := runner.TaxonResolver().ResolveTaxon(ctx, c.db, input.ID, input.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve taxon for import ID %s: %v", input.ID, err)
	}

	return nil, nil
}

func (c *ImportController) GetMethodsResolution(ctx context.Context, input *UUIDInput) (*BodyTransporter[[]models.SamplingMethodResolution], error) {
	runner, ok := c.manager.GetRunner(input.ID)
	if !ok {
		return nil, fmt.Errorf("import runner not found for import ID %s", input.ID)
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
	runner, ok := c.manager.GetRunner(input.ID)
	if !ok {
		return nil, fmt.Errorf("import runner not found for import ID %s", input.ID)
	}

	resolution, err := runner.ResolveMethod(input.ID, input.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve method for import ID %s: %v", input.ID, err)
	}

	return &BodyTransporter[models.SamplingMethodResolution]{Body: resolution}, nil
}

func (c *ImportController) GetFixativesResolution(ctx context.Context, input *UUIDInput) (*BodyTransporter[[]models.SamplingFixativeResolution], error) {
	runner, ok := c.manager.GetRunner(input.ID)
	if !ok {
		return nil, fmt.Errorf("import runner not found for import ID %s", input.ID)
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
	runner, ok := c.manager.GetRunner(input.ID)
	if !ok {
		return nil, fmt.Errorf("import runner not found for import ID %s", input.ID)
	}

	resolution, err := runner.ResolveFixative(input.ID, input.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve fixative for import ID %s: %v", input.ID, err)
	}

	return &BodyTransporter[models.SamplingFixativeResolution]{Body: resolution}, nil
}

func (c *ImportController) Materialize(ctx context.Context, input *UUIDInput) (*BodyTransporter[models.ImportBatch], error) {
	runner, ok := c.manager.GetRunner(input.ID)
	if !ok {
		return nil, fmt.Errorf("import runner not found for import ID %s", input.ID)
	}

	batch, err := runner.Materialize()
	if err != nil {
		return nil, fmt.Errorf("failed to materialize import ID %s: %v", input.ID, err)
	}

	return &BodyTransporter[models.ImportBatch]{Body: *batch}, nil
}

func (c *ImportController) DeleteWorkflow(ctx context.Context, input *UUIDInput) (*struct{}, error) {
	runner, ok := c.manager.GetRunner(input.ID)
	if !ok {
		return nil, fmt.Errorf("import runner not found for import ID %s", input.ID)
	}

	if err := runner.Delete(); err != nil {
		return nil, fmt.Errorf("failed to delete workflow for import ID %s: %v", input.ID, err)
	}

	return &struct{}{}, nil
}

func (c *ImportController) RegisterRoutes(r *router.Router) {
	importsAPI := r.RouteGroup("/imports/batch").WithTags([]string{"Batch Imports"})

	router.NewSpec(
		importsAPI,
		"ListImports",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List import workflows",
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
			Summary: "List import workflows for the current user",
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
		"DeleteWorkflow",
		huma.Operation{
			Method:      http.MethodDelete,
			Path:        "{id}",
			Summary:     "Delete the import workflow",
			Description: "Delete the import workflow and all associated data, including staging tables and taxon resolutions. This operation is irreversible.",
		},
		c.DeleteWorkflow,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleMaintainer)).
		Register(r)

	sse.Register(r.API,
		middleware.WithAccessPolicy(
			huma.Operation{
				OperationID: "ImportStatus",
				Method:      http.MethodGet,
				Path:        "/imports/batch/{id}/status",
				Summary:     "Get import status updates via Server-Sent Events (SSE)",
				Tags:        []string{"Batch Imports"},
				Responses: map[string]*huma.Response{
					"200": {
						Content: map[string]*huma.MediaType{
							"text/event-stream": {
								Schema: r.API.OpenAPI().Components.Schemas.Schema(reflect.TypeFor[imports.ImportEvent](), true, ""),
							},
						},
					},
				},
			},
			auth.Role(biomedb.UserRoleContributor),
		),
		map[string]any{
			// FIXME: event type do not seem to work with hey-api
			"status": imports.ImportEvent{},
		},
		c.TrackImportStatus)

}
