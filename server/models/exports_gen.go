// This file is auto-generated *DO NOT EDIT*

package models


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var ExportFormatValues = []ExportFormat{
	ExportFormatCSV,
	ExportFormatJSON,
	ExportFormatDWC,
}

// Register enum in OpenAPI specification
func (u ExportFormat) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["ExportFormat"] == nil {
    schemaRef := r.Schema(reflect.TypeFor[string](), true, "ExportFormat")
    schemaRef.Title = "ExportFormat"
    for _, v := range ExportFormatValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["ExportFormat"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/ExportFormat"}
}

func (m *ExportFormat) Fake(f *gofakeit.Faker) (any, error) {
	return string(ExportFormatValues[f.IntN(len(ExportFormatValues) - 1)]), nil
}