package sites

import (
	"context"
	"darco/proto/controllers"
	"darco/proto/models/location"
	"darco/proto/resolvers"
	"darco/proto/router"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterDatasetRoutes(r router.Router) {
	datasets_API := r.RouteGroup("/datasets").WithTags([]string{"Location"})
	router.Register(datasets_API, "CreateSiteDataset",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodPost,
			Summary:     "Create site dataset",
			Description: "Create a new site dataset with new or existing sites",
		}, CreateSiteDataset)

	router.Register(datasets_API, "ListSiteDatasets",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List site datasets",
			Description: "List all site datasets",
		}, controllers.ListHandler(location.ListSiteDatasets))
}

type CreateSiteDatasetInput struct {
	Body location.SiteDatasetInput
	resolvers.AccessRestricted[resolvers.Contributor]
}
type CreateSiteDatasetOutput struct {
	Body location.SiteDataset
}

func CreateSiteDataset(ctx context.Context, input *CreateSiteDatasetInput) (*CreateSiteDatasetOutput, error) {
	dataset, errs := input.Body.Validate(input.DB())
	if errs != nil {
		return nil, huma.Error422UnprocessableEntity("Invalid input", errs...)
	}
	created, err := dataset.Create(input.DB())
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create dataset", err)
	}
	return &CreateSiteDatasetOutput{Body: *created}, nil
}
