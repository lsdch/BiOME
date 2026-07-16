package api

// import (
// 	"context"
// 	"net/http"

// 	"github.com/lsdch/biome/controllers"
// 	"github.com/lsdch/biome/models/occurrence"
// 	"github.com/lsdch/biome/models/vocabulary"
// 	"github.com/lsdch/biome/resolvers"
// 	"github.com/lsdch/biome/router"

// 	"github.com/danielgtaylor/huma/v2"
// )

// func init() {
// 	accessPointsAPI := func(r *router.Router) router.Group {
// 		return r.RouteGroup("/access-points").
// 			WithTags([]string{"Sampling"})
// 	}

// 	router.RegisterSpec(
// 		accessPointsAPI,
// 		"ListAccessPoints",
// 		huma.Operation{
// 			Path:    "/",
// 			Method:  http.MethodGet,
// 			Summary: "List access points",
// 		},
// 		controllers.ListHandler[*struct {
// 			resolvers.AuthResolver
// 		}](occurrence.ListAccessPoints),
// 	)

// 	samplingAPI := func(r *router.Router) router.Group {
// 		return r.RouteGroup("/samplings").
// 			WithTags([]string{"Sampling"})
// 	}

// 	router.RegisterSpec(
// 		samplingAPI,
// 		"CreateSampling",
// 		huma.Operation{
// 			Path:    "/",
// 			Method:  http.MethodPost,
// 			Summary: "Create sampling action",
// 		},
// 		controllers.CreateHandler[occurrence.SamplingInputAtSite])

// 	router.RegisterSpec(
// 		samplingAPI,
// 		"UpdateSampling",
// 		huma.Operation{
// 			Path:    "/{id}",
// 			Method:  http.MethodPatch,
// 			Summary: "Update sampling action",
// 		},
// 		controllers.UpdateByIDHandler[occurrence.SamplingUpdate])

// 	router.RegisterSpec(
// 		samplingAPI,
// 		"DeleteSampling",
// 		huma.Operation{
// 			Path:    "/{id}",
// 			Method:  http.MethodDelete,
// 			Summary: "Delete sampling action",
// 		},
// 		controllers.DeleteByIDHandler(occurrence.DeleteSampling))

// 	router.RegisterSpec(
// 		samplingAPI,
// 		"SamplingAddOccurrence",
// 		huma.Operation{
// 			Path:        "/{id}/occurrences/",
// 			Method:      http.MethodPost,
// 			Summary:     "Add occurrence from sampling",
// 			Description: "Register new occurrence resulting from the sampling action",
// 		},
// 		SamplingAddOccurrence)

// 	fixativesAPI := func(r *router.Router) router.Group {
// 		return r.RouteGroup("/fixatives").
// 			WithTags([]string{"Sampling"})
// 	}

// 	router.RegisterSpec(
// 		fixativesAPI,
// 		"ListFixatives",
// 		huma.Operation{
// 			Path:    "/",
// 			Method:  http.MethodGet,
// 			Summary: "List fixatives",
// 		},
// 		controllers.ListHandler[*struct {
// 			resolvers.AuthResolver
// 		}](vocabulary.ListFixatives),
// 	)

// 	router.RegisterSpec(
// 		fixativesAPI,
// 		"CreateFixative",
// 		huma.Operation{
// 			Path:    "/",
// 			Method:  http.MethodPost,
// 			Summary: "Create fixative",
// 		},
// 		controllers.CreateHandler[vocabulary.FixativeInput])

// 	router.RegisterSpec(
// 		fixativesAPI,
// 		"UpdateFixative",
// 		huma.Operation{
// 			Path:    "/{code}",
// 			Method:  http.MethodPatch,
// 			Summary: "Update fixative",
// 		},
// 		controllers.UpdateByCodeHandler[vocabulary.FixativeUpdate])

// 	router.RegisterSpec(
// 		fixativesAPI,
// 		"DeleteFixative",
// 		huma.Operation{
// 			Path:    "/{code}",
// 			Method:  http.MethodDelete,
// 			Summary: "Delete fixative",
// 		},
// 		controllers.DeleteByCodeHandler(vocabulary.DeleteFixative))
// }

// type SamplingAddOccurrenceInput struct {
// 	resolvers.AccessRestricted[resolvers.Contributor]
// 	controllers.NumberInput
// 	Body occurrence.OccurrenceInput
// }

// func SamplingAddOccurrence(ctx context.Context, input *SamplingAddOccurrenceInput) (*RegisterOccurrenceOutput, error) {
// 	sampling := input.Identifier()
// 	biomat := input.Body
// 	created, err := biomat.Save(input.DB(), sampling)
// 	if err != nil {
// 		return nil, controllers.StatusError(err)
// 	}
// 	return &RegisterOccurrenceOutput{Body: created}, nil
// }
