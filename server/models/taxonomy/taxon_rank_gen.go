// This file is auto-generated *DO NOT EDIT*

package taxonomy


import (
	"reflect"

  "fmt"
	"github.com/danielgtaylor/huma/v2"
  "github.com/go-faker/faker/v4"
	"math/rand"
)



var TaxonRankValues = []TaxonRank{
	Kingdom,
	Phylum,
	Class,
	Family,
	Genus,
	Species,
	Subspecies,
}

// Register enum in OpenAPI specification
func (u TaxonRank) Schema(r huma.Registry) *huma.Schema {
	schemaRef := r.Schema(reflect.TypeOf(""), true, "TaxonRank")
  schemaRef.Title = "TaxonRank"
  for _, v := range TaxonRankValues {
	  schemaRef.Enum = append(schemaRef.Enum, string(v))
  }
  r.Map()["TaxonRank"] = schemaRef


  schema := r.Schema(reflect.TypeOf(""), true, "TaxonRank")
  schema.Ref = "#/components/schemas/TaxonRank"
	return schema
}

func init () {
  // Faker
  faker.AddProvider("TaxonRank",
    func(v reflect.Value) (interface{}, error) {
      idx := rand.Intn(len(TaxonRankValues))
      fmt.Printf("Called provided for TaxonRank: %s\n", TaxonRank(TaxonRankValues[idx]))
      return string(TaxonRankValues[idx]), nil
    })
}

// EdgeDB Marshalling
func (m TaxonRank) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *TaxonRank) UnmarshalEdgeDBStr(data []byte) error {
	*m = TaxonRank(string(data))
	return nil
}