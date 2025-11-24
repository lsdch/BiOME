package occurrences

import (
	"context"
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func registerBioMatRoutes(r router.Router) {
	biomat_API := r.RouteGroup("/occurrences").WithTags([]string{"Occurrences"})

	router.Register(biomat_API, "ListOccurrences",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List occurrences",
		}, controllers.ListHandlerWithOpts[*struct {
			resolvers.AuthResolver
			occurrence.ListOccurrencesOptions
		}](occurrence.ListOccurrences))

	router.Register(biomat_API, "GetOccurrence",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodGet,
			Summary: "Get occurrence",
		}, controllers.GetByCodeHandler(occurrence.GetOccurrence))

	router.Register(biomat_API, "CreateOccurrence",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create occurrence",
		}, CreateOccurrence)

	router.Register(biomat_API, "UpdateOccurrence",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update occurrence",
		}, controllers.UpdateByCodeHandler[occurrence.OccurrenceUpdate])

	router.Register(biomat_API, "DeleteOccurrence",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodDelete,
			Summary:     "Delete occurrence",
			Description: "Delete an occurrence record by its code",
		}, controllers.DeleteByCodeHandler(occurrence.DeleteOccurrence))
}

type CreateOccurrenceInput struct {
	resolvers.AccessRestricted[resolvers.Contributor]
	Body struct {
		Site        occurrence.SiteInput       `json:"site"`
		Sampling    occurrence.SamplingInput   `json:"sampling"`
		Biomaterial occurrence.OccurrenceInput `json:"bio_material"`
	}
}

func CreateOccurrence(ctx context.Context, input *CreateOccurrenceInput) (*RegisterOccurrenceOutput, error) {
	site, err := input.Body.Site.Save(input.DB())
	if err != nil {
		return nil, controllers.StatusError(err)
	}
	sampling, err := input.Body.Sampling.Save(input.DB(), site.Code)
	if err != nil {
		return nil, controllers.StatusError(err)
	}
	bioMaterial, err := input.Body.Biomaterial.Save(input.DB(), sampling.Number)
	if err != nil {
		return nil, controllers.StatusError(err)
	}
	return &RegisterOccurrenceOutput{bioMaterial}, nil
}
