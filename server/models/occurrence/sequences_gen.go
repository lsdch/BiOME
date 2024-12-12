// This file is auto-generated *DO NOT EDIT*

package occurrence


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var ExternalSeqCategoryValues = []ExternalSeqCategory{
	NCBI,
	PersCom,
}

// Register enum in OpenAPI specification
func (u ExternalSeqCategory) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["ExternalSeqCategory"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "ExternalSeqCategory")
    schemaRef.Title = "ExternalSeqCategory"
    for _, v := range ExternalSeqCategoryValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["ExternalSeqCategory"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/ExternalSeqCategory"}
}

func (m *ExternalSeqCategory) Fake(f *gofakeit.Faker) (any, error) {
	return string(ExternalSeqCategoryValues[f.IntN(len(ExternalSeqCategoryValues) - 1)]), nil
}

// EdgeDB Marshalling
func (m ExternalSeqCategory) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *ExternalSeqCategory) UnmarshalEdgeDBStr(data []byte) error {
	*m = ExternalSeqCategory(string(data))
	return nil
}