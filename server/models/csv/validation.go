package csvmodels

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-playground/validator/v10"
	"github.com/lsdch/biome/lib/app_errors"
)

type CSVParseError struct {
	RowNumber int32
	Err       error
}

func (e *CSVParseError) AppError() *app_errors.AppError {
	return &app_errors.AppError{
		Category: app_errors.ErrorCategoryImport,
		Code:     app_errors.ErrorCodeInvalidCSVRow,
		ErrorModel: huma.ErrorModel{
			Status: http.StatusUnprocessableEntity,
			Title:  "Error parsing CSV",
			Detail: fmt.Sprintf("Failed parsing row %d", e.RowNumber),
			Errors: []*huma.ErrorDetail{
				{Value: e.Err, Location: strconv.Itoa(int(e.RowNumber)), Message: e.Err.Error()},
			},
		},
	}
}

func (e *CSVParseError) Error() string {
	return fmt.Sprintf("error parsing CSV at row %d: %v", e.RowNumber, e.Err)
}

type CSVImportParams struct {
	Reader    io.Reader
	Separator rune
}

type RowValidator interface {
	RowNumber() int32
	Validate(v *validator.Validate) error
}

func ValidateRows[T RowValidator](rows []T, v *validator.Validate) error {
	validationErrors := BatchValidationErrors{Errors: make([]*RowValidationErrors, 0, len(rows))}
	for rowNum, row := range rows {
		if err := row.Validate(v); err != nil {
			var validationErrs validator.ValidationErrors
			if errors.As(err, &validationErrs) {
				validationErrors.Errors = append(validationErrors.Errors, &RowValidationErrors{
					RowNumber: int32(rowNum + 2),
					Errors:    validationErrs,
				})
			}
		}
	}
	if len(validationErrors.Errors) > 0 {
		return &validationErrors
	}
	return nil
}

type RowValidationErrors struct {
	RowNumber int32
	Errors    validator.ValidationErrors
}

func (e *RowValidationErrors) Error() string {
	return fmt.Sprintf("row %d: %v", e.RowNumber, e.Errors)
}

func (e *RowValidationErrors) Unwrap() error {
	return e.Errors
}

type BatchValidationErrors struct {
	Errors []*RowValidationErrors
}

func (e *BatchValidationErrors) Error() string {
	var sb strings.Builder
	sb.WriteString("validation errors:\n")
	for _, rowErr := range e.Errors {
		sb.WriteString(rowErr.Error())
	}
	return sb.String()
}

func (e *BatchValidationErrors) Unwrap() []error {
	errs := make([]error, len(e.Errors))

	for i, rowErr := range e.Errors {
		errs[i] = rowErr
	}

	return errs
}
