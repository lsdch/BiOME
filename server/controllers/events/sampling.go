package events

import (
	"darco/proto/controllers"
	"darco/proto/models/occurrence"
	"darco/proto/models/vocabulary"
	"darco/proto/resolvers"
	"darco/proto/router"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func registerSamplingRoutes(r router.Router) {

	/**
	 * Samplings
	 */
	samplingAPI := r.RouteGroup("/samplings").
		WithTags([]string{"Sampling"})

	router.Register(samplingAPI, "CreateSampling",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create sampling action",
		},
		controllers.CreateHandler[occurrence.SamplingInput])

	router.Register(samplingAPI, "UpdateSampling",
		huma.Operation{
			Path:    "/{id}",
			Method:  http.MethodPatch,
			Summary: "Update sampling action",
		},
		controllers.UpdateByIDHandler[occurrence.SamplingUpdate])

	router.Register(samplingAPI, "DeleteSampling",
		huma.Operation{
			Path:    "/{id}",
			Method:  http.MethodDelete,
			Summary: "Delete sampling action",
		},
		controllers.DeleteByIDHandler(occurrence.DeleteSampling))

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

	router.Register(fixativesAPI, "UpdateFixative",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update fixative",
		},
		controllers.UpdateByCodeHandler[vocabulary.FixativeUpdate])

	router.Register(fixativesAPI, "DeleteFixative",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete fixative",
		},
		controllers.DeleteByCodeHandler(vocabulary.DeleteFixative))

}
