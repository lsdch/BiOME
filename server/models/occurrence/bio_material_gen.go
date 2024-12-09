// This file is auto-generated *DO NOT EDIT*

package occurrence


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var BioMaterialCategoryValues = []BioMaterialCategory{
	Internal,
	External,
}

// Register enum in OpenAPI specification
func (u BioMaterialCategory) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["BioMaterialCategory"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "BioMaterialCategory")
    schemaRef.Title = "BioMaterialCategory"
    for _, v := range BioMaterialCategoryValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["BioMaterialCategory"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/BioMaterialCategory"}
}

func (m *BioMaterialCategory) Fake(f *gofakeit.Faker) (any, error) {
	return string(BioMaterialCategoryValues[f.IntN(len(BioMaterialCategoryValues) - 1)]), nil
}

// EdgeDB Marshalling
func (m BioMaterialCategory) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *BioMaterialCategory) UnmarshalEdgeDBStr(data []byte) error {
	*m = BioMaterialCategory(string(data))
	return nil
}