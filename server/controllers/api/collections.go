package api

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func init() {
	collectionsAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/collections").
			WithTags([]string{"Occurrences"})
	}

	router.RegisterSpec(
		collectionsAPI,
		"ListCollections",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List collections",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListCollections),
	)
}
