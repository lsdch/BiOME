package services

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services/geoapify"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterGeoapifyRoutes(r router.Router) {
	genesAPI := r.RouteGroup("/geoapify").
		WithTags([]string{"Services"})

	router.Register(genesAPI, "ListGeoapifyUsage",
		huma.Operation{Path: "/",
			Method:  http.MethodGet,
			Summary: "List Geoapify usage",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](geoapify.GeoapifyUsageList),
	)
}
