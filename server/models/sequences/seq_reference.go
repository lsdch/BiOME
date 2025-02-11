package sequences

import (
	"github.com/edgedb/edgedb-go"
	"github.com/lsdch/biome/models/references"
)

type SeqReference struct {
	ID        edgedb.UUID           `edgedb:"id" json:"id" format:"uuid"`
	DB        references.DataSource `edgedb:"db" json:"db"`
	Accession string                `edgedb:"accession" json:"accession"`
	IsOrigin  bool                  `edgedb:"is_origin" json:"is_origin"`
	Code      string                `edgedb:"code" json:"-"`
}

type SeqReferenceInput struct {
	DB        string `json:"db"`
	Accession string `json:"accession"`
	IsOrigin  bool   `json:"is_origin"`
}
