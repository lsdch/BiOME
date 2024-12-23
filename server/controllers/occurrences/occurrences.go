package occurrences

import (
	"darco/proto/controllers"
	"darco/proto/models/occurrence"
	"darco/proto/resolvers"
	"darco/proto/router"
	"net/http"

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
