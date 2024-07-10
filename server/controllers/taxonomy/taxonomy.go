package taxonomy

import (
	"context"
	"darco/proto/controllers"
	"darco/proto/db"
	"darco/proto/models/taxonomy"
	_ "darco/proto/models/validations"
	"darco/proto/resolvers"
	"darco/proto/router"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(r router.Router) {
	taxonomyAPI := r.RouteGroup("/taxonomy").
		WithTags([]string{"Taxonomy"})

	router.Register(taxonomyAPI, "GetTaxonomy",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "Get taxonomy",
		}, GetTaxonomy)

	taxaAPI := r.RouteGroup("/taxonomy/taxa").
		WithTags([]string{"Taxonomy"})

	router.Register(taxaAPI, "ListTaxa",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List taxa",
		}, ListTaxa)

	router.Register(taxaAPI, "GetTaxon",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodGet,
			Summary: "Get taxon",
			Errors:  []int{http.StatusNotFound},
		},
		GetTaxon)

	router.Register(taxaAPI, "CreateTaxon",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create taxon",
			Errors:  []int{http.StatusBadRequest},
		},
		controllers.CreateHandler[taxonomy.TaxonInput])

	router.Register(taxaAPI, "UpdateTaxon",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update taxon",
			Errors:  []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound},
		}, controllers.UpdateByCodeHandler[taxonomy.TaxonUpdate](taxonomy.FindByID))

	router.Register(taxaAPI, "DeleteTaxon",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete taxon",
			Errors:  []int{http.StatusNotFound, http.StatusUnauthorized},
		},
		controllers.DeleteByCodeHandler(taxonomy.Delete))
}

type ListTaxaInput struct {
	taxonomy.ListFilters
}
type ListTaxaOutput struct{ Body []taxonomy.TaxonWithParentRef }

func ListTaxa(ctx context.Context, input *ListTaxaInput) (*ListTaxaOutput, error) {
	taxa, err := taxonomy.ListTaxa(db.Client(), input.ListFilters)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to retrieve list of taxa", err)
	}
	return &ListTaxaOutput{Body: taxa}, nil
}

type GetTaxonInput struct{ controllers.CodeInput }
type GetTaxonOutput struct{ Body taxonomy.TaxonWithRelatives }

func GetTaxon(ctx context.Context, input *GetTaxonInput) (*GetTaxonOutput, error) {
	taxon, err := taxonomy.FindByCode(db.Client(), input.Code)
	if err != nil {
		return nil, huma.Error404NotFound(
			fmt.Sprintf("Taxon %s does not exist", input.Code),
		)
	}
	return &GetTaxonOutput{Body: taxon}, nil
}

type GetTaxonomyInput struct {
	resolvers.AuthResolver
	taxonomy.TaxonomyQuery
}
type GetTaxonomyOutput struct {
	Body []taxonomy.Taxonomy
}

func GetTaxonomy(ctx context.Context, input *GetTaxonomyInput) (*GetTaxonomyOutput, error) {
	var taxonomy, err = taxonomy.GetTaxonomy(input.DB(), input.TaxonomyQuery)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to fetch taxonomy", err)
	}
	return &GetTaxonomyOutput{Body: taxonomy}, nil
}
