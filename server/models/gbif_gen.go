// This file is auto-generated *DO NOT EDIT*

package models


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var TaxonGBIFPriorityValues = []TaxonGBIFPriority{
	TaxonGBIFPriorityExactAccepted,
	TaxonGBIFPriorityExactNonAccepted,
	TaxonGBIFPriorityNonExact,
}

// Register enum in OpenAPI specification
func (u TaxonGBIFPriority) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["TaxonGBIFPriority"] == nil {
    schemaRef := r.Schema(reflect.TypeFor[int32](), true, "TaxonGBIFPriority")
    schemaRef.Title = "TaxonGBIFPriority"
    for _, v := range TaxonGBIFPriorityValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["TaxonGBIFPriority"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/TaxonGBIFPriority"}
}

func (m *TaxonGBIFPriority) Fake(f *gofakeit.Faker) (any, error) {
	return string(TaxonGBIFPriorityValues[f.IntN(len(TaxonGBIFPriorityValues) - 1)]), nil
}