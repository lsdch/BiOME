// This file is auto-generated *DO NOT EDIT*

package specimen


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var QuantityValues = []Quantity{
	Unknown,
	One,
	Several,
	Dozen,
	Tens,
	Hundred,
}

// Register enum in OpenAPI specification
func (u Quantity) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["Quantity"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "Quantity")
    schemaRef.Title = "Quantity"
    for _, v := range QuantityValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["Quantity"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/Quantity"}
}

func (m *Quantity) Fake(f *gofakeit.Faker) (any, error) {
	return string(QuantityValues[f.IntN(len(QuantityValues) - 1)]), nil
}

// EdgeDB Marshalling
func (m Quantity) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *Quantity) UnmarshalEdgeDBStr(data []byte) error {
	*m = Quantity(string(data))
	return nil
}