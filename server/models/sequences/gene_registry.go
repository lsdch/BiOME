package sequences

import (
	"context"
	"encoding/json"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/vocabulary"

	"github.com/edgedb/edgedb-go"
)

type Gene struct {
	vocabulary.Vocabulary `edgedb:"$inline" json:",inline"`
	IsMOTUDelimiter       bool        `edgedb:"motu" json:"is_MOTU_delimiter"`
	Meta                  people.Meta `edgedb:"meta" json:"meta"`
}

func ListGenes(db edgedb.Executor) ([]Gene, error) {
	var items = []Gene{}
	err := db.Query(context.Background(),
		`select seq::Gene { ** };`,
		&items)
	return items, err
}

type GeneInput struct {
	vocabulary.VocabularyInput `json:",inline"`
	IsMOTUDelimiter            bool `json:"is_MOTU_delimiter,omitempty" default:"false"`
}

func (i GeneInput) Save(e edgedb.Executor) (created Gene, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0,
			select (insert seq::Gene {
				label := <str>data['label'],
				code := <str>data['code'],
				description := <str>json_get(data, 'description'),
				motu := <bool>json_get(data, "is_MOTU_delimiter") ?? false
			}) { ** }
		`, &created, data)
	return
}

type GeneUpdate struct {
	vocabulary.VocabularyUpdate `edgedb:"$inline" json:",inline"`
	IsMOTUDelimiter             models.OptionalInput[bool] `edgedb:"motu" json:"is_MOTU_delimiter,omitempty"`
}

func (u GeneUpdate) Save(e edgedb.Executor, code string) (updated Gene, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$1,
			select (update seq::Gene filter .code = <str>$0 set {
				%s
			}) { ** }
		`,
		Mappings: u.FieldMappingsWith("item", map[string]string{
			"motu": "<bool>json_get(item, 'is_MOTU_delimiter')",
		}),
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, code, data)
	return
}

func DeleteGene(db edgedb.Executor, code string) (deleted Gene, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
				delete seq::Gene filter .code = <str>$0
		 	) { ** }
		`,
		&deleted, code)
	return
}
