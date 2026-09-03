// This file is auto-generated *DO NOT EDIT*

package imports


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var RunnerStatusValues = []RunnerStatus{
	Created,
	Staging,
	Staged,
	Running,
	NeedsResolution,
	ReadyToMaterialize,
	Materializing,
	Completed,
	Failed,
	Cancelled,
}

// Register enum in OpenAPI specification
func (u RunnerStatus) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["RunnerStatus"] == nil {
    schemaRef := r.Schema(reflect.TypeFor[string](), true, "RunnerStatus")
    schemaRef.Title = "RunnerStatus"
    for _, v := range RunnerStatusValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["RunnerStatus"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/RunnerStatus"}
}

func (m *RunnerStatus) Fake(f *gofakeit.Faker) (any, error) {
	return string(RunnerStatusValues[f.IntN(len(RunnerStatusValues) - 1)]), nil
}