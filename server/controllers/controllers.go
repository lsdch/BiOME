package controllers

import (
	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

type UUIDInput = struct {
	UUID edgedb.UUID `uri:"id" binding:"required"`
}

type CodeInput = struct {
	Code string `uri:"code" binding:"required"`
}

type IDParser[ID any] func(ctx *gin.Context) (ID, error)

var ParseUUIDfromURI IDParser[edgedb.UUID] = func(ctx *gin.Context) (edgedb.UUID, error) {
	var id UUIDInput
	if err := ctx.BindUri(&id); err != nil {
		ctx.Error(err)
	}
	return id.UUID, nil
}

var ParseCodeURI IDParser[string] = func(ctx *gin.Context) (string, error) {
	var code CodeInput
	if err := ctx.BindUri(&code); err != nil {
		ctx.Error(err)
	}
	return code.Code, nil
}
