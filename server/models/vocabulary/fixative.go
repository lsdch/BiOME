package vocabulary

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/people"
)

type Fixative struct {
	Vocabulary `gel:"$inline" json:",inline"`
	Meta       people.Meta `gel:"meta" json:"meta"`
}

func ListFixatives(db geltypes.Executor) ([]Fixative, error) {
	var items = []Fixative{}
	err := db.Query(context.Background(),
		`select samples::Fixative { ** } order by .label`,
		&items)
	return items, err
}

type FixativeInput VocabularyInput

func (i FixativeInput) Save(e geltypes.Executor) (created Fixative, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0
			select (insert samples::Fixative {
				label := <str>data['label'],
				code := <str>data['code'],
				description := <str>json_get(data, 'description')
			}) { ** }
		`,
		&created, data)
	return
}

type FixativeUpdate VocabularyUpdate

func (u FixativeUpdate) Save(e geltypes.Executor, code string) (updated Fixative, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$1,
			select (update samples::Fixative filter .code = <str>$0 set {
				%s
			}) { ** }
		`,
		Mappings: VocabularyUpdate(u).FieldMappingsWith("item"),
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, code, data)
	return
}

func DeleteFixative(db geltypes.Executor, code string) (deleted Fixative, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select ( delete samples::Fixative filter .code = <str>$0 ) { ** };
		`,
		&deleted, code)
	return
}
