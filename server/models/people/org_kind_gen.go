// This file is auto-generated *DO NOT EDIT*

package people


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var OrgKindValues = []OrgKind{
	Lab,
	FoundingAgency,
	SequencingPlatform,
	Other,
}

// Register enum in OpenAPI specification
func (u OrgKind) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["OrgKind"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "OrgKind")
    schemaRef.Title = "OrgKind"
    for _, v := range OrgKindValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["OrgKind"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/OrgKind"}
}

func (m *OrgKind) Fake(f *gofakeit.Faker) (any, error) {
	return string(OrgKindValues[f.IntN(len(OrgKindValues) - 1)]), nil
}

// Gel Marshalling
func (m OrgKind) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *OrgKind) UnmarshalEdgeDBStr(data []byte) error {
	*m = OrgKind(string(data))
	return nil
}