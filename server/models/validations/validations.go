package validations

import (
	"context"
	"darco/proto/models"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func EmailUniqueValidator(fl validator.FieldLevel) bool {
	var exists bool
	query := "select exists people::User filter .email = <str>$0"
	err := models.DB.QuerySingle(context.Background(), query, &exists, fl.Field().String())
	if err != nil {
		logrus.Errorf("Unique email validation query failed: %v", err)
	}
	return exists
}
