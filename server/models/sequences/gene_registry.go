package sequences

import (
	"context"
	"darco/proto/models/people"
	"darco/proto/models/vocabulary"

	"github.com/edgedb/edgedb-go"
)

type Gene struct {
	vocabulary.Vocabulary `edgedb:"$inline" json:",inline"`
	IsDelimiterMOTU       bool        `edgedb:"motu" json:"is_delimiter_MOTU"`
	Meta                  people.Meta `edgedb:"meta" json:"meta"`
}

func ListGenes(db edgedb.Executor) ([]Gene, error) {
	var items = []Gene{}
	err := db.Query(context.Background(),
		`select seq::Gene { ** };`,
		&items)
	return items, err
}
