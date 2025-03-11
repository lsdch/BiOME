// This file is auto-generated *DO NOT EDIT*

package taxonomy


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var TaxonStatusValues = []TaxonStatus{
	Accepted,
	Unreferenced,
	Unclassified,
}

// Register enum in OpenAPI specification
func (u TaxonStatus) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["TaxonStatus"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "TaxonStatus")
    schemaRef.Title = "TaxonStatus"
    for _, v := range TaxonStatusValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["TaxonStatus"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/TaxonStatus"}
}

func (m *TaxonStatus) Fake(f *gofakeit.Faker) (any, error) {
	return string(TaxonStatusValues[f.IntN(len(TaxonStatusValues) - 1)]), nil
}

// Gel Marshalling
func (m TaxonStatus) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *TaxonStatus) UnmarshalEdgeDBStr(data []byte) error {
	*m = TaxonStatus(string(data))
	return nil
}