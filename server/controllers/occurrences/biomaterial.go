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

	router.Register(biomat_API, "ListBioMaterial",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List bio-material",
			Description: "Both internal and external",
		}, controllers.ListHandlerWithOpts[*struct {
			resolvers.AuthResolver
			occurrence.ListBioMaterialOptions
		}](occurrence.ListOccurrences))

	router.Register(biomat_API, "GetBioMaterial",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodGet,
			Summary:     "Get bio-material",
			Description: "Both internal and external",
		}, controllers.GetByCodeHandler(occurrence.GetOccurrence))

	router.Register(biomat_API, "CreateExternalOccurrence",
		huma.Operation{
			Path:    "/external",
			Method:  http.MethodPost,
			Summary: "Create external occurrence",
		}, CreateExternalOccurrence)

	router.Register(biomat_API, "UpdateExternalOccurrence",
		huma.Operation{
			Path:    "/external",
			Method:  http.MethodPatch,
			Summary: "Update external occurrence",
		}, controllers.UpdateByCodeHandler[occurrence.ExternalOccurrenceUpdate])

	router.Register(biomat_API, "DeleteBioMaterial",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodDelete,
			Summary:     "Delete bio-material",
			Description: "Delete any (internal/external) bio-material record by its code",
		}, controllers.DeleteByCodeHandler(occurrence.DeleteOccurrence))
}

type CreateExternalOccurrenceInput struct {
	resolvers.AccessRestricted[resolvers.Contributor]
	Body struct {
		Site        occurrence.SiteInput               `json:"site"`
		Sampling    occurrence.SamplingInput           `json:"sampling"`
		Biomaterial occurrence.ExternalOccurrenceInput `json:"bio_material"`
	}
}

func CreateExternalOccurrence(ctx context.Context, input *CreateExternalOccurrenceInput) (*RegisterOccurrenceOutput, error) {
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
