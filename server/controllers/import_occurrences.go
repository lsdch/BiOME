package controllers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/middleware"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services/imports"
)

type ImportController struct {
	importService *imports.ImportService
}

func NewImportController(importService *imports.ImportService) *ImportController {
	return &ImportController{
		importService: importService,
	}
}

type ImportHashInput struct {
	Hash string `path:"hash" required:"true" doc:"The unique hash of an import."`
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

func (c *ImportController) ImportOccurrences(ctx context.Context, input *OccurrenceBatchInput) (*imports.ImportResolutionState, error) {
	formData := input.RawBody.Data()
	file := formData.File

	state, err := c.importService.InitImportWorkflow(ctx, input.Body.Label, imports.CSVImportParams{
		Reader:    file,
		Separator: input.Body.Separator,
	})
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (c *ImportController) TrackImportStatus(ctx context.Context, input *ImportHashInput, send sse.Sender) {
}

func (c *ImportController) RegisterRoutes(r *router.Router) {
	importsAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/imports/").
			WithTags([]string{"Batch Imports"})
	}

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
				Path:        "/imports/{hash}/status",
				Summary:     "Get import status updates via Server-Sent Events (SSE)",
			},
			auth.Role(biomedb.UserRoleContributor),
		),
		map[string]any{},
		c.TrackImportStatus)

}
