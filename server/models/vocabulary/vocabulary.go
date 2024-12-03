package vocabulary

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/people"
	"encoding/json"
	"fmt"
	"maps"

	"github.com/edgedb/edgedb-go"
)

type Vocabulary struct {
	ID          edgedb.UUID        `edgedb:"id" json:"id" format:"uuid"`
	Label       string             `edgedb:"label" json:"label"`
	Code        string             `edgedb:"code" json:"code"`
	Description edgedb.OptionalStr `edgedb:"description" json:"description,omitempty"`
}

type VocabularyInput struct {
	Label       string                       `json:"label"`
	Code        string                       `json:"code"`
	Description models.OptionalInput[string] `json:"description,omitempty"`
}

type VocabularyUpdate struct {
	Label       models.OptionalInput[string] `edgedb:"label" json:"label,omitempty"`
	Code        models.OptionalInput[string] `edgedb:"code" json:"code,omitempty"`
	Description models.OptionalNull[string]  `edgedb:"description" json:"description,omitempty"`
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

type Fixative struct {
	Vocabulary `edgedb:"$inline" json:",inline"`
	Meta       people.Meta `edgedb:"meta" json:"meta"`
}

func ListFixatives(db edgedb.Executor) ([]Fixative, error) {
	var items = []Fixative{}
	err := db.Query(context.Background(),
		`select samples::Fixative { ** } order by .label`,
		&items)
	return items, err
}

type FixativeInput struct {
	VocabularyInput `json:",inline"`
}

func (i FixativeInput) Save(e edgedb.Executor) (created Fixative, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`select (insert samples::Fixative { ** })`,
		&created, data)
	return
}
