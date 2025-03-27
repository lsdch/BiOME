package events

import (
	"context"
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/controllers/occurrences"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/models/vocabulary"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

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
		controllers.CreateHandler[occurrence.SamplingInputWithEvent])

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
	 * OCCURRENCE
	 */
	router.Register(samplingAPI, "SamplingAddExternalOccurrence",
		huma.Operation{
			Tags:        []string{"Occurrences"},
			Path:        "/{id}/occurrences/external",
			Method:      http.MethodPost,
			Summary:     "Add occurrence from sampling",
			Description: "Register new occurrence resulting from the sampling action",
		},
		SamplingAddExternalOccurrence)

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

type SamplingAddExternalOccurrenceInput struct {
	resolvers.AccessRestricted[resolvers.Contributor]
	controllers.UUIDInput
	Body occurrence.ExternalBioMatInput
}

func SamplingAddExternalOccurrence(ctx context.Context, input *SamplingAddExternalOccurrenceInput) (*occurrences.RegisterOccurrenceOutput, error) {
	sampling := input.Identifier()
	biomat := input.Body
	created, err := biomat.Save(input.DB(), sampling)
	if err != nil {
		return nil, controllers.StatusError(err)
	}
	return &occurrences.RegisterOccurrenceOutput{Body: created}, nil
}
