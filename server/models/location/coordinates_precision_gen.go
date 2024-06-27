// This file is auto-generated *DO NOT EDIT*

package location


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var CoordinatesPrecisionValues = []CoordinatesPrecision{
	M100,
	KM1,
	KM10,
	KM100,
	Unknown,
}

// Register enum in OpenAPI specification
func (u CoordinatesPrecision) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["CoordinatesPrecision"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "CoordinatesPrecision")
    schemaRef.Title = "CoordinatesPrecision"
    for _, v := range CoordinatesPrecisionValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["CoordinatesPrecision"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/CoordinatesPrecision"}
}

func (m *CoordinatesPrecision) Fake(f *gofakeit.Faker) (any, error) {
	return string(CoordinatesPrecisionValues[f.IntN(len(CoordinatesPrecisionValues) - 1)]), nil
}

// EdgeDB Marshalling
func (m CoordinatesPrecision) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *CoordinatesPrecision) UnmarshalEdgeDBStr(data []byte) error {
	*m = CoordinatesPrecision(string(data))
	return nil
}