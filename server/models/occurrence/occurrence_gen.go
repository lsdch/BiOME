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

// Gel Marshalling
func (m OccurrenceCategory) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *OccurrenceCategory) UnmarshalEdgeDBStr(data []byte) error {
	*m = OccurrenceCategory(string(data))
	return nil
}



var OccurrenceElementValues = []OccurrenceElement{
	BioMaterialElement,
	SequenceElement,
}

// Register enum in OpenAPI specification
func (u OccurrenceElement) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["OccurrenceElement"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "OccurrenceElement")
    schemaRef.Title = "OccurrenceElement"
    for _, v := range OccurrenceElementValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["OccurrenceElement"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/OccurrenceElement"}
}

func (m *OccurrenceElement) Fake(f *gofakeit.Faker) (any, error) {
	return string(OccurrenceElementValues[f.IntN(len(OccurrenceElementValues) - 1)]), nil
}

// Gel Marshalling
func (m OccurrenceElement) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *OccurrenceElement) UnmarshalEdgeDBStr(data []byte) error {
	*m = OccurrenceElement(string(data))
	return nil
}