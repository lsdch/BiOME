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

func (c *TaxonomyController) GetTaxaAtRank(
	ctx context.Context,
	input *struct {
		Rank models.TaxonRank `path:"rank"`
	},
) (*BodyTransporter[[]models.Taxon], error) {
	taxa, err := c.service.GetTaxaByRank(ctx, c.db, input.Rank)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.Taxon]{Body: taxa}, nil
}

func (c *TaxonomyController) GetTaxonByID(
	ctx context.Context,
	input *UUIDInput,
) (*BodyTransporter[*models.TaxonWithFullLineage], error) {
	taxonWithFullLineage, err := c.service.GetTaxonWithFullLineage(ctx, c.db, input.ID)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[*models.TaxonWithFullLineage]{Body: taxonWithFullLineage}, nil
}

func (c *TaxonomyController) CreateTaxon(
	ctx context.Context,
	input *BodyTransporter[*models.CreateTaxonInput],
) (*BodyTransporter[*models.Taxon], error) {
	taxon, err := c.service.CreateTaxon(ctx, c.db, input.Body)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[*models.Taxon]{Body: taxon}, nil
}

func (c *TaxonomyController) DeleteTaxon(
	ctx context.Context,
	input *UUIDInput,
) (*struct{}, error) {
	err := c.db.WithTx(ctx, func(tx *db.Tx) error {
		return c.service.DeleteTaxonWithOccurrences(ctx, tx, input.ID, true)
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *TaxonomyController) GetGBIFKingdoms(
	ctx context.Context,
	input *struct{},
) (*BodyTransporter[[]models.TaxonGBIF], error) {
	kingdoms, err := c.service.GetGBIFKingdoms(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.TaxonGBIF]{Body: kingdoms}, nil
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

	router.NewSpec(
		group,
		"GetTaxon",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/taxa/{id}",
			Summary: "Get a taxon by ID",
		},
		c.GetTaxonByID,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		group,
		"CreateTaxon",
		huma.Operation{
			Method:  http.MethodPost,
			Path:    "/taxa",
			Summary: "Create a new taxon",
		},
		c.CreateTaxon,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleAdmin)).Register(r)

	router.NewSpec(
		group,
		"GetTaxaAtRank",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/rank/{rank}",
			Summary: "Get all taxa at a specific rank",
		},
		c.GetTaxaAtRank,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		group,
		"DeleteTaxon",
		huma.Operation{
			Method:      http.MethodDelete,
			Path:        "/taxa/{id}",
			Summary:     "Delete a taxon by ID",
			Description: "Deletes a taxon and all of its descendants. This operation will also delete all occurrences associated with the taxon and its descendants.",
		},
		c.DeleteTaxon,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleAdmin)).Register(r)

	router.NewSpec(
		group,
		"GetGBIFKingdoms",
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/gbif/kingdoms",
			Summary:     "Get GBIF kingdoms",
			Description: "Retrieves a list of all GBIF kingdoms. Response is cached to avoid repeated calls to the GBIF API.",
		},
		c.GetGBIFKingdoms,
	).WithAccessPolicy(auth.Public()).Register(r)

}
