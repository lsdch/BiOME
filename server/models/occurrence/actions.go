package occurrence

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models/taxonomy"
)

type Spotting struct {
	geltypes.Optional
	ID         geltypes.UUID        `gel:"id" json:"-" format:"uuid"`
	TargetTaxa []taxonomy.Taxon     `gel:"target_taxa" json:"target_taxa,omitempty"`
	Comments   geltypes.OptionalStr `gel:"comments" json:"comments,omitempty"`
}

type SpottingUpdate []string

func (u SpottingUpdate) Save(e geltypes.Executor, eventID geltypes.UUID) (spottings []taxonomy.Taxon, err error) {
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
