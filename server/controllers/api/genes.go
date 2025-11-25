package api

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/sequences"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func init() {
	genesAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/genes").
			WithTags([]string{"Sequences"})
	}

	router.RegisterSpec(
		genesAPI,
		"ListGenes",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List genes",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](sequences.ListGenes),
	)

	router.RegisterSpec(
		genesAPI,
		"CreateGene",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create gene",
		},
		controllers.CreateHandler[sequences.GeneInput])

	router.RegisterSpec(
		genesAPI,
		"UpdateGene",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update gene",
		},
		controllers.UpdateByCodeHandler[sequences.GeneUpdate])

	router.RegisterSpec(
		genesAPI,
		"DeleteGene",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete gene",
		},
		controllers.DeleteByCodeHandler(sequences.DeleteGene))
}
