// This file is auto-generated *DO NOT EDIT*

package models


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var SortOrderValues = []SortOrder{
	SortAsc,
	SortDesc,
}

// Register enum in OpenAPI specification
func (u SortOrder) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["SortOrder"] == nil {
    schemaRef := r.Schema(reflect.TypeFor[string](), true, "SortOrder")
    schemaRef.Title = "SortOrder"
    for _, v := range SortOrderValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["SortOrder"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/SortOrder"}
}

func (m *SortOrder) Fake(f *gofakeit.Faker) (any, error) {
	return string(SortOrderValues[f.IntN(len(SortOrderValues) - 1)]), nil
}