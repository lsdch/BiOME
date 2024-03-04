package validations

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func validateUnique(typecast string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		bindings := ParseEdgeDBBindingsOrDie(fl.Param(), typecast)
		return bindings.UniqueQuery(fl.Field().Interface())
	}
}

func validateUniqueWithBindings(bindings BindingEdgeDB) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return bindings.UniqueQuery(fl.Field().Interface())
	}
}

// Checks that a string property is unique.
// Expects a parameter that specifies the object property in the database.
//
// Format : <module name>::<object name>.<property name>
//
// Example : unique_str=people::User.email
var UniqueStrValidator = CustomValidator{
	Tag:     "unique_str",
	Handler: validateUnique("str"),
	Message: func(fl validator.FieldError) string {
		return fmt.Sprintf("'%s' is already in use", fl.Value())
	},
}

// Checks that an email address is not already in use
var UniqueEmailValidator = CustomValidator{
	Tag: "unique_email",
	Handler: validateUniqueWithBindings(BindingEdgeDB{
		ObjectName:   "people::User",
		PropertyName: "email",
		TypeCast:     "str",
	}),
	Message: func(fl validator.FieldError) string {
		return "An account is already registered with this address"
	},
}

// Checks that a login is not already in use in the database
var UniqueLoginValidator = CustomValidator{
	Tag: "unique_login",
	Handler: validateUniqueWithBindings(BindingEdgeDB{
		ObjectName:   "people::User",
		PropertyName: "login",
		TypeCast:     "str",
	}),
	Message: func(fl validator.FieldError) string {
		return "This login is already used"
	},
}
