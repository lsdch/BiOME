package location

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/location"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(r router.Router) {
	locationAPI := r.RouteGroup("/locations").
		WithTags([]string{"Location", "Countries"})

	router.Register(locationAPI, "ListCountries",
		huma.Operation{
			Path:    "/countries",
			Method:  http.MethodGet,
			Summary: "List countries",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](location.ListCountries))

	router.Register(locationAPI, "getSitesCountByCountry",
		huma.Operation{
			Path:    "/countries/sites-count",
			Method:  http.MethodGet,
			Summary: "Get country list with sites count",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](location.SitesCountByCountry))
}
