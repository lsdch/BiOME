// This file is auto-generated *DO NOT EDIT*

package taxonomy


import (
	"reflect"

  "fmt"
	"github.com/danielgtaylor/huma/v2"
  "github.com/go-faker/faker/v4"
	"math/rand"
)



var TaxonStatusValues = []TaxonStatus{
	Accepted,
	Synonym,
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

func init () {
  // Faker
  faker.AddProvider("TaxonStatus",
    func(v reflect.Value) (interface{}, error) {
      idx := rand.Intn(len(TaxonStatusValues))
      fmt.Printf("Called provided for TaxonStatus: %s\n", TaxonStatus(TaxonStatusValues[idx]))
      return string(TaxonStatusValues[idx]), nil
    })
}

// EdgeDB Marshalling
func (m TaxonStatus) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *TaxonStatus) UnmarshalEdgeDBStr(data []byte) error {
	*m = TaxonStatus(string(data))
	return nil
}