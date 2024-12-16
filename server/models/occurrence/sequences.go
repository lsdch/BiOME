package occurrence

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/people"
	"darco/proto/models/references"
	"darco/proto/models/sequences"

	"github.com/edgedb/edgedb-go"
)

type LegacySeqID struct {
	ID            int32  `edgedb:"id" json:"id"`
	Code          string `edgedb:"code" json:"code"`
	AlignmentCode string `edgedb:"alignment_code" json:"alignment_code"`
}

type SequenceInner struct {
	Code     string                       `edgedb:"code" json:"code"`
	Label    edgedb.OptionalStr           `edgedb:"label" json:"label"`
	Sequence edgedb.OptionalStr           `edgedb:"sequence" json:"sequence"`
	Gene     sequences.Gene               `edgedb:"gene" json:"gene"`
	LegacyID models.Optional[LegacySeqID] `edgedb:"legacy" json:"legacy"`
}

type ExtSeqSpecifics struct {
	Origin             sequences.ExtSeqOrigin              `edgedb:"origin" json:"origin"`
	PublishedIn        models.Optional[references.Article] `edgedb:"published_in" json:"published_in,omitempty"`
	ReferencedIn       []sequences.SeqReference            `edgedb:"referenced_in" json:"referenced_in,omitempty"`
	SpecimenIdentifier string                              `edgedb:"specimen_identifier" json:"specimen_identifier"`
	OriginalTaxon      edgedb.OptionalStr                  `edgedb:"original_taxon" json:"original_taxon"`
	SourceSample       models.Optional[BioMaterial]        `edgedb:"source_sample" json:"source_sample"`
}

type GenericSequence[SamplingType any] struct {
	GenericOccurrence[SamplingType] `edgedb:"$inline" json:",inline"`
	SequenceInner                   `edgedb:"$inline" json:",inline"`
	Category                        OccurrenceCategory               `edgedb:"category" json:"category"`
	Event                           EventInner                       `edgedb:"event" json:"event"`
	External                        models.Optional[ExtSeqSpecifics] `edgedb:"external" json:"external,omitempty"`
	Meta                            people.Meta                      `edgedb:"meta" json:"meta"`
}

type Sequence GenericSequence[SamplingInner]

type SequenceWithDetails GenericSequence[Sampling]

func GetSequence(db edgedb.Executor, code string) (seq SequenceWithDetails, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select seq::SequenceWithType {
				**,
				gene: { * },
				required event := .sampling.event { *, site: {name, code} },
				sampling: {
					*,
					target_taxa: { * },
					fixatives: { * },
					methods: { * },
					habitats: { * },
					samples: { **, identification: { **, identified_by: { * } } },
					occurring_taxa: { * }
				},
				identification: { **, identified_by: { * } },
				external := [is seq::ExternalSequence]{
					origin,
					referenced_in: { ** },
					published_in: { ** },
					specimen_identifier,
					original_taxon,
					source_sample : {
						[is occurrence::BioMaterial].*,
						identification: { ** }
					}
				}
			} filter .code = <str>$0
		`,
		&seq, code)
	return
}

func ListSequences(db edgedb.Executor) ([]Sequence, error) {
	var items = []Sequence{}
	err := db.Query(context.Background(),
		`#edgeql
			select seq::SequenceWithType {
				**,
				gene: { * },
				required event := .sampling.event { *, site: {name, code} },
				identification: { **, identified_by: { * } },
				external := [is seq::ExternalSequence]{
					origin,
					referenced_in: { ** },
					published_in: { ** },
					specimen_identifier,
					original_taxon,
				}
			}
		`,
		&items)
	return items, err
}

func DeleteSequence(db edgedb.Executor, code string) (deleted Sequence, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
			 delete seq::Sequence filter .code = <str>$0
		 	) {
				**,
				gene: { * },
				required event := .sampling.event { *, site: {name, code} },
				identification: { **, identified_by: { * } },
				external := [is seq::ExternalSequence]{
					origin,
					referenced_in: { ** },
					published_in: { ** },
					specimen_identifier,
					original_taxon,
				}
			}
		`,
		&deleted, code)
	return
}

type ExternalSequence struct {
	Occurrence      `edgedb:"$inline" json:",inline"`
	SequenceInner   `edgedb:"$inline" json:",inline"`
	ExtSeqSpecifics `edgedb:"$inline" json:",inline"`
}
