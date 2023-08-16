package taxonomy

import (
	"darco/proto/models/taxonomy"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAnchors(ctx *gin.Context) {
	anchors, err := taxonomy.GetAnchorTaxa()
	if err != nil {
		ctx.Error(err).SetMeta(gin.H{
			"msg": "Failed to fetch taxonomy data",
		})
	} else {
		ctx.JSON(http.StatusOK, anchors)
	}
}
