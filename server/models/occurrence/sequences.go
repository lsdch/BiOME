package occurrence

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/sequences"
	"github.com/sirupsen/logrus"
)

type LegacySeqID struct {
	ID            int32  `gel:"id" json:"id"`
	Code          string `gel:"code" json:"code"`
	AlignmentCode string `gel:"alignment_code" json:"alignment_code"`
}

type SequenceInner struct {
	ID             geltypes.UUID `gel:"id" json:"id" format:"uuid"`
	CodeIdentifier `gel:"$inline" json:",inline"`
	WithCategory   `gel:"$inline" json:",inline"`
	Label          geltypes.OptionalStr         `gel:"label" json:"label,omitempty"`
	Sequence       geltypes.OptionalStr         `gel:"sequence" json:"sequence,omitempty"`
	Gene           sequences.Gene               `gel:"gene" json:"gene"`
	Legacy         models.Optional[LegacySeqID] `gel:"legacy" json:"legacy,omitempty"`
	IsIdentifying  bool                         `gel:"is_identifying" json:"is_identifying"`
	Comments       geltypes.OptionalStr         `gel:"comments" json:"comments,omitempty"`
}

type SequenceInput struct {
	Code          string                            `json:"code"`
	Label         models.OptionalInput[string]      `json:"label,omitempty"`
	Sequence      models.OptionalInput[string]      `json:"sequence,omitempty"`
	Gene          string                            `json:"gene"`
	Legacy        models.OptionalInput[LegacySeqID] `json:"legacy,omitempty"`
	Comments      models.OptionalInput[string]      `json:"comments,omitempty"`
	IsIdentifying models.OptionalInput[bool]        `json:"is_identifying,omitempty"`
	ReferencedIn  []sequences.SeqReferenceInput     `json:"referenced_in,omitempty"`
}

func (seq *SequenceInput) WithCreatedMetadata(c *CreatedMetadata) *SequenceInput {
	for i := range seq.ReferencedIn {
		(&seq.ReferencedIn[i]).WithDataSourceCode(c.DataSources)
	}
	return seq
}

type Sequence struct {
	SequenceInner `gel:"$inline" json:",inline"`
	ReferencedIn  []sequences.SeqReference `gel:"referenced_in" json:"referenced_in,omitempty"`
	Meta          people.Meta              `gel:"meta" json:"meta"`
}

type SequenceListItem struct {
	SequenceInner      `gel:"$inline" json:",inline"`
	Occurrence         GenericOccurrence[SamplingInnerWithSite] `gel:"occurrence" json:"occurrence"`
	Identification     Identification                           `gel:"identification" json:"identification"`
	SpecimenIdentifier geltypes.OptionalStr                     `gel:"specimen_identifier" json:"specimen_identifier,omitempty"`
	Meta               people.Meta                              `gel:"meta" json:"meta"`
}

type AssembledSequenceSpecifics struct {
	AlignmentCode string          `gel:"alignment_code" json:"alignment_code"`
	AssembledBy   []people.Person `gel:"assembled_by" json:"assembled_by,omitempty"`
	// Not implemented yet
	// Chromatograms []sequences.Chromatogram `gel:"chromatograms" json:"chromatograms,omitempty"`
	// Specimen      samples.Specimen          `gel:"specimen" json:"specimen"`
}

// type SequenceListItem GenericSequence[SamplingInnerWithSite]
type SequenceWithDetails struct {
	Sequence           `gel:"$inline" json:",inline"`
	Occurrence         Occurrence[SamplingWithSite] `gel:"occurrence" json:"occurrence"`
	Identification     Identification               `gel:"identification" json:"identification"`
	SpecimenIdentifier string                       `gel:"specimen_identifier" json:"specimen_identifier"`

	Internal models.Optional[AssembledSequenceSpecifics] `gel:"internal" json:"internal,omitempty"`
}

