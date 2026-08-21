package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models"
	csvmodels "github.com/lsdch/biome/models/csv"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
	"github.com/lsdch/biome/stores"
	"github.com/sirupsen/logrus"
	"github.com/uber/h3-go/v4"
)

type OccurrenceController struct {
	db      *db.DB
	service *services.OccurrencesService
}

func NewOccurrenceController(db *db.DB, service *services.OccurrencesService) *OccurrenceController {
	return &OccurrenceController{
		db:      db,
		service: service,
	}
}

func (c *OccurrenceController) GetOccurrence(ctx context.Context, input *ULIDPath) (*BodyTransporter[models.OccurrenceWithDetails], error) {
	logrus.Debugf("GetOccurrence: ulid=%s", input.ULID.String())
	occurrence, err := c.service.GetOccurrenceWithDetails(ctx, c.db, input.ULID)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[models.OccurrenceWithDetails]{Body: *occurrence}, nil
}

func (c *OccurrenceController) CreateOccurrence(ctx context.Context, input *BodyTransporter[models.FullOccurrenceInput]) (*BodyTransporter[models.OccurrenceWithDetails], error) {
	var occurrence *models.OccurrenceWithDetails
	err := c.db.WithTx(ctx, func(tx *db.Tx) error {
		var txErr error
		occurrence, txErr = c.service.CreateOccurrence(ctx, tx, input.Body)
		return txErr
	})
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[models.OccurrenceWithDetails]{Body: *occurrence}, nil
}

type CreateOccurrenceAtSamplingInput struct {
	UUIDInput
	Body models.OccurrenceInput
}

func (c *OccurrenceController) CreateOccurrenceAtSampling(
	ctx context.Context,
	input *CreateOccurrenceAtSamplingInput,
) (*BodyTransporter[models.OccurrenceWithMetadata], error) {
	var occurrence *models.OccurrenceWithMetadata
	err := c.db.WithTx(ctx, func(tx *db.Tx) error {
		var txErr error
		occurrence, txErr = c.service.CreateOccurrenceAtSampling(ctx, tx, input.ID, input.Body)
		return txErr
	})
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[models.OccurrenceWithMetadata]{Body: *occurrence}, nil
}

func (c *OccurrenceController) ListOccurrences(ctx context.Context, input *stores.ListOccurrencesParams) (*BodyTransporter[models.PaginatedList[models.Occurrence]], error) {
	occurrences, err := c.service.ListOccurrences(ctx, c.db, *input)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[models.PaginatedList[models.Occurrence]]{Body: occurrences}, nil
}

func (c *OccurrenceController) ListOccurrencesAtProximity(ctx context.Context, input *models.ListSamplingsAtProximityInput) (*BodyTransporter[[]*models.SamplingWithOccurrencesAndDistance], error) {
	samplings, err := c.service.ListOccurrencesAtProximity(ctx, c.db, *input)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]*models.SamplingWithOccurrencesAndDistance]{Body: samplings}, nil
}

func (c *OccurrenceController) OccurrencesTaxaOverview(ctx context.Context, input *struct{}) (*BodyTransporter[[]models.OccurrenceOverviewItem], error) {
	overview, err := c.service.OccurrencesByTaxaOverview(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.OccurrenceOverviewItem]{Body: overview}, nil
}

func (c *OccurrenceController) OccurringTaxaAtCell(ctx context.Context, input *struct {
	stores.ListOccurrencesParams
	Resolution int64  `path:"resolution" minimum:"0" maximum:"12"`
	Cell       string `path:"cell"`
}) (*BodyTransporter[[]models.OccurrenceOverviewItem], error) {
	cell := h3.CellFromString(input.Cell)
	overview, err := c.service.ListOccurringTaxaAtCell(ctx, c.db, cell, input.Resolution, input.ListOccurrencesParams)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.OccurrenceOverviewItem]{Body: overview}, nil
}

func (c *OccurrenceController) ListCollectionNames(ctx context.Context, input *struct{}) (*BodyTransporter[[]string], error) {
	collections, err := c.service.ListCollectionNames(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]string]{Body: collections}, nil
}

