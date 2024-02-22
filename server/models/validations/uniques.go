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
	Tag:     "unique_email",
	Handler: emailUnique.Validator(),
	Message: func(fl validator.FieldError) string {
		return "An account is already registered with this address"
	},
}

// Checks that a login is not already in use in the database
var UniqueLoginValidator = CustomValidator{
	Tag:     "unique_login",
	Handler: loginUnique.Validator(),
	Message: func(fl validator.FieldError) string {
		return "This login is already used"
	},
}

func validateUnique(typecast string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		bindings := ParseEdgeDBBindings(fl.Param(), "str")
		return bindings.UniqueQuery(fl.Field().Interface())
	}
}

func validateUniqueWithBindings(bindings BindingEdgeDB) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return bindings.UniqueQuery(fl.Field().Interface())
	}
}

var UniqueStrValidator = CustomValidator{
	Tag:     "unique_str",
	Handler: validateUnique("str"),
	Message: func(fl validator.FieldError) string {
		return fmt.Sprintf("'%s' is already in use", fl.Value())
	},
}
