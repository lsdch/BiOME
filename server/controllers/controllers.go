package controllers

import (
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

// Functions that parse an identifier from Gin context
type IDParser[ID any] func(ctx *gin.Context) (ID, error)

// Retrieves an UUID from URI parameters, using the 'id' key
var ParseUUIDfromURI IDParser[edgedb.UUID] = func(ctx *gin.Context) (edgedb.UUID, error) {
	var strUUID = ctx.Param("id")
	uuid, err := edgedb.ParseUUID(strUUID)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
	}
	return uuid, err
}

// Retrieves a string identifier from URI parameters, using the 'code' key
var ParseCodeURI IDParser[string] = func(ctx *gin.Context) (string, error) {
	return ctx.Param("code"), nil
}

// Replaces nil array with empty array
func NonNilArray[T any](array []T) []T {
	if array != nil {
		return array
	} else {
		return []T{}
	}
}
