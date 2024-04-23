package person

import (
	"darco/proto/controllers"
	"darco/proto/router"
	"net/http"

	"darco/proto/models/people"
	_ "darco/proto/models/validations"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(r router.Router) {
	personsAPI := r.RouteGroup("/persons").
		WithTags([]string{"People", "Person"})

	router.Register(personsAPI, "ListPersons",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List persons",
			Errors:  []int{500},
		},
		controllers.ListHandler(people.ListPersons),
	)

	router.Register(personsAPI, "CreatePerson",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create person",
			Errors:  []int{400, 500},
		},
		controllers.CreateHandler[people.PersonInput, people.Person])

	router.Register(personsAPI, "UpdatePerson",
		huma.Operation{
			Path:    "/{id}",
			Method:  http.MethodPatch,
			Summary: "Update person",
			Errors:  []int{400, 500},
		}, controllers.UpdateByIDHandler[people.PersonUpdate](people.FindPerson))

	router.Register(personsAPI, "DeletePerson",
		huma.Operation{
			Path:    "/{id}",
			Method:  http.MethodDelete,
			Summary: "Delete person",
			Errors:  []int{400, 500},
		},
		controllers.DeleteByIDHandler(people.DeletePerson))
}
