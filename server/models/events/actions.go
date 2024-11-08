package events

import (
	"darco/proto/models/taxonomy"

	"github.com/edgedb/edgedb-go"
)

type Spotting struct {
	ID         edgedb.UUID        `edgedb:"id" json:"id" format:"uuid"`
	TargetTaxa []taxonomy.Taxon   `edgedb:"target_taxa" json:"target_taxa"`
	Comments   edgedb.OptionalStr `edgedb:"comments" json:"comments,omitempty"`
}
