package vocabulary

import (
	"fmt"
	"maps"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
)

type Vocabulary struct {
	ID          geltypes.UUID        `gel:"id" json:"id" format:"uuid"`
	Label       string               `gel:"label" json:"label"`
	Code        string               `gel:"code" json:"code"`
	Description geltypes.OptionalStr `gel:"description" json:"description,omitempty"`
}

type VocabularyInput struct {
	Label       string                       `json:"label"`
	Code        string                       `json:"code"`
	Description models.OptionalInput[string] `json:"description,omitempty"`
}

type VocabularyUpdate struct {
	Label       models.OptionalInput[string] `gel:"label" json:"label,omitempty"`
	Code        models.OptionalInput[string] `gel:"code" json:"code,omitempty"`
	Description models.OptionalNull[string]  `gel:"description" json:"description,omitempty"`
}

// FieldMappingsWith defines Vocabulary field mappings to be used with db.UpdateQuery.
// Variadic parameters allow adding extra mappings,
// e.g. when VocularyUpdate is embedded in another struct
func (v VocabularyUpdate) FieldMappingsWith(jsonItem string, extend ...map[string]string) map[string]string {
	m := map[string]string{
		"label":       fmt.Sprintf("<str>%s['label']", jsonItem),
		"code":        fmt.Sprintf("<str>%s['code']", jsonItem),
		"description": fmt.Sprintf("<str>%s['description']", jsonItem),
	}
	for _, e := range extend {
		maps.Copy(m, e)
	}
	return m
}
