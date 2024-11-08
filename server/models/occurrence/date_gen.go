// This file is auto-generated *DO NOT EDIT*

package occurrence


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var DatePrecisionValues = []DatePrecision{
	Year,
	Month,
	Day,
	Unknown,
}

// Register enum in OpenAPI specification
func (u DatePrecision) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["DatePrecision"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "DatePrecision")
    schemaRef.Title = "DatePrecision"
    for _, v := range DatePrecisionValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["DatePrecision"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/DatePrecision"}
}

func (m *DatePrecision) Fake(f *gofakeit.Faker) (any, error) {
	return string(DatePrecisionValues[f.IntN(len(DatePrecisionValues) - 1)]), nil
}

// EdgeDB Marshalling
func (m DatePrecision) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *DatePrecision) UnmarshalEdgeDBStr(data []byte) error {
	*m = DatePrecision(string(data))
	return nil
}