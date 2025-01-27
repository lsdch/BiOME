package occurrence

import (
	"context"
	"encoding/json"

	"github.com/lsdch/biome/models/taxonomy"

	"github.com/edgedb/edgedb-go"
)

type Spotting struct {
	edgedb.Optional
	ID         edgedb.UUID        `edgedb:"id" json:"-" format:"uuid"`
	TargetTaxa []taxonomy.Taxon   `edgedb:"target_taxa" json:"target_taxa,omitempty"`
	Comments   edgedb.OptionalStr `edgedb:"comments" json:"comments,omitempty"`
}

type SpottingUpdate []string

func (u SpottingUpdate) Save(e edgedb.Executor, eventID edgedb.UUID) (spottings []taxonomy.Taxon, err error) {
	if len(u) == 0 {
		return
	}
	data, _ := json.Marshal(u)
	err = e.Query(context.Background(),
		`#edgeql
			with ev := (
				update events::Event
        filter .id = <uuid>$0
        set {
					spottings := assert_distinct((
						select taxonomy::Taxon
						filter .name in <str>json_array_unpack(<json>$1)
					))
				}
			)
      select (ev.spottings) { * }
		`,
		&spottings, eventID, data)
	return
}
