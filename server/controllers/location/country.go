package country

import (
	"darco/proto/controllers"
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
	countriesAPI := r.RouteGroup("/countries").
		WithTags([]string{"Location", "Countries"})

	router.Register(countriesAPI, "ListCountries",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List countries",
		}, controllers.ListHandler(country.List))
}