func GetSequence(db geltypes.Executor, code string) (seq SequenceWithDetails, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select seq::GenericSequenceWithDetails {
				**,
				referenced_in: { ** },
				gene: { * },
				identification: { **, identified_by: { * } },
				occurrence: {
					*,
					identification: { **, identified_by: { ** } },
					sampling: {
						*,
						target_taxa: { * },
						fixatives: { * },
						methods: { * },
						habitats: { * },
						occurrences: { * },
						occurring_taxa: { * },
						site: { *, country: { * } },
						meta: { * }
					},
				},
				internal: {
					alignment_code,
					assembled_by: { * },
					# chromatograms: { * },
					# specimen: { *, biomat: { *, sampling: { *, site: { *, country: { * } } } } }
				}
			} filter .code = <str>$0
		`,
		&seq, code)
	return
}

func ListSequences(db geltypes.Executor) ([]SequenceListItem, error) {
	var items = []SequenceListItem{}
	err := db.Query(context.Background(),
		`#edgeql
			select seq::GenericSequenceWithDetails {
				*,
				gene: { * },
				identification: { ** },
				occurrence: {
					id,
					code,
					identification: { **, identified_by: { ** } },
					sampling: {
						id,
						number,
						performed_on,
						site: { *, country: { * } }
					}
				},
				meta: { * }
			}
		`,
		&items)
	return items, err
}

func DeleteSequence(db geltypes.Executor, code string) (deleted SequenceListItem, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
			 delete seq::Sequence filter .code = <str>$0
		 	) {
				*,
				gene: { * },
				required identification := (
					[is AssembledSequence].identification ?? .biomat.identification
				) { **, identified_by: { * } },
				required occurrence := (
					assert_exists(
						([is AssembledSequence].specimen.biomat ?? [is ExternalSequence].biomat ?? {}),
						message := "Failed to find occurrence for seq::Sequence subtype " ++ __source__.__type__.name
					)
				) {
					id,
					code,
					identification: { **, identified_by: { ** } },
					sampling: { id, number, performed_on, site: { *, country: { * } } }
				},
			}
		`,
		&deleted, code)
	return
}

type ExternalSequence struct {
	Sequence           `gel:"$inline" json:",inline"`
	SpecimenIdentifier geltypes.OptionalStr `gel:"specimen_identifier" json:"specimen_identifier,omitempty"`
}

type ExternalSequenceInput struct {
	SequenceInput      `json:",inline"`
	SpecimenIdentifier string `json:"specimen_identifier"`
	// Origin             sequences.ExtSeqOrigin        `json:"origin"`
}

func (i *ExternalSequenceInput) GenerateCode(taxonCode string, samplingCode string) {
	i.Code = fmt.Sprintf("%s[%s]%s|%s",
		taxonCode,
		samplingCode,
		i.SpecimenIdentifier,
		i.Gene,
	)
}

func (i ExternalSequenceInput) Save(e geltypes.Executor, occurrenceCode string) (created ExternalSequence, err error) {
	data, _ := json.Marshal(i)
	logrus.Infof("ExternalSequence: %s", string(data))
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$1,
			select (insert seq::ExternalSequence {
				biomat := (select occurrence::ExternalOccurrence filter .code = <str>$0),
				code := <str>data['code'],
				label := <str>json_get(data, 'label'),
				sequence := <str>json_get(data, 'sequence'),
				gene := seq::geneByCode(<str>data['gene']),
				legacy := <tuple<id: int32, code: str, alignment_code: str>>json_get(data, 'legacy'),
				# origin := <seq::ExtSeqOrigin>json_get(data, 'origin'),
				referenced_in := (
          for ref in json_array_unpack(json_get(data, 'referenced_in'))
					insert references::SeqReference {
            db := references::dataSourceByCode(<str>ref['db']),
            accession := <str>ref['accession'],
          }
				),
				specimen_identifier := <str>json_get(data, 'specimen_identifier'),
			}) {
				*,
				gene: { * },
				referenced_in: { ** },
			}
		`, &created, occurrenceCode, data)
	return
}
