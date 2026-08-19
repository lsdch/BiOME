// This file is auto-generated *DO NOT EDIT*

package csvmodels


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var CSVDelimiterValues = []CSVDelimiter{
	CSVDelimiterComma,
	CSVDelimiterSemicolon,
	CSVDelimiterTab,
}

// Register enum in OpenAPI specification
func (u CSVDelimiter) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["CSVDelimiter"] == nil {
    schemaRef := r.Schema(reflect.TypeFor[string](), true, "CSVDelimiter")
    schemaRef.Title = "CSVDelimiter"
    for _, v := range CSVDelimiterValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["CSVDelimiter"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/CSVDelimiter"}
}

func (m *CSVDelimiter) Fake(f *gofakeit.Faker) (any, error) {
	return string(CSVDelimiterValues[f.IntN(len(CSVDelimiterValues) - 1)]), nil
}



var CSVQuoteCharValues = []CSVQuoteChar{
	CSVQuoteCharDouble,
	CSVQuoteCharSingle,
}

// Register enum in OpenAPI specification
func (u CSVQuoteChar) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["CSVQuoteChar"] == nil {
    schemaRef := r.Schema(reflect.TypeFor[string](), true, "CSVQuoteChar")
    schemaRef.Title = "CSVQuoteChar"
    for _, v := range CSVQuoteCharValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["CSVQuoteChar"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/CSVQuoteChar"}
}

func (m *CSVQuoteChar) Fake(f *gofakeit.Faker) (any, error) {
	return string(CSVQuoteCharValues[f.IntN(len(CSVQuoteCharValues) - 1)]), nil
}