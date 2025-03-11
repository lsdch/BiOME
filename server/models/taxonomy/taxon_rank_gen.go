// This file is auto-generated *DO NOT EDIT*

package taxonomy


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var TaxonRankValues = []TaxonRank{
	Kingdom,
	Phylum,
	Class,
	Order,
	Family,
	Genus,
	Species,
	Subspecies,
}

// Register enum in OpenAPI specification
func (u TaxonRank) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["TaxonRank"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "TaxonRank")
    schemaRef.Title = "TaxonRank"
    for _, v := range TaxonRankValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["TaxonRank"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/TaxonRank"}
}

func (m *TaxonRank) Fake(f *gofakeit.Faker) (any, error) {
	return string(TaxonRankValues[f.IntN(len(TaxonRankValues) - 1)]), nil
}

// Gel Marshalling
func (m TaxonRank) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *TaxonRank) UnmarshalEdgeDBStr(data []byte) error {
	*m = TaxonRank(string(data))
	return nil
}