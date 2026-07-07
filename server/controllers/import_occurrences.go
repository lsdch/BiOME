package controllers

import (
	"context"
	"net/http"

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

type OccurrenceInputCSV struct {
	File huma.FormFile `form:"occurrences" contentType:"text/csv" required:"true"`
}

type OccurrenceBatchInputMetadata struct {
	Label     string `json:"label"`
	Separator rune   `json:"separator"`
}

type OccurrenceBatchInput struct {
	RawBody huma.MultipartFormFiles[OccurrenceInputCSV]
	Body    OccurrenceBatchInputMetadata
}

func (c *ImportController) ImportOccurrences(
	ctx context.Context,
	input *OccurrenceBatchInput,
) (*BodyTransporter[models.ImportWorkflow], error) {
	formData := input.RawBody.Data()
	file := formData.File

	runner, err := c.manager.NewWorkflow(ctx, input.Body.Label)
	if err != nil {
		return nil, err
	}

	err = runner.StartWorkflowCSV(file, input.Body.Separator)
	if err != nil {
		return nil, err
	}

	return &BodyTransporter[models.ImportWorkflow]{Body: runner.Workflow()}, nil
}

func (c *ImportController) ListImports(ctx context.Context, _ *struct{}) (*BodyTransporter[[]imports.ImportEvent], error) {
	return &BodyTransporter[[]imports.ImportEvent]{Body: c.manager.Snapshots()}, nil
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

func (c *ImportController) RegisterRoutes(r *router.Router) {
	importsAPI := r.RouteGroup("/imports/").WithTags([]string{"Batch Imports"})

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
		"ImportOccurrences",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Import occurrence data",
		},
		c.ImportOccurrences,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).
		Register(r)

	sse.Register(r.API,
		middleware.WithAccessPolicy(
			huma.Operation{
				OperationID: "ImportStatus",
				Method:      http.MethodGet,
				Path:        "/imports/{id}/status",
				Summary:     "Get import status updates via Server-Sent Events (SSE)",
			},
			auth.Role(biomedb.UserRoleContributor),
		),
		map[string]any{},
		c.TrackImportStatus)

}
