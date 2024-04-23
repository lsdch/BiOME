// This file is auto-generated *DO NOT EDIT*

package people


import (
	"reflect"

  "fmt"
	"github.com/danielgtaylor/huma/v2"
  "github.com/go-faker/faker/v4"
	"math/rand"
)



var InstitutionKindValues = []InstitutionKind{
	Lab,
	FoundingAgency,
	SequencingPlatform,
	Other,
}

// Register enum in OpenAPI specification
func (u InstitutionKind) Schema(r huma.Registry) *huma.Schema {
	schemaRef := r.Schema(reflect.TypeOf(""), true, "InstitutionKind")
  schemaRef.Title = "InstitutionKind"
  for _, v := range InstitutionKindValues {
	  schemaRef.Enum = append(schemaRef.Enum, string(v))
  }
  r.Map()["InstitutionKind"] = schemaRef


  schema := r.Schema(reflect.TypeOf(""), true, "InstitutionKind")
  schema.Ref = "#/components/schemas/InstitutionKind"
	return schema
}

func init () {
  // Faker
  faker.AddProvider("InstitutionKind",
    func(v reflect.Value) (interface{}, error) {
      idx := rand.Intn(len(InstitutionKindValues))
      fmt.Printf("Called provided for InstitutionKind: %s\n", InstitutionKind(InstitutionKindValues[idx]))
      return string(InstitutionKindValues[idx]), nil
    })
}

// EdgeDB Marshalling
func (m InstitutionKind) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *InstitutionKind) UnmarshalEdgeDBStr(data []byte) error {
	*m = InstitutionKind(string(data))
	return nil
}