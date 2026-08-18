// This file is auto-generated *DO NOT EDIT*

package models


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var OccurrenceSortKeyValues = []OccurrenceSortKey{
	OccurrenceSortKeyCode,
	OccurrenceSortKeySiteName,
	OccurrenceSortKeySiteCode,
	OccurrenceSortKeyEventDate,
	OccurrenceSortKeyTaxonName,
	OccurrenceSortKeyIdentifiedOn,
	OccurrenceSortKeyCreatedAt,
	OccurrenceSortKeyUpdatedAt,
}

// Register enum in OpenAPI specification
func (u OccurrenceSortKey) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["OccurrenceSortKey"] == nil {
    schemaRef := r.Schema(reflect.TypeFor[string](), true, "OccurrenceSortKey")
    schemaRef.Title = "OccurrenceSortKey"
    for _, v := range OccurrenceSortKeyValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["OccurrenceSortKey"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/OccurrenceSortKey"}
}

func (m *OccurrenceSortKey) Fake(f *gofakeit.Faker) (any, error) {
	return string(OccurrenceSortKeyValues[f.IntN(len(OccurrenceSortKeyValues) - 1)]), nil
}