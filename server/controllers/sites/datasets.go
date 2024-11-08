package sites

import (
	"context"
	"darco/proto/controllers"
	"darco/proto/db"
	"darco/proto/models/occurrence"
	"darco/proto/models/people"
	"darco/proto/resolvers"
	"darco/proto/router"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type GetSiteDatasetInput struct {
	resolvers.AuthResolver
	Slug string `path:"slug"`
}

func (i GetSiteDatasetInput) Identifier() string {
	return i.Slug
}

func RegisterDatasetRoutes(r router.Router) {
	datasets_API := r.RouteGroup("/datasets").WithTags([]string{"Location"})

	router.Register(datasets_API, "GetSiteDataset",
		huma.Operation{
			Path:        "/{slug}",
			Method:      http.MethodPost,
			Summary:     "Get site dataset",
			Description: "Get infos for a site dataset",
		}, controllers.GetHandler[*GetSiteDatasetInput](occurrence.FindDataset))

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
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListSiteDatasets))

	router.Register(datasets_API, "UpdateSiteDataset",
		huma.Operation{
			Path:        "/{slug}",
			Method:      http.MethodPatch,
			Summary:     "Update site dataset",
			Description: "Update properties of a site dataset",
		}, controllers.UpdateHandler[*UpdateSiteDatasetInput])
}

type UpdateSiteDatasetInput struct {
	resolvers.AuthRequired
	Slug string `path:"slug"`
	controllers.UpdateInput[occurrence.SiteDatasetUpdate, string, occurrence.SiteDataset]
}

func (u UpdateSiteDatasetInput) Identifier() string {
	return u.Slug
}

func (u *UpdateSiteDatasetInput) Resolve(ctx huma.Context) []error {
	if err := u.AuthRequired.Resolve(ctx); err != nil {
		return err
	}
	dataset, err := occurrence.FindDataset(u.DB(), u.Slug)
	if err != nil {
		if db.IsNoData(err) {
			return []error{huma.Error404NotFound("Item not found", err)}
		}
		return []error{err}
	}
	if !(dataset.IsMaintainer(u.UserInner) || u.User.Role.IsGreaterEqual(people.Admin)) {
		return []error{huma.Error403Forbidden("Access restricted to admins or dataset maintainers")}
	}

	return nil
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
	created, err := dataset.Create(input.DB())
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create dataset", err)
	}
	return &CreateSiteDatasetOutput{Body: *created}, nil
}
