package occurrence

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/references"
	"github.com/lsdch/biome/models/sequences"
	"github.com/lsdch/biome/models/taxonomy"

	"github.com/edgedb/edgedb-go"
)

type LegacySeqID struct {
	ID            int32  `edgedb:"id" json:"id"`
	Code          string `edgedb:"code" json:"code"`
	AlignmentCode string `edgedb:"alignment_code" json:"alignment_code"`
}

type SequenceInner struct {
	CodeIdentifier `edgedb:"$inline" json:",inline"`
	Label          edgedb.OptionalStr           `edgedb:"label" json:"label"`
	Sequence       edgedb.OptionalStr           `edgedb:"sequence" json:"sequence,omitempty"`
	Gene           sequences.Gene               `edgedb:"gene" json:"gene"`
	Legacy         models.Optional[LegacySeqID] `edgedb:"legacy" json:"legacy,omitempty"`
	Category       OccurrenceCategory           `edgedb:"category" json:"category"`
	// Comments edgedb.OptionalStr           `edgedb:"comments" json:"comments,omitempty"`
}

type SequenceInnerInput struct {
	Code     string                            `json:"code"`
	Label    models.OptionalInput[string]      `json:"label,omitempty"`
	Sequence models.OptionalInput[string]      `json:"sequence,omitempty"`
	Gene     string                            `json:"gene"`
	Legacy   models.OptionalInput[LegacySeqID] `json:"legacy,omitempty"`
	Comments models.OptionalInput[string]      `json:"comments,omitempty"`
}

type ExtSeqSpecifics[BioMat any] struct {
	Origin             sequences.ExtSeqOrigin           `edgedb:"origin" json:"origin"`
	ReferencedIn       []sequences.SeqReference         `edgedb:"referenced_in" json:"referenced_in,omitempty"`
	SpecimenIdentifier string                           `edgedb:"specimen_identifier" json:"specimen_identifier"`
	OriginalTaxon      edgedb.OptionalStr               `edgedb:"original_taxon" json:"original_taxon"`
	SourceSample       models.Optional[BioMat]          `edgedb:"source_sample" json:"source_sample"`
	PublishedIn        []references.OccurrenceReference `edgedb:"published_in" json:"published_in,omitempty"`
}

type GenericSequence[SamplingType any] struct {
	ID                              edgedb.UUID `edgedb:"id" json:"id" format:"uuid"`
	GenericOccurrence[SamplingType] `edgedb:"$inline" json:",inline"`
	SequenceInner                   `edgedb:"$inline" json:",inline"`
	Comments                        edgedb.OptionalStr                            `edgedb:"comments" json:"comments,omitempty"`
	Event                           EventInner                                    `edgedb:"event" json:"event"`
	External                        models.Optional[ExtSeqSpecifics[BioMaterial]] `edgedb:"external" json:"external,omitempty"`
	Meta                            people.Meta                                   `edgedb:"meta" json:"meta"`
}

type Sequence GenericSequence[SamplingInner]

type SequenceWithDetails GenericSequence[Sampling]

func GetSequence(db edgedb.Executor, code string) (seq SequenceWithDetails, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select seq::Sequence {
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
					specimen_identifier,
					original_taxon,
					published_in: { ** },
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
			select seq::Sequence {
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
					source_sample : {
						[is occurrence::BioMaterial].*,
						identification: { ** }
					}
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
	Occurrence                        `edgedb:"$inline" json:",inline"`
	SequenceInner                     `edgedb:"$inline" json:",inline"`
	ExtSeqSpecifics[BioMaterialInner] `edgedb:"$inline" json:",inline"`
	Meta                              people.Meta `edgedb:"meta" json:"meta"`
}

func (s ExternalSequence) AsOccurrence() OccurrenceWithCategory {
	return OccurrenceWithCategory{
		Occurrence:        s.Occurrence,
		Category:          External,
		OccurrenceElement: SequenceElement,
	}
}

type ExternalSequenceInput struct {
	SequenceInnerInput `json:",inline"`
	Origin             sequences.ExtSeqOrigin                `json:"origin"`
	PublishedIn        []references.OccurrenceReferenceInput `json:"published_in,omitempty"`
	ReferencedIn       []sequences.SeqReferenceInput         `json:"referenced_in,omitempty"`
	SpecimenIdentifier string                                `json:"specimen_identifier"`
	OriginalTaxon      models.OptionalInput[string]          `json:"original_taxon,omitempty"`
	SourceSample       models.OptionalInput[string]          `json:"source_sample,omitempty"`
	Identification     IdentificationInput                   `json:"identification"`
}

func (i *ExternalSequenceInput) UseSamplingCode(samplingCode string) {
	i.Code = fmt.Sprintf("%s[%s]%s|%s",
		taxonomy.TaxonCode(i.Identification.Taxon),
		samplingCode,
		i.SpecimenIdentifier,
		strings.ToLower(string(i.Origin)),
	)
}

func (i ExternalSequenceInput) Save(e edgedb.Executor, samplingID edgedb.UUID) (created ExternalSequence, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$1,
			select (insert seq::ExternalSequence {
				sampling := (<events::Sampling><uuid>$0),
				code := <str>data['code'],
				label := <str>json_get(data, 'label'),
				sequence := <str>json_get(data, 'sequence'),
				gene := seq::geneByCode(<str>data['gene']),
				legacy := <tuple<id: int32, code: str, alignment_code: str>>json_get(data, 'legacy'),
				origin := <seq::ExtSeqOrigin>json_get(data, 'origin'),
				published_in := (
					with pubs := json_array_unpack(json_get(data, 'published_in'))
					select assert_distinct((for p in pubs union (
						select references::Article {
							@original_source := <bool>json_get(p, 'original')
						} filter .code = <str>p['code']
					)))
				),
				identification := (
					with identification := data['identification']
					insert occurrence::Identification {
						identified_by := people::personByAlias(<str>identification['identified_by']),
						identified_on := date::from_json_with_precision(identification['identified_on']),
						taxon := taxonomy::taxonByName(<str>identification['taxon']),
					}
				),
				referenced_in := (
          for ref in json_array_unpack(json_get(data, 'referenced_in'))
					insert seq::SeqReference {
            db := seq::seqDbByCode(<str>ref['db']),
            accession := <str>ref['accession'],
            is_origin := <bool>json_get(ref, 'is_origin'),
          }
				),
				specimen_identifier := <str>json_get(data, 'specimen_identifier'),
				original_taxon := <str>json_get(data, 'original_taxon'),
				source_sample := (
					with source_sample := <str>json_get(data, 'source_sample')
					select if exists source_sample
					then occurrence::externalBiomatByCode(source_sample)
					else <occurrence::ExternalBioMat>{}
				),
			}) {
				**,
				referenced_in: { ** },
				source_sample: { id, code, category, is_type, comments },
        identification: { ** }
			}
		`, &created, samplingID, data)
	return
}
