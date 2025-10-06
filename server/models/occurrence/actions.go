package occurrence

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models/taxonomy"
)

type Flagging struct {
	ID                geltypes.UUID        `gel:"id" json:"-" format:"uuid"`
	TargetTaxa        []taxonomy.Taxon     `gel:"target_taxa" json:"target_taxa,omitempty"`
	AbioticParameters []AbioticParameter   `gel:"abiotic_parameters" json:"abiotic_parameters,omitempty"`
	Indications       geltypes.OptionalStr `gel:"indications" json:"indications,omitempty"`
}

type FlaggingInput struct {
	TargetTaxa        []string             `json:"target_taxa,omitempty" doc:"List of taxon names"`
	AbioticParameters []string             `json:"abiotic_parameters,omitempty" doc:"List of abiotic parameter codes"`
	Indications       geltypes.OptionalStr `json:"indications,omitempty"`
}

func (u FlaggingInput) Save(e geltypes.Executor, siteCode string) (flaggings []Flagging, err error) {
	data, _ := json.Marshal(u)
	err = e.Query(context.Background(),
		`#edgeql
			with data := <json>$0
      select (
				insert events::Flagging {
					site := (select location::Site filter .code = <str>$1),
					target_taxa := (
						select taxonomy::Taxon filter .name in json_array_unpack(<array<str>>json_get(<json>$0, 'target_taxa'))
					),
					abiotic_parameters := (
						select events::AbioticParameter filter .code in json_array_unpack(<array<str>>json_get(<json>$0, 'abiotic_parameters'))
					),
					indications := <str>json_get(<json>$0, 'indications')
				}
			) { *, target_taxa: { * }, abiotic_parameters: { * } }
		`,
		&flaggings, siteCode, data)
	return
}
