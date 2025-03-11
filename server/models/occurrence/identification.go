package occurrence

import (
	"fmt"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/taxonomy"
)

type Identification struct {
	ID           geltypes.UUID         `gel:"id" json:"id" format:"uuid"`
	Taxon        taxonomy.Taxon        `gel:"taxon" json:"taxon"`
	IdentifiedBy people.OptionalPerson `gel:"identified_by" json:"identified_by"`
	IdentifiedOn DateWithPrecision     `gel:"identified_on" json:"identified_on"`
	Meta         people.Meta           `gel:"meta" json:"meta"`
}

type IdentificationInput struct {
	Taxon        string                       `gel:"taxon" json:"taxon"`
	IdentifiedBy models.OptionalInput[string] `gel:"identified_by" json:"identified_by"`
	IdentifiedOn DateWithPrecisionInput       `gel:"identified_on" json:"identified_on"`
}

type IdentificationUpdate struct {
	Taxon        models.OptionalInput[string]                 `gel:"taxon" json:"taxon,omitempty"`
	IdentifiedBy models.OptionalNull[string]                  `gel:"identified_by" json:"identified_by,omitempty"`
	IdentifiedOn models.OptionalInput[DateWithPrecisionInput] `gel:"identified_on" json:"identified_on,omitempty"`
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
