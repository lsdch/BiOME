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

func RegisterRoutes(r router.Router) {
	sites_API := r.RouteGroup("/sites").WithTags([]string{"Location"})

	router.Register(sites_API, "CreateSiteDataset",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodPost,
			Summary:     "Create site dataset",
			Description: "Create a new site dataset with new or existing sites",
		}, CreateSiteDataset)

	huma.Register(r.API, huma.Operation{
		Path:    "/access-points",
		Method:  http.MethodGet,
		Summary: "List access points",
		Tags:    sites_API.Tags,
	}, controllers.ListHandler(location.ListAccessPoints))
}

type CreateSiteDatasetInput struct {
	Body location.SiteDatasetInput
	resolvers.AccessRestricted[resolvers.Contributor]
}
type CreateSiteDatasetOutput struct {
	Body location.SiteDataset
}

func CreateSiteDataset(ctx context.Context, input *CreateSiteDatasetInput) (*CreateSiteDatasetOutput, error) {
	_, errs := input.Body.Validate(input.DB())
	if errs != nil {
		return nil, huma.Error422UnprocessableEntity("Invalid input", errs...)
	}
	return &CreateSiteDatasetOutput{}, nil
}
