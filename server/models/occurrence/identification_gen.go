// This file is auto-generated *DO NOT EDIT*

package occurrence


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var IdentificationQualifierValues = []IdentificationQualifier{
	IdQualifCF,
	IdQualifAFF,
}

// Register enum in OpenAPI specification
func (u IdentificationQualifier) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["IdentificationQualifier"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "IdentificationQualifier")
    schemaRef.Title = "IdentificationQualifier"
    for _, v := range IdentificationQualifierValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["IdentificationQualifier"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/IdentificationQualifier"}
}

func (m *IdentificationQualifier) Fake(f *gofakeit.Faker) (any, error) {
	return string(IdentificationQualifierValues[f.IntN(len(IdentificationQualifierValues) - 1)]), nil
}

// Gel Marshalling
func (m IdentificationQualifier) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *IdentificationQualifier) UnmarshalEdgeDBStr(data []byte) error {
	*m = IdentificationQualifier(string(data))
	return nil
}