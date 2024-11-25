package sequences

import (
	"context"
	"darco/proto/models/people"
	"darco/proto/models/vocabulary"
	"encoding/json"

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

func (i GeneInput) Create(e edgedb.Executor) (created Gene, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data = <json>$0,
			select (insert seq::Gene {
				label := <str>data['label'],
				code := <str>data['code'],
				description := <str>json_get(data, 'description')
				motu := <bool>json_get(data, "is_MOTU_delimiter")
			})
		`, &created, data)
	return
}
