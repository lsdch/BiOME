package country

import (
	country "darco/proto/models/location"
	"net/http"

	gin "github.com/gin-gonic/gin"
)

func Setup(ctx *gin.Context) {
	err := country.Setup()
	if err != nil {
		ctx.Error(err)
	} else {
		ctx.Status(http.StatusOK)
	}
}

// List countries
// swagger:route GET /countries/
// @Summary List Countries
// @Tags Location
// @Success 200 {array} country.Country
// @Router /countries/ [get]
func List(ctx *gin.Context) {
	countries, err := country.List()
	if err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, countries)
	}
}
