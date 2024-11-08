package events

import (
	"darco/proto/controllers"
	"darco/proto/models/occurrence"
	"darco/proto/resolvers"
	"darco/proto/router"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func registerProgramRoutes(r router.Router) {
	programsAPI := r.RouteGroup("/programs").
		WithTags([]string{"Events"})

	router.Register(programsAPI, "ListPrograms",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List programs",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListPrograms))

	router.Register(programsAPI, "CreateProgram",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create program",
		},
		controllers.CreateHandler[occurrence.ProgramInput, occurrence.Program])

	router.Register(programsAPI, "UpdateProgram",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update program",
		},
		controllers.UpdateByCodeHandler[occurrence.ProgramUpdate])

	router.Register(programsAPI, "DeleteProgram",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete program",
		},
		controllers.DeleteByCodeHandler(occurrence.DeleteProgram))
}
