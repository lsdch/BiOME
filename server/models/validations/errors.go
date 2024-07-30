package validations

import (
	"fmt"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-playground/validator/v10"
)

// A slice of validation errors which allows bulk manipulation such as updating error locations
type ValidationErrors []*huma.ErrorDetail

func (e ValidationErrors) Errors() []error {
	var errors = make([]error, len(e))
	for _, err := range e {
		errors = append(errors, err)
	}
	return errors
}

// Prepends an error location string to all errors in the slice.
// Don't forget '.' separators when needed.
func (e ValidationErrors) WithLocation(location string) []error {
	for _, err := range []*huma.ErrorDetail(e) {
		err.Location = location + err.Location
	}
	return e.Errors()
}

// Structured validation error built from go-playground/validator to be consumed by API clients
// @description A validation error to be fixed in the input
type InputValidationError struct {
	Field     string `json:"field" example:"age" binding:"required"`
	Value     any    `json:"value" binding:"required"`
	Tag       string `json:"tag" example:"min"`
	Param     string `json:"param" example:"0"`
	Message   string `json:"message" example:"Must be a positive number" binding:"required"`
	ErrString string `json:"error" example:"Key: 'Person.age' Error:Field validation for 'age' failed on the 'min' tag"`
} //@name InputValidationError

func (vErr InputValidationError) Error() string {
	return vErr.ErrString
}

// Transforms [validator.ValidationErrors] to a slice of richer [InputValidationError].
func InputValidationErrors(errors validator.ValidationErrors) []InputValidationError {
	out := make([]InputValidationError, len(errors))
	for i, err := range errors {
		out[i] = InputValidationError{
			Field:     err.Field(),
			Value:     err.Value(),
			Tag:       err.Tag(),
			Param:     err.Param(),
			ErrString: err.Error(),
			Message:   fieldErrorMsg(err),
		}
	}
	return out
}

// Validation errors indexed by JSON field name
type FieldErrors map[string][]InputValidationError // @name FieldErrors

// Indexes validation errors by JSON field name
func ValidationErrorsByField(errors []InputValidationError) FieldErrors {
	indexed_errors := make(map[string][]InputValidationError)
	for _, err := range errors {
		indexed_errors[err.Field] = append(indexed_errors[err.Field], err)
	}
	return indexed_errors
}

// Custom error messages for some built-in validation tags
func fieldErrorMsg(err validator.FieldError) string {
	switch ValidationTag(err.Tag()) {
	case Required:
		return "This field is required"
	case Email:
		return "Invalid email address"
	case Minimum:
		switch err.Kind() {
		case reflect.String:
			return fmt.Sprintf("Required length is at least %s characters", err.Param())
		case reflect.Slice, reflect.Map, reflect.Array:
			return fmt.Sprintf("At least %s values are required", err.Param())
		default:
			return fmt.Sprintf("Minimum value is %s", err.Param())
		}
	case Maximum:
		switch err.Kind() {
		case reflect.String:
			return fmt.Sprintf("Maximum length is %s characters", err.Param())
		case reflect.Slice, reflect.Map, reflect.Array:
			return fmt.Sprintf("At most %s values can be provided", err.Param())
		default:
			return fmt.Sprintf("Maximum value is %s", err.Param())
		}
	case NotEqual:
		switch err.Kind() {
		case reflect.Slice, reflect.Map, reflect.Array:
			return fmt.Sprintf("Number of items must be different from %s", err.Param())
		default:
			return fmt.Sprintf("Must be different from %s", err.Param())
		}
	case GreaterThan:
		return fmt.Sprintf("Must be greater than %s", err.Param())
	case GreaterOrEqual:
		return fmt.Sprintf("Must be greater or equal to %s", err.Param())
	case Alpha:
		return "Only alphabetic characters allowed"
	case AlphaUnicode:
		return "Only alphabetic characters allowed"
	case Numeric:
		return "Only numbers allowed"
	case StartsWith:
		return fmt.Sprintf("Must start with \"%s\"", err.Param())
	case EndsWith:
		return fmt.Sprintf("Must end with \"%s\"", err.Param())
	}

	for _, validator := range validators {
		if ValidationTag(err.Tag()) == validator.Tag {
			return validator.Message(err)
		}
	}

	for _, customTag := range customTags {
		if err.Tag() == customTag.Alias {
			return customTag.Message
		}
	}

	return "Invalid value"
}
