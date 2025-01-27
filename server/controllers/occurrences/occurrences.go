package occurrences

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(r router.Router) {

	occurAPI := r.RouteGroup("/occurrences").
		WithTags([]string{"Occurrences"})

	router.Register(occurAPI, "OccurrenceOverview",
		huma.Operation{
			Path:    "/overview",
			Method:  http.MethodGet,
			Summary: "Occurrences overview",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.OccurrenceOverview),
	)
}
