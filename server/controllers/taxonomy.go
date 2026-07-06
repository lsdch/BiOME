package controllers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
)

type TaxonomyController struct {
	db      *db.DB
	service *services.TaxonomyService
}

func NewTaxonomyController(db *db.DB, service *services.TaxonomyService) *TaxonomyController {
	return &TaxonomyController{
		db:      db,
		service: service,
	}
}

func (c *TaxonomyController) SearchTaxa(
	ctx context.Context,
	input *models.ListTaxaParams,
) (*BodyTransporter[[]models.Taxon], error) {
	taxa, err := c.service.ListTaxa(ctx, c.db, *input)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.Taxon]{Body: taxa}, nil
}

func (c *TaxonomyController) RegisterRoutes(r *router.Router) {
	group := r.RouteGroup("/taxonomy").WithTags([]string{"Taxonomy"})

	router.NewSpec(
		group,
		"SearchTaxa",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/search",
			Summary: "Search for taxa",
		},
		c.SearchTaxa,
	).WithAccessPolicy(auth.Public()).Register(r)
}
