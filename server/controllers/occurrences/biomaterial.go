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
	biomat_API := r.RouteGroup("/bio-material").WithTags([]string{"Occurrences"})

	router.Register(biomat_API, "ListBioMaterial",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List bio-material",
			Description: "Both internal and external",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListBioMaterials))

	router.Register(biomat_API, "GetBioMaterial",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodGet,
			Summary:     "Get bio-material",
			Description: "Both internal and external",
		}, controllers.GetByCodeHandler(occurrence.GetBioMaterial))

	router.Register(biomat_API, "CreateExternalBioMat",
		huma.Operation{
			Path:    "/external",
			Method:  http.MethodPost,
			Summary: "Create external bio-material",
		}, CreateExternalBioMat)

	router.Register(biomat_API, "UpdateExternalBioMat",
		huma.Operation{
			Path:    "/external",
			Method:  http.MethodPatch,
			Summary: "Update external bio-material",
		}, controllers.UpdateByCodeHandler[occurrence.ExternalBioMatUpdate])

	router.Register(biomat_API, "DeleteBioMaterial",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodDelete,
			Summary:     "Delete bio-material",
			Description: "Delete any (internal/external) bio-material record by its code",
		}, controllers.DeleteByCodeHandler(occurrence.DeleteBioMaterial))
}

type CreateExternalBioMatInput struct {
	resolvers.AccessRestricted[resolvers.Contributor]
	Body struct {
		Site        occurrence.SiteInput           `json:"site"`
		Event       occurrence.EventInput          `json:"event"`
		Sampling    occurrence.SamplingInput       `json:"sampling"`
		Biomaterial occurrence.ExternalBioMatInput `json:"bio_material"`
	}
}

func CreateExternalBioMat(ctx context.Context, input *CreateExternalBioMatInput) (*RegisterOccurrenceOutput, error) {
	site, err := input.Body.Site.Save(input.DB())
	if err != nil {
		return nil, controllers.StatusError(err)
	}
	event, err := input.Body.Event.Save(input.DB(), site.Code)
	if err != nil {
		return nil, controllers.StatusError(err)
	}
	sampling, err := input.Body.Sampling.Save(input.DB(), event.ID)
	if err != nil {
		return nil, controllers.StatusError(err)
	}
	bioMaterial, err := input.Body.Biomaterial.Save(input.DB(), sampling.ID)
	if err != nil {
		return nil, controllers.StatusError(err)
	}
	return &RegisterOccurrenceOutput{bioMaterial}, nil
}
