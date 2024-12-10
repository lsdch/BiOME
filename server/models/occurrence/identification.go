package occurrence

import (
	"darco/proto/models"
	"darco/proto/models/people"
	"darco/proto/models/taxonomy"

	"github.com/edgedb/edgedb-go"
)

type Identification struct {
	ID           edgedb.UUID           `edgedb:"id" json:"id" format:"uuid"`
	Taxon        taxonomy.Taxon        `edgedb:"taxon" json:"taxon"`
	IdentifiedBy people.OptionalPerson `edgedb:"identified_by" json:"identified_by"`
	IdentifiedOn DateWithPrecision     `edgedb:"identified_on" json:"identified_on"`
	Meta         people.Meta           `edgedb:"meta" json:"meta"`
}

type IdentificationInput struct {
	Taxon        string                       `edgedb:"taxon" json:"taxon"`
	IdentifiedBy models.OptionalInput[string] `edgedb:"identified_by" json:"identified_by"`
	IdentifiedOn DateWithPrecisionInput       `edgedb:"identified_on" json:"identified_on"`
}
