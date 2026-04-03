package taxonomy

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/taxonomy"
	"github.com/lsdch/biome/router"
)

func init() {
	taxonomyAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/taxonomy/by-rank").
			WithTags([]string{"Taxonomy"})
	}

	router.RegisterSpec(
		taxonomyAPI,
		"GetTaxonomyAtRank",
		huma.Operation{
			Path:    "/{rank}",
			Method:  http.MethodGet,
			Summary: "Get taxonomy",
		},
		GetTaxonomyAtRank)

	taxaAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/taxonomy/taxa").
			WithTags([]string{"Taxonomy"})
	}

	router.RegisterSpec(
		taxaAPI,
		"ListTaxa",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List taxa",
		},
		ListTaxa)

	router.RegisterSpec(
		taxaAPI,
		"GetTaxon",
		huma.Operation{
			Path:    "/{name}",
			Method:  http.MethodGet,
			Summary: "Get taxon",
			Errors:  []int{http.StatusNotFound},
		},
		GetTaxon)

	router.RegisterSpec(
		taxaAPI,
		"CreateTaxon",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create taxon",
			Errors:  []int{http.StatusBadRequest},
		},
		controllers.CreateHandlerWithInput[*CreateTaxonInput, taxonomy.TaxonInput, taxonomy.TaxonWithRelatives])

	router.RegisterSpec(
		taxaAPI,
		"UpdateTaxon",
		huma.Operation{
			Path:    "/{name}",
			Method:  http.MethodPatch,
			Summary: "Update taxon",
			Errors:  []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound},
		},
		controllers.UpdateByNameHandler[taxonomy.TaxonUpdate])

	router.RegisterSpec(
		taxaAPI,
		"DeleteTaxon",
		huma.Operation{
			Path:    "/{name}",
			Method:  http.MethodDelete,
			Summary: "Delete taxon",
			Errors:  []int{http.StatusNotFound, http.StatusUnauthorized},
		},
		controllers.DeleteByNameHandler(taxonomy.Delete))
}