func (c *OccurrenceController) ListOccurrencesH3(ctx context.Context, input *struct {
	stores.ListOccurrencesParams
	Resolution int64 `path:"resolution" minimum:"0" maximum:"12"`
}) (*BodyTransporter[[]models.H3CellWithRichness], error) {
	cells, err := c.service.ListOccurrencesH3(ctx, c.db, input.Resolution, input.ListOccurrencesParams)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.H3CellWithRichness]{Body: cells}, nil
}

func (c *OccurrenceController) ListSamplingsH3(ctx context.Context, input *struct {
	stores.ListSamplingsParams
	Resolution int64 `path:"resolution" minimum:"0" maximum:"12"`
}) (*BodyTransporter[[]models.H3CellWithRichness], error) {
	cells, err := c.service.ListSamplingsH3(ctx, c.db, input.Resolution, input.ListSamplingsParams)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.H3CellWithRichness]{Body: cells}, nil
}

func (c *OccurrenceController) ListSamplingsWithOccurrences(ctx context.Context, input *struct {
	stores.ListOccurrencesParams
}) (*BodyTransporter[[]models.SamplingWithOccurrences], error) {
	samplings, err := c.service.ListSamplingsWithOccurrences(ctx, c.db, input.ListOccurrencesParams)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.SamplingWithOccurrences]{Body: samplings}, nil
}

type ExportSamplingsWithOccurrencesInput struct {
	stores.ListOccurrencesParams
	Format   models.ExportFormat `query:"format" default:"json"`
	Filename string              `query:"filename" required:"true"`
	csvmodels.ExportCSVOptions
}

func (c *OccurrenceController) ExportSamplingsWithOccurrences(ctx context.Context, input *ExportSamplingsWithOccurrencesInput) (*huma.StreamResponse, error) {
	data, err := c.service.ListSamplingsWithOccurrences(ctx, c.db, input.ListOccurrencesParams)
	if err != nil {
		return nil, err
	}
	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"", input.Filename)
	switch input.Format {
	case "json":
		return &huma.StreamResponse{
			Body: func(ctx huma.Context) {
				writer := ctx.BodyWriter()
				ctx.SetHeader("Content-Type", "application/json")
				ctx.SetHeader("Content-Disposition", contentDisposition)
				json.NewEncoder(writer).Encode(data)
			},
		}, nil
	case "csv", "tsv":
		return &huma.StreamResponse{
			Body: func(ctx huma.Context) {
				writer := ctx.BodyWriter()
				ctx.SetHeader("Content-Type", "text/tab-separated-values")
				ctx.SetHeader("Content-Disposition", contentDisposition)
				header := []string{"sampling_id", "site_name", "site_code", "latitude", "longitude", "coordinates_precision_m", "sampling_date", "occurrence_id", "taxon_id", "taxon_name", "taxon_rank"}
				_, err := writer.Write(
					[]byte(strings.Join(header, ",") + "\n"),
				)
				if err != nil {
					logrus.Error("Error writing CSV header:", err)
				}
				for _, row := range data {
					line := []string{
						row.Sampling.ID.String(),
						row.Sampling.Site.Name.GetWithDefault(""),
						row.Sampling.Site.Code.GetWithDefault(""),
						strconv.FormatFloat(row.Sampling.Coordinates.Latitude, 'f', -1, 64),
						strconv.FormatFloat(row.Sampling.Coordinates.Longitude, 'f', -1, 64),
						models.MapOptional(row.Sampling.Coordinates.Precision, func(v int32) string { return strconv.Itoa(int(v)) }).GetWithDefault(""),
						models.MapOptional(row.Sampling.PerformedOn, func(v models.DateWithPrecision) string { return v.String() }).GetWithDefault(""),
					}
					for _, occ := range row.Occurrences {
						occLine := append(line,
							occ.ID.String(),
							occ.Identification.Taxon.ID.String(),
							occ.Identification.Taxon.Name,
							string(occ.Identification.Taxon.Rank),
						)
						_, err := writer.Write(
							[]byte(strings.Join(occLine, ",") + "\n"),
						)
						if err != nil {
							logrus.Error("Error writing CSV line:", err)
						}
					}
				}
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", input.Format)
	}
}

func (c *OccurrenceController) ListSamplingsWithOccurrencesAtCell(ctx context.Context, input *struct {
	stores.ListOccurrencesParams
	Resolution int64  `path:"resolution" minimum:"0" maximum:"12"`
	Cell       string `path:"cell"`
}) (*BodyTransporter[[]models.SamplingWithOccurrences], error) {
	cell := h3.CellFromString(input.Cell)
	samplings, err := c.service.ListSamplingsWithOccurrencesAtCell(ctx, c.db, cell, input.Resolution, input.ListOccurrencesParams)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.SamplingWithOccurrences]{Body: samplings}, nil
}

