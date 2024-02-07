package controllers

import (
	"darco/proto/models"
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ParamsUpdate[ID any, Item models.Updatable[Item]](ctx *gin.Context, db *edgedb.Client, parseIdentifier IDParser[ID], find models.ItemFinder[ID, Item]) (*Item, error) {
	id, err := parseIdentifier(ctx)
	if err != nil {
		return nil, err
	}
	item, err := find(db, id)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return nil, err
	}
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.Error(err)
		return nil, err
	}
	return &item, nil
}

type UpdateParamBinder[ID any, Item models.Updatable[Item]] func(ctx *gin.Context, db *edgedb.Client, find models.ItemFinder[ID, Item]) (*Item, error)

func BindUpdateByID[Item models.Updatable[Item]](ctx *gin.Context, db *edgedb.Client, find models.ItemFinder[edgedb.UUID, Item]) (*Item, error) {
	return ParamsUpdate[edgedb.UUID, Item](ctx, db, ParseUUIDfromURI, find)
}

func BindUpdateByCode[Item models.Updatable[Item]](ctx *gin.Context, db *edgedb.Client, find models.ItemFinder[string, Item]) (*Item, error) {
	return ParamsUpdate[string, Item](ctx, db, ParseCodeURI, find)
}

func UpdateItem[ID any, Item models.Updatable[Item]](
	ctx *gin.Context,
	db *edgedb.Client,
	find models.ItemFinder[ID, Item],
	bindParams UpdateParamBinder[ID, Item],
) {
	item, err := bindParams(ctx, db, find)
	if err != nil {
		return
	}

	updated, err := (*item).Update(db)
	if err != nil {
		logrus.Errorf("Item update failed : %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

func UpdateByCode[Item models.Updatable[Item]](
	ctx *gin.Context,
	db *edgedb.Client,
	find models.ItemFinder[string, Item],
) {
	UpdateItem[string, Item](ctx, db, find, BindUpdateByCode[Item])
}

func UpdateByID[Item models.Updatable[Item]](
	ctx *gin.Context,
	db *edgedb.Client,
	find models.ItemFinder[edgedb.UUID, Item],
) {
	UpdateItem[edgedb.UUID, Item](ctx, db, find, BindUpdateByID[Item])
}
