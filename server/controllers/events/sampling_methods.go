package events

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func registerSamplingMethodsRoutes(r router.Router) {

	samplingMethodsAPI := r.RouteGroup("/sampling-methods").
		WithTags([]string{"Sampling"})

	router.Register(samplingMethodsAPI, "ListSamplingMethods",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List sampling methods",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListSamplingMethods))

	router.Register(samplingMethodsAPI, "CreateSamplingMethod",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create sampling method",
		},
		controllers.CreateHandler[occurrence.SamplingMethodInput])

	router.Register(samplingMethodsAPI, "UpdateSamplingMethod",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update sampling method",
		},
		controllers.UpdateByCodeHandler[occurrence.SamplingMethodUpdate])

	router.Register(samplingMethodsAPI, "DeleteSamplingMethod",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete sampling method",
		},
		controllers.DeleteByCodeHandler(occurrence.DeleteSamplingMethod))

}
