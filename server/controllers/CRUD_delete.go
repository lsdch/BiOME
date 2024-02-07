package controllers

import (
	"darco/proto/models"
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

type ItemDelete[ID any, Item models.Updatable[Item]] func(db *edgedb.Client, id ID) (Item, error)

func delete[ID any, Item models.Updatable[Item]](
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
	ctx.JSON(http.StatusOK, deleted)
}

func DeleteByCode[Item models.Updatable[Item]](
	ctx *gin.Context,
	db *edgedb.Client,
	deleteItem ItemDelete[string, Item],
) {
	delete[string, Item](ctx, db, deleteItem, ParseCodeURI)
}

func DeleteByID[Item models.Updatable[Item]](
	ctx *gin.Context,
	db *edgedb.Client,
	deleteItem ItemDelete[edgedb.UUID, Item],
) {
	delete[edgedb.UUID, Item](ctx, db, deleteItem, ParseUUIDfromURI)
}
