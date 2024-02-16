package controllers

import (
	"darco/proto/models"
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

func CreateItem[Item models.Creatable[Created], Created any](
	ctx *gin.Context,
	db *edgedb.Client,
) {
	var item Item
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.Error(err)
		return
	}
	created, err := item.Create(db)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusAccepted, created)
}
