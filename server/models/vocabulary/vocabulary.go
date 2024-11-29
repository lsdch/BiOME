package vocabulary

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/people"
	"encoding/json"

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
