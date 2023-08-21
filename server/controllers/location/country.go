package country

import (
	country "darco/proto/models/location"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(ctx *gin.Context) {
	err := country.Setup()
	if err != nil {
		ctx.Error(err)
	} else {
		ctx.Status(http.StatusOK)
	}
}
