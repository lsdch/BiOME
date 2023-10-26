package validations

import (
	"context"
	"darco/proto/models"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UniqueValidator struct {
	model           string
	field           string
	edgedb_typecast string
}

func (uq *UniqueValidator) Query() string {
	return fmt.Sprintf("select exists %s filter .%s = <%s>$0", uq.model, uq.field, uq.edgedb_typecast)
}

func (uq *UniqueValidator) Validator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		var exists bool
		err := models.DB.QuerySingle(context.Background(), uq.Query(), &exists, fl.Field().String())
		if err != nil {
			logrus.Errorf("Unique validation query failed: %v with query %s", err, uq.Query())
		}
		return exists
	}
}

var EmailUnique = UniqueValidator{
	model:           "people::User",
	field:           "email",
	edgedb_typecast: "str",
}
