// This file is auto-generated *DO NOT EDIT*

package app_errors


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var ErrorCodeValues = []ErrorCode{
	ErrorCodeInconsistentTaxa,
	ErrorCodeInvalidCSVRow,
}

// Register enum in OpenAPI specification
func (u ErrorCode) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["ErrorCode"] == nil {
    schemaRef := r.Schema(reflect.TypeFor[string](), true, "ErrorCode")
    schemaRef.Title = "ErrorCode"
    for _, v := range ErrorCodeValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["ErrorCode"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/ErrorCode"}
}

func (m *ErrorCode) Fake(f *gofakeit.Faker) (any, error) {
	return string(ErrorCodeValues[f.IntN(len(ErrorCodeValues) - 1)]), nil
}



var ErrorCategoryValues = []ErrorCategory{
	ErrorCategoryImport,
}

// Register enum in OpenAPI specification
func (u ErrorCategory) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["ErrorCategory"] == nil {
    schemaRef := r.Schema(reflect.TypeFor[string](), true, "ErrorCategory")
    schemaRef.Title = "ErrorCategory"
    for _, v := range ErrorCategoryValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["ErrorCategory"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/ErrorCategory"}
}

func (m *ErrorCategory) Fake(f *gofakeit.Faker) (any, error) {
	return string(ErrorCategoryValues[f.IntN(len(ErrorCategoryValues) - 1)]), nil
}