// This file is auto-generated *DO NOT EDIT*

package occurrence


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var BioMaterialTypeValues = []BioMaterialType{
	Internal,
	External,
}

// Register enum in OpenAPI specification
func (u BioMaterialType) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["BioMaterialType"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "BioMaterialType")
    schemaRef.Title = "BioMaterialType"
    for _, v := range BioMaterialTypeValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["BioMaterialType"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/BioMaterialType"}
}

func (m *BioMaterialType) Fake(f *gofakeit.Faker) (any, error) {
	return string(BioMaterialTypeValues[f.IntN(len(BioMaterialTypeValues) - 1)]), nil
}

// EdgeDB Marshalling
func (m BioMaterialType) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *BioMaterialType) UnmarshalEdgeDBStr(data []byte) error {
	*m = BioMaterialType(string(data))
	return nil
}