package controllers

import (
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ItemDelete[ID any, Item any] func(db *edgedb.Client, id ID) (Item, error)

func delete[ID any, Item any](
	ctx *gin.Context,
	db *edgedb.Client,
	delete ItemDelete[ID, Item],
	bindID IDParser[ID],
) {
	id, err := bindID(ctx)
	if err != nil {
		return
	}
	deleted, err := delete(db, id)
	if err != nil {
		ctx.Error(err)
		return
	}
	logrus.Debugf("Deleted item %+v", deleted)
	ctx.JSON(http.StatusOK, deleted)
}

func DeleteByCode[Item any](
	ctx *gin.Context,
	db *edgedb.Client,
	deleteItem ItemDelete[string, Item],
) {
	delete[string, Item](ctx, db, deleteItem, ParseCodeURI)
}

func DeleteByID[Item any](
	ctx *gin.Context,
	db *edgedb.Client,
	deleteItem ItemDelete[edgedb.UUID, Item],
) {
	delete[edgedb.UUID, Item](ctx, db, deleteItem, ParseUUIDfromURI)
}
