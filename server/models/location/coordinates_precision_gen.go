// This file is auto-generated *DO NOT EDIT*

package location


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var CoordinatePrecisionValues = []CoordinatePrecision{
	M10,
	M100,
	KM1,
	KM10,
	KM100,
	Unknown,
}

// Register enum in OpenAPI specification
func (u CoordinatePrecision) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["CoordinatePrecision"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "CoordinatePrecision")
    schemaRef.Title = "CoordinatePrecision"
    for _, v := range CoordinatePrecisionValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["CoordinatePrecision"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/CoordinatePrecision"}
}

func (m *CoordinatePrecision) Fake(f *gofakeit.Faker) (any, error) {
	return string(CoordinatePrecisionValues[f.IntN(len(CoordinatePrecisionValues) - 1)]), nil
}

// EdgeDB Marshalling
func (m CoordinatePrecision) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *CoordinatePrecision) UnmarshalEdgeDBStr(data []byte) error {
	*m = CoordinatePrecision(string(data))
	return nil
}