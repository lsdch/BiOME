// This file is auto-generated *DO NOT EDIT*

package occurrence


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var BioMatSortKeyValues = []BioMatSortKey{
	BioMatSortCode,
	BioMatSortSite,
	BioMatSortSamplingDate,
	BioMatSortIdentifiedOn,
	BioMatSortTaxon,
	BioMatSortIdentifiedBy,
	BioMatSortLastUpdated,
}

// Register enum in OpenAPI specification
func (u BioMatSortKey) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["BioMatSortKey"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "BioMatSortKey")
    schemaRef.Title = "BioMatSortKey"
    for _, v := range BioMatSortKeyValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["BioMatSortKey"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/BioMatSortKey"}
}

func (m *BioMatSortKey) Fake(f *gofakeit.Faker) (any, error) {
	return string(BioMatSortKeyValues[f.IntN(len(BioMatSortKeyValues) - 1)]), nil
}

// Gel Marshalling
func (m BioMatSortKey) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *BioMatSortKey) UnmarshalEdgeDBStr(data []byte) error {
	*m = BioMatSortKey(string(data))
	return nil
}