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

// Error handling middleware
func ErrorHandler(c *gin.Context) {
	c.Next() // execute all the handlers

	// at this point, all the handlers finished. Let's read the errors!
	// in this example we only will use the **last error typed as public**
	// but you could iterate over all them since c.Errors is a slice!
	err := c.Errors.Last()
	if err == nil {
		return
	}

	var dbErr edgedb.Error
	if errors.As(err, &dbErr) && dbErr.Category(edgedb.NoDataError) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		apiErr := validations.ValidationErrorsByField(validations.InputValidationErrors(validationErr))
		logrus.Debugf("Validation error: %v", apiErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, apiErr)
		return
	}

	var manualValidationErr validations.InputValidationError
	if errors.As(err, &manualValidationErr) {
		logrus.Debugf("Validation error: %v", manualValidationErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if err.Meta != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Meta)
		return
	}

	c.JSON(int(err.Type), err.Err.Error())
}
