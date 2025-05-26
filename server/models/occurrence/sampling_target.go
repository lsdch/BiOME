package occurrence

import (
	"reflect"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/danielgtaylor/huma/v2"
)

type SamplingTargetKind string

//generate:enum
const (
	Community     SamplingTargetKind = "Community"
	UnknownTarget SamplingTargetKind = "Unknown"
	Taxa          SamplingTargetKind = "Taxa"
)

type SamplingTargetKindOrNever SamplingTargetKind

const NeverSampled SamplingTargetKindOrNever = "Never"

var SamplingTargetKindOrNeverValues = []SamplingTargetKindOrNever{
	SamplingTargetKindOrNever(Community),
	SamplingTargetKindOrNever(UnknownTarget),
	SamplingTargetKindOrNever(Taxa),
	NeverSampled,
}

// Register enum in OpenAPI specification
func (u SamplingTargetKindOrNever) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["SamplingTargetKindOrNever"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "SamplingTargetKindOrNever")
		schemaRef.Title = "SamplingTargetKindOrNever"
		for _, v := range SamplingTargetKindOrNeverValues {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["SamplingTargetKindOrNever"] = schemaRef
	}

	return &huma.Schema{Ref: "#/components/schemas/SamplingTargetKindOrNever"}
}

func (m *SamplingTargetKindOrNever) Fake(f *gofakeit.Faker) (any, error) {
	return string(SamplingTargetKindOrNeverValues[f.IntN(len(SamplingTargetKindOrNeverValues)-1)]), nil
}

// Gel Marshalling
func (m SamplingTargetKindOrNever) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *SamplingTargetKindOrNever) UnmarshalEdgeDBStr(data []byte) error {
	*m = SamplingTargetKindOrNever(string(data))
	return nil
}
