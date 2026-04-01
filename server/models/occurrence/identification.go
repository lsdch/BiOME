package occurrence

import (
	"fmt"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/taxonomy"
)

type BaseIdentification struct {
	Taxon        taxonomy.Taxon            `gel:"taxon" json:"taxon"`
	IdentifiedOn OptionalDateWithPrecision `gel:"identified_on" json:"identified_on,omitzero"`
	Confer       bool                      `gel:"confer" json:"confer"`
	Addendum     geltypes.OptionalStr      `gel:"addendum" json:"addendum,omitzero"`
}

type Identification struct {
	BaseIdentification `gel:"$inline"`
	ID                 geltypes.UUID `gel:"id" json:"id" format:"uuid"`
	IdentifiedBy       []string      `gel:"identified_by" json:"identified_by,omitempty"`
	Meta               people.Meta   `gel:"meta" json:"meta"`
}

type IdentificationInput struct {
	Taxon        string                                       `gel:"taxon" json:"taxon"`
	IdentifiedBy []string                                     `gel:"identified_by" json:"identified_by,omitzero"`
	IdentifiedOn models.OptionalInput[DateWithPrecisionInput] `gel:"identified_on" json:"identified_on,omitzero"`
	Confer       bool                                         `gel:"confer" json:"confer,omitzero"`
	Addendum     models.OptionalInput[string]                 `gel:"addendum" json:"addendum,omitzero"`
}

// func (id *IdentificationInput) WithPersonAliases(aliases map[string]string) IdentificationInput {
// 	if p, ok := id.IdentifiedBy.Get(); ok {
// 		if alias, ok := aliases[p]; ok {
// 			id.IdentifiedBy.SetValue(alias)
// 		}
// 	}
// 	return *id
// }

type IdentificationUpdate struct {
	Taxon        models.OptionalInput[string]                 `gel:"taxon" json:"taxon,omitempty"`
	IdentifiedBy models.OptionalNull[[]string]                `gel:"identified_by" json:"identified_by,omitempty"`
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
			"identified_by": `<str>json_array_unpack(data['identified_by'])`,
			"identified_on": `#edgeql
				date::from_json_with_precision(<json>data['identified_on'])
			`,
		},
	}
}
