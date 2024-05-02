package country

import (
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

	router.Register(locationAPI, "CreateHabitat",
		huma.Operation{
			Path:    "/habitats",
			Method:  http.MethodPost,
			Summary: "Create habitat",
		}, controllers.CreateHandler[location.HabitatInput])
}
