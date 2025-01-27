package datasets

import (
	"context"
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterSiteDatasetsRoutes(r router.Router) {
	group := r.RouteGroup("/datasets/sites").WithTags([]string{"Datasets"})

	router.Register(group, "GetSiteDataset",
		huma.Operation{
			Path:        "/{slug}",
			Method:      http.MethodGet,
			Summary:     "Get site dataset",
			Description: "Get infos for a site dataset",
		}, controllers.GetHandler[*GetDatasetInput](occurrence.GetSiteDataset))

	router.Register(group, "ListSiteDatasets",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List site datasets",
			Description: "List all site datasets",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListSiteDatasets))

	router.Register(group, "CreateSiteDataset",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodPost,
			Summary:     "Create site dataset",
			Description: "Create a new site dataset with new or existing sites",
		}, CreateSiteDataset)

}

type CreateSiteDatasetInput struct {
	Body occurrence.SiteDatasetInput
	resolvers.AccessRestricted[resolvers.Contributor]
}
type CreateSiteDatasetOutput struct {
	Body occurrence.SiteDataset
}

func CreateSiteDataset(ctx context.Context, input *CreateSiteDatasetInput) (*CreateSiteDatasetOutput, error) {
	dataset, errs := input.Body.Validate(input.DB())
	if errs != nil {
		return nil, huma.Error422UnprocessableEntity("Invalid input", errs...)
	}
	created, err := dataset.Save(input.DB())
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create dataset", err)
	}
	return &CreateSiteDatasetOutput{Body: *created}, nil
}
