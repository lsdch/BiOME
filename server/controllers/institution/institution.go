package institution

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	_ "github.com/lsdch/biome/models/validations"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(r router.Router) {
	institutionsAPI := r.RouteGroup("/institutions").
		WithTags([]string{"People", "Institution"})

	router.Register(institutionsAPI, "ListInstitutions",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List institutions",
			Errors:  []int{500},
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](people.ListInstitutions),
	)

	router.Register(institutionsAPI, "CreateInstitution",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create institution",
			Errors:  []int{400, 500},
		},
		controllers.CreateHandler[people.InstitutionInput, people.Institution])

	router.Register(institutionsAPI, "UpdateInstitution",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update institution",
			Errors:  []int{400, 500},
		}, controllers.UpdateByCodeHandler[people.InstitutionUpdate])

	router.Register(institutionsAPI, "DeleteInstitution",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete institution",
			Errors:  []int{400, 500},
		},
		controllers.DeleteByCodeHandler(people.DeleteInstitution))
}
