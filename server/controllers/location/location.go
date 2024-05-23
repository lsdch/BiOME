package location

import (
	"context"
	"darco/proto/controllers"
	"darco/proto/models/location"
	"darco/proto/resolvers"
	"darco/proto/router"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
	gin "github.com/gin-gonic/gin"
)

func Setup(ctx *gin.Context, db *edgedb.Client) {
	err := location.Setup(db)
	if err != nil {
		ctx.Error(err)
	} else {
		ctx.Status(http.StatusOK)
	}
}

func RegisterRoutes(r router.Router) {
	locationAPI := r.RouteGroup("/locations").
		WithTags([]string{"Location", "Countries"})

	router.Register(locationAPI, "ListCountries",
		huma.Operation{
			Path:    "/countries",
			Method:  http.MethodGet,
			Summary: "List countries",
		}, controllers.ListHandler(location.List))

	router.Register(locationAPI, "ListHabitats",
		huma.Operation{
			Path:    "/habitats",
			Method:  http.MethodGet,
			Summary: "List habitats",
		}, controllers.ListHandler(location.ListHabitats))

	router.Register(locationAPI, "ListHabitatGroups",
		huma.Operation{
			Path:    "/habitat-groups",
			Method:  http.MethodGet,
			Summary: "List habitat groups",
		}, controllers.ListHandler(location.ListHabitatGroups))

	router.Register(locationAPI, "CreateHabitat",
		huma.Operation{
			Path:    "/habitats",
			Method:  http.MethodPost,
			Summary: "Create habitat",
		}, controllers.CreateHandler[location.HabitatInput])

	router.Register(locationAPI, "CreateHabitatGroup",
		huma.Operation{
			Path:    "/habitat-groups",
			Method:  http.MethodPost,
			Summary: "Create habitat group",
		}, controllers.CreateHandler[location.HabitatGroupInput])

	router.Register(locationAPI, "DeleteHabitatGroup",
		huma.Operation{
			Path:    "/habitat-groups/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete habitat group",
		}, controllers.DeleteByCodeHandler(location.DeleteHabitatGroup))

	router.Register(locationAPI, "UpdateHabitatGroup",
		huma.Operation{
			Path:    "/habitat-groups/{code}",
			Method:  http.MethodPatch,
			Summary: "Update habitat group",
		}, controllers.UpdateByCodeHandler[location.HabitatGroupUpdate](location.FindHabitatGroup))

	router.Register(locationAPI, "setHabitatGroupDepends",
		huma.Operation{
			Path:    "/habitat-groups/{code}/set-depends",
			Method:  http.MethodPost,
			Summary: "Set dependency of habitat group to habitat",
		}, SetHabitatGroupDepends)
}

type SetHabitatGroupDependsInput struct {
	resolvers.AccessRestricted[resolvers.Admin]
	controllers.CodeInput
	Depends string `query:"set-depends"`
}
type SetHabitatGroupDependsOutput struct {
	Body location.HabitatGroup
}

func SetHabitatGroupDepends(ctx context.Context, input *SetHabitatGroupDependsInput) (*SetHabitatGroupDependsOutput, error) {
	group, err := location.SetHabitatGroupParent(input.DB(), input.Code, input.Depends)
	if err != nil {
		huma.Error500InternalServerError("Failed to set habitat group dependency, err")
	}
	return &SetHabitatGroupDependsOutput{Body: *group}, nil
}
