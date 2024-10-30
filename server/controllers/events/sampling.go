package events

import (
	"darco/proto/controllers"
	"darco/proto/models/events"
	"darco/proto/models/vocabulary"
	"darco/proto/resolvers"
	"darco/proto/router"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func registerSamplingRoutes(r router.Router) {

	/**
	 * Sampling methods
	 */
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
		}](events.ListSamplingMethods))

	router.Register(samplingMethodsAPI, "CreateSamplingMethod",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create sampling method",
		},
		controllers.CreateHandler[events.SamplingMethodInput])

	/**
	 * FIXATIVES
	 */
	fixativesAPI := r.RouteGroup("/fixatives").
		WithTags([]string{"Sampling"})

	router.Register(fixativesAPI, "ListFixatives",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List fixatives",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](vocabulary.ListFixatives))

	router.Register(fixativesAPI, "CreateFixative",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create fixative",
		},
		controllers.CreateHandler[vocabulary.FixativeInput])

	/**
	 * Abiotic parameters
	 */
	abioticAPI := r.RouteGroup("/abiotic").
		WithTags([]string{"Sampling"})

	router.Register(abioticAPI, "ListAbioticParameters",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List abiotic parameters",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](events.ListAbioticParameters))

	router.Register(abioticAPI, "CreateAbioticParameter",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create abiotic parameter",
		},
		controllers.CreateHandler[events.AbioticParameterInput])

}
