package controllers

import (
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

type FetchItemList[Item any] func(db *edgedb.Client) ([]Item, error)

func ListItems[Item any](ctx *gin.Context, db *edgedb.Client, listFn FetchItemList[Item]) {
	items, err := listFn(db)
	if err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, NonNilArray[Item](items))
	}
}
