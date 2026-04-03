package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/dataset"
	"github.com/lsdch/biome/models/occurrence"
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

func init() {
	datasetsAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/datasets").WithTags([]string{"Datasets"})
	}

	occDatasetsAPI := func(r *router.Router) router.Group {
		return datasetsAPI(r).RouteGroup("/occurrence")
	}

	seqDatasetsAPI := func(r *router.Router) router.Group {
		return datasetsAPI(r).RouteGroup("/sequence")
	}

	siteDatasetsAPI := func(r *router.Router) router.Group {
		return datasetsAPI(r).RouteGroup("/site")
	}

	router.RegisterSpec(
		datasetsAPI,
		"ListDatasets",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List all datasets",
			Description: "List all datasets with optional filters and category discriminator",
		}, controllers.ListHandlerWithOpts[*struct {
			resolvers.AuthResolver
			dataset.ListDatasetOptions
		}](dataset.ListDatasets))

	router.RegisterSpec(
		datasetsAPI,
		"GetDataset",
		huma.Operation{
			Path:        "/{slug}",
			Method:      http.MethodGet,
			Summary:     "Get dataset",
			Description: "Retrieve dataset infos by slug",
		}, controllers.GetHandler[*GetDatasetInput](dataset.GetDataset))

	router.RegisterSpec(
		datasetsAPI,
		"TogglePinDataset",
		huma.Operation{
			Path:        "/pin/{slug}",
			Method:      http.MethodPatch,
			Summary:     "Pin/unpin dataset",
			Description: "Pin or unpin dataset from from dashboard priority display",
		}, PinUnpinDataset)

	router.RegisterSpec(
		datasetsAPI,
		"UpdateDataset",
		huma.Operation{
			Path:        "/edit/{slug}",
			Method:      http.MethodPatch,
			Summary:     "Update dataset",
			Description: "Update dataset metadata",
		}, controllers.UpdateHandler[*UpdateDatasetInput])

	/* --------------------------------
	 * Occurrence datasets
	 * -------------------------------- */

	router.RegisterSpec(
		occDatasetsAPI,
		"GetOccurrenceDataset",
		huma.Operation{
			Path:        "/{slug}",
			Method:      http.MethodGet,
			Summary:     "Get occurrence dataset",
			Description: "Get infos for an occurrence dataset",
		}, controllers.GetHandler[*GetDatasetInput](occurrence.GetOccurrenceDataset))

	router.RegisterSpec(
		occDatasetsAPI,
		"UpdateOccurrenceCodesInDataset",
		huma.Operation{
			Path:        "/{slug}/update-codes",
			Method:      http.MethodPatch,
			Summary:     "Update occurrence codes in dataset",
			Description: "Update occurrence codes based on the current taxon and sampling data",
		}, controllers.GetHandler[*GetDatasetInput](occurrence.UpdateOccurrenceCodesInDataset))

	router.RegisterSpec(
		occDatasetsAPI,
		"ListOccurrenceDatasets",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List occurrence datasets",
			Description: "List all occurrence datasets",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListOccurrenceDatasets))

	/* --------------------------------
	 * Sequence datasets
	 * -------------------------------- */

	router.RegisterSpec(
		seqDatasetsAPI,
		"GetSequenceDataset",
		huma.Operation{
			Path:        "/{slug}",
			Method:      http.MethodGet,
			Summary:     "Get sequence dataset",
			Description: "Get infos for an sequence dataset",
		}, controllers.GetHandler[*GetDatasetInput](occurrence.GetSequenceDataset))

	router.RegisterSpec(
		seqDatasetsAPI,
		"ListSequenceDatasets",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List sequence datasets",
			Description: "List all sequence datasets",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListSequenceDatasets))

	/* --------------------------------
	 * Site datasets
	 * -------------------------------- */

	router.RegisterSpec(
		siteDatasetsAPI,
		"GetSiteDataset",
		huma.Operation{
			Path:        "/{slug}",
			Method:      http.MethodGet,
			Summary:     "Get site dataset",
			Description: "Get infos for a site dataset",
		}, controllers.GetHandler[*GetDatasetInput](occurrence.GetSiteDataset))

	router.RegisterSpec(
		siteDatasetsAPI,
		"ListSiteDatasets",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List site datasets",
			Description: "List all site datasets",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListSiteDatasets))

	router.RegisterSpec(
		siteDatasetsAPI,
		"CreateSiteDataset",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodPost,
			Summary:     "Create site dataset",
			Description: "Create a new site dataset with new or existing sites",
		}, CreateSiteDataset)
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
