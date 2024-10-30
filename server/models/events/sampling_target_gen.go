// This file is auto-generated *DO NOT EDIT*

package events


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var SamplingTargetKindValues = []SamplingTargetKind{
	Community,
	UnknownTarget,
	Taxa,
}

// Register enum in OpenAPI specification
func (u SamplingTargetKind) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["SamplingTargetKind"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "SamplingTargetKind")
    schemaRef.Title = "SamplingTargetKind"
    for _, v := range SamplingTargetKindValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["SamplingTargetKind"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/SamplingTargetKind"}
}

func (m *SamplingTargetKind) Fake(f *gofakeit.Faker) (any, error) {
	return string(SamplingTargetKindValues[f.IntN(len(SamplingTargetKindValues) - 1)]), nil
}

// EdgeDB Marshalling
func (m SamplingTargetKind) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *SamplingTargetKind) UnmarshalEdgeDBStr(data []byte) error {
	*m = SamplingTargetKind(string(data))
	return nil
}