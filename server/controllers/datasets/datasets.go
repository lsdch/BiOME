package datasets

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/dataset"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"
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

	router.Register(datasets_API, "ListDatasets",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List all datasets",
			Description: "List all datasets with optional filters and category discriminator",
		}, controllers.ListHandlerWithOpts[*struct {
			resolvers.AuthResolver
			dataset.ListDatasetOptions
		}](dataset.ListDatasets))

	router.Register(datasets_API, "TogglePinDataset",
		huma.Operation{
			Path:        "/pin/{slug}",
			Method:      http.MethodPatch,
			Summary:     "Pin/unpin dataset",
			Description: "Pin or unpin dataset from from dashboard priority display",
		}, PinUnpinDataset)

	router.Register(datasets_API, "UpdateDataset",
		huma.Operation{
			Path:        "/edit/{slug}",
			Method:      http.MethodPatch,
			Summary:     "Update dataset",
			Description: "Update dataset metadata",
		}, controllers.UpdateHandler[*UpdateDatasetInput],
	)

	RegisterSiteDatasetsRoutes(r)
	RegisterOccurrenceDatasetsRoutes(r)
	RegisterSeqDatasetsRoutes(r)
	registerProgramRoutes(r)
}

type UpdateDatasetInput struct {
	resolvers.AuthRequired
	Slug string `path:"slug"`
	controllers.UpdateInput[dataset.DatasetUpdate, string, dataset.Dataset]
}

func (i UpdateDatasetInput) Identifier() string {
	return i.Slug
}

type PinDatasetInput struct {
	resolvers.AccessRestricted[resolvers.Admin]
	controllers.SlugInput
}

func PinUnpinDataset(ctx context.Context, input *PinDatasetInput) (*controllers.UpdateHandlerOutput[dataset.Dataset], error) {

	pinned, err := dataset.TogglePinDataset(input.DB(), input.Slug)

	if err = controllers.StatusError(err); err != nil {
		return nil, err
	}
	return &controllers.UpdateHandlerOutput[dataset.Dataset]{
		Body: pinned,
	}, nil
}

// router.Register(datasets_API, "UpdateDataset",
// 	huma.Operation{
// 		Path:        "/{slug}",
// 		Method:      http.MethodPatch,
// 		Summary:     "Update site dataset",
// 		Description: "Update properties of a site dataset",
// 	}, controllers.UpdateHandler[*UpdateDatasetInput])
// type UpdateDatasetInput struct {
// 	resolvers.AuthRequired
// 	Slug string `path:"slug"`
// 	controllers.UpdateInput[occurrence.DatasetUpdate, string, occurrence.AbstractDataset]
// }

// func (u UpdateDatasetInput) Identifier() string {
// 	return u.Slug
// }

// func (u *UpdateDatasetInput) Resolve(ctx huma.Context) []error {
// 	if err := u.AuthRequired.Resolve(ctx); err != nil {
// 		return err
// 	}
// 	dataset, err := occurrence.FindDataset(u.DB(), u.Slug)
// 	if err != nil {
// 		if db.IsNoData(err) {
// 			return []error{huma.Error404NotFound("Item not found", err)}
// 		}
// 		return []error{err}
// 	}
// 	if !(dataset.IsMaintainer(u.UserInner) || u.User.Role.IsGreaterEqual(people.Admin)) {
// 		return []error{huma.Error403Forbidden("Access restricted to admins or dataset maintainers")}
// 	}

// 	return nil
// }
