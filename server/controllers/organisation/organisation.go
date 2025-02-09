package organisation

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
	organisationsAPI := r.RouteGroup("/organisations").
		WithTags([]string{"People", "Organisation"})

	router.Register(organisationsAPI, "ListOrganisations",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List organisations",
			Errors:  []int{500},
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](people.ListOrganisations),
	)

	router.Register(organisationsAPI, "CreateOrganisation",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create organisation",
			Errors:  []int{400, 500},
		},
		controllers.CreateHandler[people.OrganisationInput, people.Organisation])

	router.Register(organisationsAPI, "UpdateOrganisation",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update organisation",
			Errors:  []int{400, 500},
		}, controllers.UpdateByCodeHandler[people.OrganisationUpdate])

	router.Register(organisationsAPI, "DeleteOrganisation",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete organisation",
			Errors:  []int{400, 500},
		},
		controllers.DeleteByCodeHandler(people.DeleteOrganisation))
}
