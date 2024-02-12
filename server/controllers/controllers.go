package controllers

import (
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

type CodeInput = struct {
	Code string `uri:"code" binding:"required"`
}

type IDParser[ID any] func(ctx *gin.Context) (ID, error)

var ParseUUIDfromURI IDParser[edgedb.UUID] = func(ctx *gin.Context) (edgedb.UUID, error) {
	var strUUID = ctx.Param("id")
	uuid, err := edgedb.ParseUUID(strUUID)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
	}
	return uuid, err
}

var ParseCodeURI IDParser[string] = func(ctx *gin.Context) (string, error) {
	var code CodeInput
	if err := ctx.BindUri(&code); err != nil {
		ctx.Error(err)
	}
	return code.Code, nil
}

func NonNilArray[T any](array []T) []T {
	if array != nil {
		return array
	} else {
		return []T{}
	}
}
