package habitats

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(r router.Router) {
	habitatsAPI := r.RouteGroup("/habitats").
		WithTags([]string{"Habitats"})

	router.Register(habitatsAPI, "ListHabitatGroups",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List habitats",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListHabitatGroups))

	router.Register(habitatsAPI, "CreateHabitatGroup",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create habitat group",
		}, controllers.CreateHandler[occurrence.HabitatGroupInput])

	router.Register(habitatsAPI, "DeleteHabitatGroup",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete habitat group",
		}, controllers.DeleteByCodeHandler(occurrence.DeleteHabitatGroup))

	router.Register(habitatsAPI, "UpdateHabitatGroup",
		huma.Operation{
			Path:    "/{label}",
			Method:  http.MethodPatch,
			Summary: "Update habitat group",
		}, controllers.UpdateHandler[*struct {
			resolvers.AccessRestricted[resolvers.Maintainer]
			controllers.LabelInput
			controllers.UpdateInput[occurrence.HabitatGroupUpdate, string, occurrence.HabitatGroup]
		}])

}
