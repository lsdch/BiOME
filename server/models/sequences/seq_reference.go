package sequences

import (
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models/references"
)

type SeqReference struct {
	ID        geltypes.UUID         `gel:"id" json:"id" format:"uuid"`
	DB        references.DataSource `gel:"db" json:"db"`
	Accession string                `gel:"accession" json:"accession"`
	IsOrigin  bool                  `gel:"is_origin" json:"is_origin"`
	Code      string                `gel:"code" json:"-"`
}

type SeqReferenceInput struct {
	DB        string `json:"db"`
	Accession string `json:"accession"`
	IsOrigin  bool   `json:"is_origin"`
}
