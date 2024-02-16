package controllers

import (
	"darco/proto/models"
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func UpdateItem[ID any, Item models.Updatable[ID, Updated], Updated any](
	ctx *gin.Context,
	db *edgedb.Client,
	parseID IDParser[ID],
	find models.ItemFinder[edgedb.UUID, Updated],
) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}
	var item Item
	if err = ctx.ShouldBindJSON(&item); err != nil {
		ctx.Error(err)
		return
	}

	uuid, err := item.Update(db, id)
	if err != nil {
		logrus.Errorf("Item update failed: %+v", err)
		ctx.Error(err)
		return
	}
	updated, err := find(db, uuid)
	if err != nil {
		logrus.Errorf("Failed to retrieve updated item: %v", err)
		ctx.Error(err)
		return
	}
	logrus.Infof("Item updated: %+v", updated)
	ctx.JSON(http.StatusOK, updated)
}

func UpdateItemByCode[Item models.Updatable[string, Updated], Updated any](
	ctx *gin.Context,
	db *edgedb.Client,
	find models.ItemFinder[edgedb.UUID, Updated],
) {
	UpdateItem[string, Item, Updated](ctx, db, ParseCodeURI, find)
}

func UpdateItemByUUID[Item models.Updatable[edgedb.UUID, Updated], Updated any](
	ctx *gin.Context,
	db *edgedb.Client,
	find models.ItemFinder[edgedb.UUID, Updated],
) {
	UpdateItem[edgedb.UUID, Item, Updated](ctx, db, ParseUUIDfromURI, find)
}
