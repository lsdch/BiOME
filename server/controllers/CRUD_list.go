package controllers

import (
	"darco/proto/models"
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

type FetchItemList[Item models.Updatable[Item]] func(db *edgedb.Client) ([]Item, error)

func ListItems[Item models.Updatable[Item]](ctx *gin.Context, db *edgedb.Client, listFn FetchItemList[Item]) {
	items, err := listFn(db)
	if err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, NonNilArray[Item](items))
	}
}
