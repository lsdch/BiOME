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
)

type SamplingController struct {
	db      *db.DB
	service *services.SamplingService
}

func NewSamplingController(db *db.DB, service *services.SamplingService) *SamplingController {
	return &SamplingController{
		db:      db,
		service: service,
	}
}

func (c *SamplingController) ListSamplingsAtProximity(
	ctx context.Context,
	input *models.ListSamplingsAtProximityInput,
) (*BodyTransporter[[]models.SamplingWithDistance], error) {
	samplings, err := c.service.ListSamplingsAtProximity(ctx, c.db, *input)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.SamplingWithDistance]{Body: samplings}, nil
}

func (c *SamplingController) ListSamplingsH3AtProximity(
	ctx context.Context,
	input *models.ListSamplingsAtProximityInput,
) (*BodyTransporter[[]models.H3CellWithRichnessAndDistance], error) {
	samplings, err := c.service.ListSamplingsH3AtProximity(ctx, c.db, *input)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.H3CellWithRichnessAndDistance]{Body: samplings}, nil
}

func (c *SamplingController) CreateSampling(
	ctx context.Context,
	input *BodyTransporter[models.SamplingInput],
) (*BodyTransporter[models.SamplingWithDetails], error) {
	sampling := &models.SamplingWithDetails{}
	err := c.db.WithTx(ctx, func(q *db.Tx) error {
		s, err := c.service.CreateSampling(ctx, q, input.Body)
		sampling = s
		return err
	})
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[models.SamplingWithDetails]{Body: *sampling}, nil
}

func (c *SamplingController) ListAccessPoints(
	ctx context.Context,
	_ *struct{},
) (*BodyTransporter[[]string], error) {
	accessPoints, err := c.service.ListSamplingAccessPoints(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]string]{Body: accessPoints}, nil
}

func (c *SamplingController) ListSamplingYears(
	ctx context.Context,
	_ *struct{},
) (*BodyTransporter[[]int32], error) {
	years, err := c.service.ListSamplingYears(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]int32]{Body: years}, nil
}

func (c *SamplingController) RegisterRoutes(r *router.Router) {

	baseGroup := r.RouteGroup("/").WithTags([]string{"Samplings"})

	router.NewSpec(baseGroup, "ListSamplingYears", huma.Operation{
		Method:      http.MethodGet,
		Path:        "years",
		Summary:     "List sampling years",
		Description: "List all unique years for samplings.",
	},
		c.ListSamplingYears,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(baseGroup, "ListSamplingsAtProximity", huma.Operation{
		Method:      http.MethodGet,
		Path:        "proximity",
		Summary:     "List samplings at proximity",
		Description: "List samplings at proximity to a given point, within a given radius and date range.",
	},
		c.ListSamplingsAtProximity,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(baseGroup, "ListSamplingsH3AtProximity", huma.Operation{
		Method:      http.MethodGet,
		Path:        "proximity/h3",
		Summary:     "List samplings at proximity (H3)",
		Description: "List H3 cells at proximity to a given point, within a given radius and date range.",
	},
		c.ListSamplingsH3AtProximity,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(baseGroup, "ListAccessPoints", huma.Operation{
		Method:      http.MethodGet,
		Path:        "access-points",
		Summary:     "List access points",
		Description: "List all unique access points for samplings.",
	},
		c.ListAccessPoints,
	).WithAccessPolicy(auth.Public()).Register(r)

	samplingsGroup := r.RouteGroup("/samplings").WithTags([]string{"Samplings"})

	router.NewSpec(samplingsGroup, "CreateSampling",
		huma.Operation{
			Method:      http.MethodPost,
			Path:        "/",
			Summary:     "Create a new sampling",
			Description: "Create a new sampling with the provided details.",
		},
		c.CreateSampling,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).Register(r)

	c.RegisterSamplingMethodsRoutes(r)
	c.RegisterSamplingFixativesRoutes(r)
}
