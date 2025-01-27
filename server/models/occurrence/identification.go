package occurrence

import (
	"fmt"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/taxonomy"

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

type IdentificationUpdate struct {
	Taxon        models.OptionalInput[string]                 `edgedb:"taxon" json:"taxon,omitempty"`
	IdentifiedBy models.OptionalNull[string]                  `edgedb:"identified_by" json:"identified_by,omitempty"`
	IdentifiedOn models.OptionalInput[DateWithPrecisionInput] `edgedb:"identified_on" json:"identified_on,omitempty"`
}

func (u IdentificationUpdate) UpdateQuery(fieldName string) string {
	return u.FieldMappings(fieldName).Query(u)
}

func (u IdentificationUpdate) FieldMappings(fieldName string) db.UpdateQuery {
	return db.UpdateQuery{
		Frame: fmt.Sprintf(`#edgeql
			with identification := <json>$0
			update %s set {
				%%s # completed using Mappings
			}
		`, fieldName),
		Mappings: map[string]string{
			"taxon": `#edgeql
				assert_exists(
					select taxonomy::Taxon filter .name = <str>data['taxon']
				)`,
			"identified_by": `#edgeql
				assert_exists(
					select people::Person filter .alias = <str>data['identified_by']
				)`,
			"identified_on": `#edgeql
				date::from_json_with_precision(<json>data['identified_on'])
			`,
		},
	}
}
