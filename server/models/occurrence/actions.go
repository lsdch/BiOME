package occurrence

import (
	"context"
	"darco/proto/db"
	"darco/proto/models"
	"darco/proto/models/taxonomy"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
)

type Spotting struct {
	edgedb.Optional
	ID         edgedb.UUID        `edgedb:"id" json:"-" format:"uuid"`
	TargetTaxa []taxonomy.Taxon   `edgedb:"target_taxa" json:"target_taxa,omitempty"`
	Comments   edgedb.OptionalStr `edgedb:"comments" json:"comments,omitempty"`
}

type SpottingUpdate struct {
	TargetTaxa models.OptionalNull[[]string] `edgedb:"target_taxa" json:"target_taxa,omitempty"`
	Comments   models.OptionalNull[string]   `edgedb:"comments" json:"comments,omitempty"`
}

func (u SpottingUpdate) Save(e edgedb.Executor, eventID edgedb.UUID) (updated Spotting, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with data := <json>$1,
			select (
				insert events::Spotting {
					event := <events::Event><uuid>$0,
					target_taxa := (
						select taxonomy::Taxon
						filter .name in <str>json_array_unpack(json_get(data,'target_taxa'))
					),
					comments := <str>json_get(data, 'comments')
				} unless conflict on .event else (
					update events::Spotting set {
						%s
					}
				)
			) { *,  target_taxa: { * } }
		`,
		Mappings: map[string]string{
			"target_taxa": `#edgeql
				(
					select taxonomy::Taxon
					filter .name in <str>json_array_unpack(json_get(data,'target_taxa'))
				)
			`,
			"comments": "<str>json_get(data, 'comments')",
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, eventID, data)
	return
}
