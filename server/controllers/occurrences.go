package controllers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
	"github.com/lsdch/biome/stores"
	"github.com/sirupsen/logrus"
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

func (c *OccurrenceController) ListCollectionNames(ctx context.Context, input *struct{}) (*BodyTransporter[[]string], error) {
	collections, err := c.service.ListCollectionNames(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]string]{Body: collections}, nil
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
}
