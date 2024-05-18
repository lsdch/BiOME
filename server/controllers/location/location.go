package location

import (
	"context"
	"darco/proto/controllers"
	"darco/proto/models/location"
	country "darco/proto/models/location"
	"darco/proto/router"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
	gin "github.com/gin-gonic/gin"
)

func Setup(ctx *gin.Context, db *edgedb.Client) {
	err := country.Setup(db)
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
		}, controllers.ListHandler(country.List))

	router.Register(locationAPI, "ListHabitats",
		huma.Operation{
			Path:    "/habitats",
			Method:  http.MethodGet,
			Summary: "List habitats",
		}, controllers.ListHandler(country.ListHabitats))

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
}

type SetHabitatGroupParentInput struct {
	controllers.CodeInput
	Body struct {
		Parent location.Habitat
	}
}
type SetHabitatGroupParentOutput struct{}

func SetHabitatGroupParent(ctx context.Context, input *SetHabitatGroupParentInput) (*SetHabitatGroupParentOutput, error) {
	return &SetHabitatGroupParentOutput{}, nil
}
