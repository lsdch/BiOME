package sequences

import (
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models/references"
)

// SeqReference represents a reference to a sequence in a specific database.
type SeqReference struct {
	ID        geltypes.UUID         `gel:"id" json:"id" format:"uuid"`
	DB        references.DataSource `gel:"db" json:"db"`
	Accession string                `gel:"accession" json:"accession"`
	IsOrigin  bool                  `gel:"is_origin" json:"is_origin"`
	// For internal use in Gel
	Code string `gel:"code" json:"-"`
}

type SeqReferenceInput struct {
	DB        string `json:"db" doc:"Data source code identifier"`
	Accession string `json:"accession" doc:"Accession number or sequence identifier in the data source"`
	IsOrigin  bool   `json:"is_origin" doc:"Is this the origin of the sequence?"`
}
