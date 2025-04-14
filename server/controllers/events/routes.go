package events

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(r router.Router) {
	registerSamplingRoutes(r)
	registerEventsRoutes(r)
	registerAbioticParameterRoutes(r)
	registerSamplingMethodsRoutes(r)

	huma.Register(r.API, huma.Operation{
		Path:    "/access-points",
		Method:  http.MethodGet,
		Summary: "List access points",
		Tags:    []string{"Sampling"},
	}, controllers.ListHandler[*struct {
		resolvers.AuthResolver
	}](occurrence.ListAccessPoints))

}
