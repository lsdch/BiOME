// This file is auto-generated *DO NOT EDIT*

package people


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var InstitutionKindValues = []InstitutionKind{
	Lab,
	FoundingAgency,
	SequencingPlatform,
	Other,
}

// Register enum in OpenAPI specification
func (u InstitutionKind) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["InstitutionKind"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "InstitutionKind")
    schemaRef.Title = "InstitutionKind"
    for _, v := range InstitutionKindValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["InstitutionKind"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/InstitutionKind"}
}

func (m *InstitutionKind) Fake(f *gofakeit.Faker) (any, error) {
	return string(InstitutionKindValues[f.IntN(len(InstitutionKindValues) - 1)]), nil
}

// EdgeDB Marshalling
func (m InstitutionKind) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *InstitutionKind) UnmarshalEdgeDBStr(data []byte) error {
	*m = InstitutionKind(string(data))
	return nil
}