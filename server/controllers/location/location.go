package location

import (
	"context"
	"net/http"
	"reflect"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/location"
	"github.com/lsdch/biome/models/occurrence"
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

	registry := r.API.OpenAPI().Components.Schemas

	router.Register(locationAPI, "coordinatesToCountry",
		huma.Operation{
			Path:    "/coordinates",
			Method:  http.MethodPost,
			Summary: "Get country from WGS84 coordinates",
			Responses: map[string]*huma.Response{
				"200": {
					Description: "The country that contains the coordinates",
					Content: map[string]*huma.MediaType{
						"application/json": {
							Schema: registry.Schema(reflect.TypeFor[location.Country](), true, ""),
						},
					},
				},
				"204": {Description: "No country matches the provided coordinates", Content: nil},
			},
		}, CoordinatesToCountry)

}

type CoordinatesToCountryInput struct {
	resolvers.AuthResolver
	Body occurrence.LatLongCoords
}
type CoordinatesToCountryOutput struct {
	Body location.Country
}

func CoordinatesToCountry(ctx context.Context, input *CoordinatesToCountryInput) (*CoordinatesToCountryOutput, error) {
	country, err := input.Body.FindCountry(input.DB())
	if db.IsNoData(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &CoordinatesToCountryOutput{country}, nil
	}
}
