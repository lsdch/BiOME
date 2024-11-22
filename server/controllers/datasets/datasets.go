package datasets

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

type GetDatasetInput struct {
	resolvers.AuthResolver
	Slug string `path:"slug"`
}

func (i GetDatasetInput) Identifier() string {
	return i.Slug
}

func RegisterRoutes(r router.Router) {
	datasets_API := r.RouteGroup("/datasets").WithTags([]string{"Datasets"})

	router.Register(datasets_API, "GetDataset",
		huma.Operation{
			Path:        "/{slug}",
			Method:      http.MethodPost,
			Summary:     "Get site dataset",
			Description: "Get infos for a site dataset",
		}, controllers.GetHandler[*GetDatasetInput](occurrence.FindDataset))

	router.Register(datasets_API, "CreateDataset",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodPost,
			Summary:     "Create site dataset",
			Description: "Create a new site dataset with new or existing sites",
		}, CreateDataset)

	router.Register(datasets_API, "ListDatasets",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List site datasets",
			Description: "List all site datasets",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListDatasets))

	router.Register(datasets_API, "UpdateDataset",
		huma.Operation{
			Path:        "/{slug}",
			Method:      http.MethodPatch,
			Summary:     "Update site dataset",
			Description: "Update properties of a site dataset",
		}, controllers.UpdateHandler[*UpdateDatasetInput])
}

type UpdateDatasetInput struct {
	resolvers.AuthRequired
	Slug string `path:"slug"`
	controllers.UpdateInput[occurrence.DatasetUpdate, string, occurrence.Dataset]
}

func (u UpdateDatasetInput) Identifier() string {
	return u.Slug
}

func (u *UpdateDatasetInput) Resolve(ctx huma.Context) []error {
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

type CreateDatasetInput struct {
	Body occurrence.DatasetInput
	resolvers.AccessRestricted[resolvers.Contributor]
}
type CreateDatasetOutput struct {
	Body occurrence.Dataset
}

func CreateDataset(ctx context.Context, input *CreateDatasetInput) (*CreateDatasetOutput, error) {
	dataset, errs := input.Body.Validate(input.DB())
	if errs != nil {
		return nil, huma.Error422UnprocessableEntity("Invalid input", errs...)
	}
	created, err := dataset.Create(input.DB())
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create dataset", err)
	}
	return &CreateDatasetOutput{Body: *created}, nil
}
