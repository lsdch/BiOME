// This file is auto-generated *DO NOT EDIT*

package occurrence


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var OccurrenceCategoryValues = []OccurrenceCategory{
	Internal,
	External,
}

// Register enum in OpenAPI specification
func (u OccurrenceCategory) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["OccurrenceCategory"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "OccurrenceCategory")
    schemaRef.Title = "OccurrenceCategory"
    for _, v := range OccurrenceCategoryValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["OccurrenceCategory"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/OccurrenceCategory"}
}

func (m *OccurrenceCategory) Fake(f *gofakeit.Faker) (any, error) {
	return string(OccurrenceCategoryValues[f.IntN(len(OccurrenceCategoryValues) - 1)]), nil
}

// EdgeDB Marshalling
func (m OccurrenceCategory) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *OccurrenceCategory) UnmarshalEdgeDBStr(data []byte) error {
	*m = OccurrenceCategory(string(data))
	return nil
}