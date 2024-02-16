package validations

import (
	"context"
	"darco/proto/db"
	"encoding/json"
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
	return fmt.Sprintf("select exists %s filter .%s = <%s><json>$0", uq.model, uq.field, uq.edgedb_typecast)
}

func (uq *UniqueValidator) Validator() validator.Func {
	return func(fl validator.FieldLevel) (exists bool) {
		args, _ := json.Marshal(fl.Field().String())
		err := db.Client().QuerySingle(context.Background(), uq.Query(), &exists, args)
		if err != nil {
			logrus.Errorf("Unique validation query failed: %v with query %s", err, uq.Query())
		}
		return !exists
	}
}

var emailUnique = UniqueValidator{
	model:           "people::User",
	field:           "email",
	edgedb_typecast: "str",
}

var loginUnique = UniqueValidator{
	model:           "people::User",
	field:           "login",
	edgedb_typecast: "str",
}

// Checks that an email address is not already in use
var UniqueEmailValidator = CustomValidator{
	tag:     "unique_email",
	handler: emailUnique.Validator(),
	message: func(fl validator.FieldError) string {
		return "An account is already registered with this address"
	},
}

// Checks that a login is not already in use in the database
var UniqueLoginValidator = CustomValidator{
	tag:     "unique_login",
	handler: loginUnique.Validator(),
	message: func(fl validator.FieldError) string {
		return "This login is already used"
	},
}
