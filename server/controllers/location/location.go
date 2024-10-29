package location

import (
	"darco/proto/controllers"
	"darco/proto/models/location"
	"darco/proto/resolvers"
	"darco/proto/router"
	"fmt"
	"net/http"

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

	habitatsAPI := r.RouteGroup("/habitats").
		WithTags([]string{"Location", "Habitats"})

	router.Register(habitatsAPI, "ListHabitatGroups",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List habitats",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](location.ListHabitatGroups))

	router.Register(habitatsAPI, "CreateHabitatGroup",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create habitat group",
		}, controllers.CreateHandler[location.HabitatGroupInput])

	router.Register(habitatsAPI, "DeleteHabitatGroup",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete habitat group",
		}, controllers.DeleteByCodeHandler(location.DeleteHabitatGroup))

	router.Register(habitatsAPI, "UpdateHabitatGroup",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update habitat group",
		}, controllers.UpdateByCodeHandler[location.HabitatGroupUpdate])

}

type HabitatGroupInput[R resolvers.RoleSpecifier] struct {
	resolvers.AccessRestricted[R]
	controllers.CodeInput
	HabitatGroup location.HabitatGroup
}
type HabitatGroupOutput struct {
	Body location.HabitatGroup
}

func (i *HabitatGroupInput[R]) Resolve(ctx huma.Context) []error {
	group, err := location.FindHabitatGroup(i.DB(), i.Code)
	if err != nil {
		return []error{huma.Error404NotFound(fmt.Sprintf("Habitat group '%s' does not exist", i.Code))}
	}
	i.HabitatGroup = *group
	return nil
}
