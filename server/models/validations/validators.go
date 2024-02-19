package validations

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Tag     ValidationTag
	Handler validator.Func
	Message func(fl validator.FieldError) string
}

var loginRegex = regexp.MustCompile("^[a-zA-Z0-9.]+$")
var LoginValidator = CustomValidator{
	Tag: "login",
	Handler: func(fl validator.FieldLevel) bool {
		return loginRegex.MatchString(fl.Field().String())
	},
	Message: func(fl validator.FieldError) string {
		return "Only alphanumeric and '.' characters allowed"
	},
}

var validators = []CustomValidator{
	UniqueEmailValidator,
	UniqueLoginValidator,
	LoginValidator,
	ExistAllValidator,
	ExistValidator,
}

func RegisterCustomValidator(v CustomValidator) error {
	validators = append(validators, v)
	return nil
}

type CustomTag struct {
	Alias   string
	Tags    string
	Message string
}

var customTags = []CustomTag{
	{Alias: "nullalpha", Tags: "eq=|alpha",
		Message: "Only alphabetic characters allowed"},
	{Alias: "nullalphanum", Tags: "eq=|alphanum",
		Message: "Only alphanumeric characters allowed"},
	{Alias: "nullalphaunicode", Tags: "eq=|alphaunicode",
		Message: "Only alphabetic characters allowed"},
	{Alias: "nullalphanumunicode", Tags: "eq=|alphanumunicode",
		Message: "Only alphabetic characters allowed"},
	{Alias: "nullemail", Tags: "eq=|email",
		Message: "Only alphabetic characters allowed"},
}

func RegisterValidators(engine *validator.Validate) {
	for _, validator := range validators {
		engine.RegisterValidation(string(validator.Tag), validator.Handler)
	}
	for _, tag := range customTags {
		engine.RegisterAlias(tag.Alias, tag.Tags)
	}
}
