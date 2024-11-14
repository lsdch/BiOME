package occurrence

import (
	"darco/proto/models/people"
	"darco/proto/models/taxonomy"
)

type Identification struct {
	Taxon        taxonomy.Taxon        `edgedb:"taxon" json:"taxon"`
	IdentifiedBy people.OptionalPerson `edgedb:"identified_by" json:"identified_by"`
	Date         DateWithPrecision     `edgedb:"date" json:"date"`
	Meta         people.Meta           `edgedb:"meta" json:"meta"`
}
