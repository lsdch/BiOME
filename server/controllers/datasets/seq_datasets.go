package datasets

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterSeqDatasetsRoutes(r router.Router) {
	group := r.RouteGroup("/datasets/sequences").WithTags([]string{"Datasets"})

	router.Register(group, "GetSequenceDataset",
		huma.Operation{
			Path:        "/{slug}",
			Method:      http.MethodGet,
			Summary:     "Get sequence dataset",
			Description: "Get infos for an sequence dataset",
		}, controllers.GetHandler[*GetDatasetInput](occurrence.GetSequenceDataset))

	router.Register(group, "ListSequenceDatasets",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List sequence datasets",
			Description: "List all sequence datasets",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListSequenceDatasets))
}
