package events

import (
	"darco/proto/controllers"
	"darco/proto/models/occurrence"
	"darco/proto/resolvers"
	"darco/proto/router"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func registerAbioticParameterRoutes(r router.Router) {
	abioticAPI := r.RouteGroup("/abiotic").
		WithTags([]string{"Sampling"})

	router.Register(abioticAPI, "ListAbioticParameters",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List abiotic parameters",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListAbioticParameters))

	router.Register(abioticAPI, "CreateAbioticParameter",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create abiotic parameter",
		},
		controllers.CreateHandler[occurrence.AbioticParameterInput])

	router.Register(abioticAPI, "UpdateAbioticParameter",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update abiotic parameter",
		},
		controllers.UpdateByCodeHandler[occurrence.AbioticParameterUpdate])

	router.Register(abioticAPI, "DeleteAbioticParameter",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete abiotic parameter",
		},
		controllers.DeleteByCodeHandler(occurrence.DeleteAbioticParameter))

}
