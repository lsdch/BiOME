package services

import (
	"darco/proto/controllers"
	"darco/proto/resolvers"
	"darco/proto/router"
	"darco/proto/services/geoapify"
	"net/http"

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
