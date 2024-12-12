package sequences

import (
	"context"
	"darco/proto/db"
	"darco/proto/models"
	"darco/proto/models/people"
	"darco/proto/models/vocabulary"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
)

type SeqDB struct {
	vocabulary.Vocabulary `edgedb:"$inline" json:",inline"`
	LinkTemplate          edgedb.OptionalStr `edgedb:"link_template" json:"link_template"`
	Meta                  people.Meta        `edgedb:"meta" json:"meta"`
}

func ListSeqDB(db edgedb.Executor) ([]SeqDB, error) {
	var items = []SeqDB{}
	err := db.Query(context.Background(),
		`#edgeql
			select seq::SeqDB { ** };
		`,
		&items)
	return items, err
}

func DeleteSeqDB(db edgedb.Executor, code string) (deleted SeqDB, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
			  delete seq::SeqDB filter .code = <str>$0
		 	) { ** };
		`,
		&deleted, code)
	return
}

type SeqDBInput struct {
	vocabulary.VocabularyInput `edgedb:"$inline" json:",inline"`
	LinkTemplate               models.OptionalInput[string] `edgedb:"link_template" json:"link_template,omitempty"`
}

func (i SeqDBInput) Save(e edgedb.Executor) (created SeqDB, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0,
			select (insert seq::SeqDB {
				label := <str>data['label'],
				code := <str>data['code'],
				description := <str>json_get(data, 'description'),
				link_template := <str>json_get(data, 'link_template')
			}) { ** }
		`, &created, data)
	return
}

type SeqDBUpdate struct {
	vocabulary.VocabularyUpdate `edgedb:"$inline" json:",inline"`
	LinkTemplate                models.OptionalNull[string] `edgedb:"link_template" json:"link_template,omitempty"`
}

func (u SeqDBUpdate) Save(e edgedb.Executor, code string) (updated SeqDB, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$1,
			select (update seq::SeqDB filter .code = <str>$0 set {
				%s
			}) { ** }
		`,
		Mappings: u.FieldMappingsWith("item", map[string]string{
			"link_template": "<str>item['link_template']",
		}),
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, code, data)
	return
}
