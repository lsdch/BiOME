package accounts

import "darco/proto/models/validations"

func weakPasswordError(field string, err error) validations.InputValidationError {
	return validations.InputValidationError{
		Field:     field,
		Message:   "Password is too weak",
		ErrString: err.Error(),
	}
}
