package validations

import (
	"github.com/danielgtaylor/huma/v2"
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
