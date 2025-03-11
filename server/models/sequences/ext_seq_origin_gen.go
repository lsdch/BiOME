// This file is auto-generated *DO NOT EDIT*

package sequences


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var ExtSeqOriginValues = []ExtSeqOrigin{
	Lab,
	DB,
	PersCom,
}

// Register enum in OpenAPI specification
func (u ExtSeqOrigin) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["ExtSeqOrigin"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "ExtSeqOrigin")
    schemaRef.Title = "ExtSeqOrigin"
    for _, v := range ExtSeqOriginValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["ExtSeqOrigin"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/ExtSeqOrigin"}
}

func (m *ExtSeqOrigin) Fake(f *gofakeit.Faker) (any, error) {
	return string(ExtSeqOriginValues[f.IntN(len(ExtSeqOriginValues) - 1)]), nil
}

// Gel Marshalling
func (m ExtSeqOrigin) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *ExtSeqOrigin) UnmarshalEdgeDBStr(data []byte) error {
	*m = ExtSeqOrigin(string(data))
	return nil
}