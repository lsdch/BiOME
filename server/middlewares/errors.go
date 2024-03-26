package middlewares

import (
	"darco/proto/models/validations"
	"errors"
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func handleNoData(ctx *gin.Context, err error) bool {
	var dbErr edgedb.Error
	if errors.As(err, &dbErr) && dbErr.Category(edgedb.NoDataError) {
		ctx.AbortWithStatus(http.StatusNotFound)
		return true
	}
	return false
}

func handleValidationErrors(ctx *gin.Context, err error) bool {
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		apiErr := validations.ValidationErrorsByField(validations.InputValidationErrors(validationErr))
		logrus.Debugf("Validation error: %+v", apiErr)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apiErr)
		return true
	}
	return false
}

func handleManualValidationErrors(ctx *gin.Context, err error) bool {
	var manualValidationErr validations.InputValidationError
	if errors.As(err, &manualValidationErr) {
		apiErr := validations.ValidationErrorsByField(
			[]validations.InputValidationError{manualValidationErr},
		)
		logrus.Debugf("Validation error: %+v", apiErr)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apiErr)
		return true
	}
	return false
}

// Error handling middleware
func ErrorHandler(c *gin.Context) {
	c.Next() // execute all the handlers

	err := c.Errors.Last()
	if err == nil {
		return
	}

	var handled = handleNoData(c, err) ||
		handleValidationErrors(c, err) ||
		handleManualValidationErrors(c, err)

	if !handled {
		if err.Meta != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Meta)
			return
		}

		c.JSON(http.StatusInternalServerError, err.Err.Error())
	}
}
