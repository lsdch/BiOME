package validations

import (
	"context"
	"darco/proto/models"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UniqueValidator struct {
	model           string
	field           string
	edgedb_typecast string
}

func (uq *UniqueValidator) Query() string {
	return fmt.Sprintf("select exists (select %s filter .%s = <%s>$0)", uq.model, uq.field, uq.edgedb_typecast)
}

func (uq *UniqueValidator) Validator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		var exists bool
		err := models.DB().QuerySingle(context.Background(), uq.Query(), &exists, fl.Field().String())
		if err != nil {
			logrus.Errorf("Unique validation query failed: %v with query %s", err, uq.Query())
		}
		return !exists
	}
}

var EmailUnique = UniqueValidator{
	model:           "people::User",
	field:           "email",
	edgedb_typecast: "str",
}

var LoginUnique = UniqueValidator{
	model:           "people::User",
	field:           "login",
	edgedb_typecast: "str",
}

type CustomValidator struct {
	tag     ValidationTag
	handler validator.Func
	message func(fl validator.FieldError) string
}

var loginRegex = regexp.MustCompile("^[a-zA-Z0-9.]+$")
var LoginValidator = CustomValidator{
	tag: "login",
	handler: func(fl validator.FieldLevel) bool {
		return loginRegex.MatchString(fl.Field().String())
	},
	message: func(fl validator.FieldError) string {
		return "Only alphanumeric and '.' characters allowed"
	},
}

// Checks that an email address is not already in use
var UniqueEmailValidator = CustomValidator{
	tag:     "unique_email",
	handler: EmailUnique.Validator(),
	message: func(fl validator.FieldError) string {
		return "An account is already registered with this address"
	},
}

// Checks that a login is not already in use in the database
var UniqueLoginValidator = CustomValidator{
	tag:     "unique_login",
	handler: LoginUnique.Validator(),
	message: func(fl validator.FieldError) string {
		return "This login is already used"
	},
}

var validators = []CustomValidator{UniqueEmailValidator, UniqueLoginValidator, LoginValidator}

func RegisterValidators(engine *validator.Validate) {
	for _, validator := range validators {
		engine.RegisterValidation(string(validator.tag), validator.handler)
	}
}
