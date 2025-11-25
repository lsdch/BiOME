package api

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func init() {
	abioticAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/abiotic").
			WithTags([]string{"Sampling"})
	}

	router.RegisterSpec(
		abioticAPI,
		"ListAbioticParameters",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List abiotic parameters",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListAbioticParameters),
	)

	router.RegisterSpec(
		abioticAPI,
		"CreateAbioticParameter",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create abiotic parameter",
		},
		controllers.CreateHandler[occurrence.AbioticParameterInput])

	router.RegisterSpec(
		abioticAPI,
		"UpdateAbioticParameter",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update abiotic parameter",
		},
		controllers.UpdateByCodeHandler[occurrence.AbioticParameterUpdate])

	router.RegisterSpec(
		abioticAPI,
		"DeleteAbioticParameter",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete abiotic parameter",
		},
		controllers.DeleteByCodeHandler(occurrence.DeleteAbioticParameter))
}