func (c *OccurrenceController) RegisterRoutes(r *router.Router) {
	// Register routes for occurrences
	occurrencesGroup := r.RouteGroup("/occurrences").WithTags([]string{"Occurrences"})

	router.NewSpec(
		occurrencesGroup,
		"GetOccurrence",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/item/{ulid}",
			Summary: "Get Occurrence with all relevant metadata",
		},
		c.GetOccurrence,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		occurrencesGroup,
		"ListOccurrences",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/",
			Summary: "List occurrences with optional filters and pagination",
		},
		c.ListOccurrences,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		occurrencesGroup,
		"ListSamplingsWithOccurrences",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/by-sampling/",
			Summary: "List samplings with occurrences",
		},
		c.ListSamplingsWithOccurrences,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		occurrencesGroup,
		"ExportSamplingsWithOccurrences",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/export/",
			Summary: "Export samplings with occurrences",
			// Responses: map[string]*huma.Response{
			// 	"200": {
			// 		Content: map[string]*huma.MediaType{},
			// 	},
			// },
		},
		c.ExportSamplingsWithOccurrences,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		occurrencesGroup,
		"ListOccurrencesAtProximity",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/proximity",
			Summary: "List samplings with occurrences within a certain distance of a given point",
		},
		c.ListOccurrencesAtProximity,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		occurrencesGroup,
		"CreateOccurrence",
		huma.Operation{
			Method:  http.MethodPost,
			Path:    "/",
			Summary: "Create a new occurrence with its sampling and taxon",
		},
		c.CreateOccurrence,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).Register(r)

	router.NewSpec(
		occurrencesGroup,
		"CreateOccurrenceAtSampling",
		huma.Operation{
			Method:  http.MethodPost,
			Path:    "/samplings/{uuid}",
			Summary: "Create a new occurrence at a specific sampling point",
		},
		c.CreateOccurrenceAtSampling,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).Register(r)

	router.NewSpec(
		occurrencesGroup,
		"OccurrencesTaxaOverview",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/taxa-overview",
			Summary: "Get an overview of occurrences count by taxa",
		},
		c.OccurrencesTaxaOverview,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		occurrencesGroup,
		"ListCollectionNames",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/collection-names",
			Summary: "List all unique collection names",
		},
		c.ListCollectionNames,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		occurrencesGroup,
		"ListOccurrencesH3",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/h3/{resolution}",
			Summary: "List occurrences aggregated by H3 cells",
		},
		c.ListOccurrencesH3,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		occurrencesGroup,
		"ListSamplingsH3",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/h3-samplings/{resolution}",
			Summary: "List samplings aggregated by H3 cells",
		},
		c.ListSamplingsH3,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		occurrencesGroup,
		"ListSamplingsWithOccurrencesAtCell",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/h3/{resolution}/{cell}/data",
			Summary: "List samplings with occurrences at a specific H3 cell",
		},
		c.ListSamplingsWithOccurrencesAtCell,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		occurrencesGroup,
		"ListOccurringTaxaAtCell",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/h3/{resolution}/{cell}/taxa",
			Summary: "List taxa occurring at a specific H3 cell with optional filters and pagination",
		},
		c.OccurringTaxaAtCell,
	).WithAccessPolicy(auth.Public()).Register(r)

}
