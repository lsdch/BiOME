package occurrence

import (
	"darco/proto/models/people"
	"darco/proto/models/taxonomy"
)

type Identification struct {
	Taxon        taxonomy.Taxon        `edgedb:"taxon" json:"taxon"`
	IdentifiedBy people.OptionalPerson `edgedb:"identified_by" json:"identified_by"`
	IdentifiedOn DateWithPrecision     `edgedb:"identified_on" json:"identified_on"`
	IsType       bool                  `edgedb:"is_type" json:"is_type"`
	Meta         people.Meta           `edgedb:"meta" json:"meta"`
}
