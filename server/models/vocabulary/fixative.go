package vocabulary

import (
	"context"
	"darco/proto/db"
	"darco/proto/models/people"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
)

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

type FixativeInput VocabularyInput

func (i FixativeInput) Save(e edgedb.Executor) (created Fixative, err error) {
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

func (u FixativeUpdate) Save(e edgedb.Executor, code string) (updated Fixative, err error) {
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

func DeleteFixative(db edgedb.Executor, code string) (deleted Fixative, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select ( delete samples::Fixative filter .code = <str>$0 ) { ** };
		`,
		&deleted, code)
	return
}
