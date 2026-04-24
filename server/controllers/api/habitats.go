package api

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func init() {
	habitatsAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/habitats").
			WithTags([]string{"Sampling"})
	}

	router.RegisterSpec(
		habitatsAPI,
		"ListHabitatGroups",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List habitats",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListHabitatGroups),
	)

	router.RegisterSpec(
		habitatsAPI,
		"CreateHabitatGroup",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create habitat group",
		},
		controllers.CreateHandler[occurrence.HabitatGroupInput])

	router.RegisterSpec(
		habitatsAPI,
		"DeleteHabitatGroup",
		huma.Operation{
			Path:    "/{label}",
			Method:  http.MethodDelete,
			Summary: "Delete habitat group",
		},
		controllers.DeleteByLabelHandler(occurrence.DeleteHabitatGroup))

	router.RegisterSpec(
		habitatsAPI,
		"UpdateHabitatGroup",
		huma.Operation{
			Path:    "/{label}",
			Method:  http.MethodPatch,
			Summary: "Update habitat group",
		},
		controllers.UpdateHandler[*struct {
			resolvers.AccessRestricted[resolvers.Maintainer]
			controllers.LabelInput
			controllers.UpdateInput[occurrence.HabitatGroupUpdate, string, occurrence.HabitatGroup]
		}])
}
