// This file is auto-generated *DO NOT EDIT*

package dataset


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var DatasetCategoryValues = []DatasetCategory{
	Site,
	Occurrence,
	Seq,
}

// Register enum in OpenAPI specification
func (u DatasetCategory) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["DatasetCategory"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "DatasetCategory")
    schemaRef.Title = "DatasetCategory"
    for _, v := range DatasetCategoryValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["DatasetCategory"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/DatasetCategory"}
}

func (m *DatasetCategory) Fake(f *gofakeit.Faker) (any, error) {
	return string(DatasetCategoryValues[f.IntN(len(DatasetCategoryValues) - 1)]), nil
}

// Gel Marshalling
func (m DatasetCategory) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *DatasetCategory) UnmarshalEdgeDBStr(data []byte) error {
	*m = DatasetCategory(string(data))
	return nil
}