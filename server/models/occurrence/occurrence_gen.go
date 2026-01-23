// This file is auto-generated *DO NOT EDIT*

package occurrence


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var TypeStatusValues = []TypeStatus{
	Holotype,
	Neotype,
	Topotype,
}

// Register enum in OpenAPI specification
func (u TypeStatus) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["TypeStatus"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "TypeStatus")
    schemaRef.Title = "TypeStatus"
    for _, v := range TypeStatusValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["TypeStatus"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/TypeStatus"}
}

func (m *TypeStatus) Fake(f *gofakeit.Faker) (any, error) {
	return string(TypeStatusValues[f.IntN(len(TypeStatusValues) - 1)]), nil
}

// Gel Marshalling
func (m TypeStatus) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}



var BioMatSortKeyValues = []BioMatSortKey{
	BioMatSortCode,
	BioMatSortSite,
	BioMatSortSamplingDate,
	BioMatSortIdentifiedOn,
	BioMatSortTaxon,
	BioMatSortIdentification,
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



var SiteSamplingStatusValues = []SiteSamplingStatus{
	IncludeAllSites,
	IncludeSampled,
	IncludeWithOccurrences,
}

// Register enum in OpenAPI specification
func (u SiteSamplingStatus) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["SiteSamplingStatus"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "SiteSamplingStatus")
    schemaRef.Title = "SiteSamplingStatus"
    for _, v := range SiteSamplingStatusValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["SiteSamplingStatus"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/SiteSamplingStatus"}
}

func (m *SiteSamplingStatus) Fake(f *gofakeit.Faker) (any, error) {
	return string(SiteSamplingStatusValues[f.IntN(len(SiteSamplingStatusValues) - 1)]), nil
}

// Gel Marshalling
func (m SiteSamplingStatus) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *SiteSamplingStatus) UnmarshalEdgeDBStr(data []byte) error {
	*m = SiteSamplingStatus(string(data))
	return nil
}