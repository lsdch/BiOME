package validations

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

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

var Validators = []CustomValidator{
	UniqueEmailValidator,
	UniqueLoginValidator,
	LoginValidator,
	ExistAllValidator,
	ExistValidator,
}

func RegisterValidators(engine *validator.Validate) {
	for _, validator := range Validators {
		engine.RegisterValidation(string(validator.tag), validator.handler)
	}
}
