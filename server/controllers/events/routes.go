package events

import (
	"darco/proto/controllers"
	"darco/proto/models/occurrence"
	"darco/proto/resolvers"
	"darco/proto/router"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(r router.Router) {

	registerProgramRoutes(r)
	registerSamplingRoutes(r)
	registerEventsRoutes(r)
	registerAbioticParameterRoutes(r)

	huma.Register(r.API, huma.Operation{
		Path:    "/access-points",
		Method:  http.MethodGet,
		Summary: "List access points",
		Tags:    []string{"Sampling"},
	}, controllers.ListHandler[*struct {
		resolvers.AuthResolver
	}](occurrence.ListAccessPoints))

}
