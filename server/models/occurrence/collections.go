package occurrence

import (
	"context"

	"github.com/geldata/gel-go/geltypes"
)

func ListCollections(db geltypes.Executor) ([]string, error) {
	var items []string
	err := db.Query(context.Background(),
		`#edgeql
			with collections := (select occurrence::Occurrence.collections order by .name)
			select distinct collections.name
		`, &items)
	return items, err
}
