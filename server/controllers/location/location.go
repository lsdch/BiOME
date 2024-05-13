package location

import (
	"context"
	"darco/proto/controllers"
	"darco/proto/db"
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
		}, ListHabitatGroups)

	router.Register(locationAPI, "CreateHabitat",
		huma.Operation{
			Path:    "/habitats",
			Method:  http.MethodPost,
			Summary: "Create habitat",
		}, controllers.CreateHandler[location.HabitatInput])
}

type HabitatGroups struct {
	Body []location.HabitatGroup
}

func ListHabitatGroups(ctx context.Context, input *struct{}) (*HabitatGroups, error) {
	data, err := location.ListHabitatGroups(db.Client())
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to retrieve grouped habitat list", err)
	}
	return &HabitatGroups{data}, nil
}
