package datasets

import (
	"darco/proto/controllers"
	"darco/proto/models/occurrence"
	"darco/proto/resolvers"
	"darco/proto/router"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterOccurrenceDatasetsRoutes(r router.Router) {
	group := r.RouteGroup("/datasets/occurrence").WithTags([]string{"Datasets"})

	router.Register(group, "GetOccurrenceDataset",
		huma.Operation{
			Path:        "/{slug}",
			Method:      http.MethodGet,
			Summary:     "Get occurrence dataset",
			Description: "Get infos for an occurrence dataset",
		}, controllers.GetHandler[*GetDatasetInput](occurrence.GetOccurrenceDataset))

	router.Register(group, "ListOccurrenceDatasets",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List occurrence datasets",
			Description: "List all occurrence datasets",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListOccurrenceDatasets))
}
