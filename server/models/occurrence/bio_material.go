package occurrence

import (
	"context"
	"encoding/json"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/references"
	"github.com/lsdch/biome/models/sequences"
	"github.com/lsdch/biome/models/specimen"
	"github.com/lsdch/biome/models/taxonomy"

	"github.com/edgedb/edgedb-go"
)

type SpecimenVoucher struct {
	Collection edgedb.OptionalStr `edgedb:"in_collection" json:"collection"`
	Item       []string           `edgedb:"item_vouchers" json:"vouchers"`
}

type ExternalBioMatSpecific struct {
	// ID                 edgedb.UUID        `edgedb:"id" json:"id" format:"uuid"`
	OriginalSource     models.Optional[references.DataSource] `edgedb:"original_source" json:"original_source,omitempty"`
	OriginalLink       edgedb.OptionalStr                     `edgedb:"original_link" json:"original_link,omitempty"`
	OriginalTaxon      edgedb.OptionalStr                     `edgedb:"original_taxon" json:"original_taxon,omitempty"`
	Quantity           specimen.Quantity                      `edgedb:"quantity" json:"quantity"`
	ContentDescription edgedb.OptionalStr                     `edgedb:"content_description" json:"content_description,omitempty"`
	Archive            SpecimenVoucher                        `edgedb:"$inline" json:"archive"`
	Comments           edgedb.OptionalStr                     `edgedb:"comments" json:"comments"`
	Content            []ExternalBioMatContent                `edgedb:"content" json:"content,omitempty"`
}

// ExternalBioMatSequence represents a sequence of an external biomaterial.
// It is intended to be embedded in external bio material occurrence details.
type ExternalBioMatSequence struct {
	ID             edgedb.UUID `edgedb:"id" json:"id" format:"uuid"`
	SequenceInner  `edgedb:"$inline" json:",inline"`
	Category       OccurrenceCategory       `edgedb:"category" json:"category"`
	Origin         sequences.ExtSeqOrigin   `edgedb:"origin" json:"origin"`
	ReferencedIn   []sequences.SeqReference `edgedb:"referenced_in" json:"referenced_in,omitempty"`
	Identification Identification           `edgedb:"identification" json:"identification"`
	Comments       edgedb.OptionalStr       `edgedb:"comments" json:"comments"`
	PublishedIn    []references.Article     `edgedb:"published_in" json:"published_in,omitempty"`
	// SourceSample            `edgedb:"source_sample" json:"source_sample"`
	AccessionNumber    edgedb.OptionalStr `edgedb:"accession_number" json:"accession_number"`
	SpecimenIdentifier string             `edgedb:"specimen_identifier" json:"specimen_identifier"`
	OriginalTaxon      edgedb.OptionalStr `edgedb:"original_taxon" json:"original_taxon"`
}

type ExternalBioMatContent struct {
	Specimen  string                   `edgedb:"specimen" json:"specimen"`
	Sequences []ExternalBioMatSequence `edgedb:"sequences" json:"sequences"`
}

type BioMaterialInner struct {
	ID             edgedb.UUID `edgedb:"id" json:"id" format:"uuid"`
	CodeIdentifier `edgedb:"$inline" json:",inline"`
	Category       OccurrenceCategory `edgedb:"category" json:"category"`
	IsType         bool               `edgedb:"is_type" json:"is_type"`
	Comments       edgedb.OptionalStr `edgedb:"comments" json:"comments,omitempty"`
}

type GenericBioMaterial[SamplingType any] struct {
	GenericOccurrence[SamplingType] `edgedb:"$inline" json:",inline"`
	BioMaterialInner                `edgedb:"$inline" json:",inline"`
	HasSequences                    bool                                    `edgedb:"has_sequences" json:"has_sequences"`
	IsHomogenous                    bool                                    `edgedb:"is_homogenous" json:"is_homogenous"`
	IsCongruent                     bool                                    `edgedb:"is_congruent" json:"is_congruent"`
	SequenceConsensus               models.Optional[taxonomy.Taxon]         `edgedb:"seq_consensus" json:"seq_consensus,omitempty"`
	External                        models.Optional[ExternalBioMatSpecific] `edgedb:"external" json:"external,omitempty"`
	Meta                            people.Meta                             `edgedb:"meta" json:"meta"`
}

type BioMaterial GenericBioMaterial[SamplingInner]

func (b BioMaterial) AsOccurrence() OccurrenceWithCategory {
	return OccurrenceWithCategory{
		Occurrence: Occurrence{
			ID:                b.BioMaterialInner.ID,
			GenericOccurrence: b.GenericOccurrence,
		},
		Category:          b.Category,
		OccurrenceElement: BioMaterialElement,
	}
}

type BioMaterialWithDetails struct {
	GenericBioMaterial[Sampling] `edgedb:"$inline" json:",inline"`
	Event                        EventInner `edgedb:"event" json:"event"`
}

func (b BioMaterialWithDetails) AsOccurrence() OccurrenceWithCategory {
	return OccurrenceWithCategory{
		Occurrence: Occurrence{
			ID: b.BioMaterialInner.ID,
			GenericOccurrence: GenericOccurrence[SamplingInner]{
				Sampling:       b.GenericOccurrence.Sampling.SamplingInner,
				Identification: b.GenericOccurrence.Identification,
			},
			Comments: b.Comments,
		},
		Category:          b.Category,
		OccurrenceElement: BioMaterialElement,
	}
}

func GetBioMaterial(db edgedb.Executor, code string) (biomat BioMaterialWithDetails, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
		with module occurrence
		select BioMaterialWithType {
			**,
			sampling: {
				*,
				target_taxa: { * },
				fixatives: { * },
				methods: { * },
				habitats: { * },
				samples: { **, identification: { **, identified_by: { * } } },
				occurring_taxa: { * }
			},
			event := .sampling.event { *, site: {name, code} },
			identification: { ** },
			external: {
				content := (
					select (group .sequences by .specimen_identifier) {
						specimen := .key.specimen_identifier,
						sequences := .elements {
							*,
							referenced_in: { ** },
							gene: { ** },
							identification: { ** }
						}
					}
				),
				original_source,
				original_link,
				in_collection,
				item_vouchers,
				quantity,
				content_description
			}
		} filter .code = <str>$0
	`,
		&biomat, code)
	return biomat, err
}

func ListBioMaterials(db edgedb.Executor) ([]BioMaterialWithDetails, error) {
	var items = []BioMaterialWithDetails{}
	err := db.Query(context.Background(),
		`#edgeql
			with module occurrence
			select BioMaterialWithType {
        **,
				event := .sampling.event { *, site: {name, code} },
				identification: { **, identified_by: { * } },
        external: {
					original_source,
          original_link,
          in_collection,
          item_vouchers,
          quantity,
          content_description
        }
      }
		`,
		&items)
	return items, err
}

func DeleteBioMaterial(db edgedb.Executor, code string) (deleted BioMaterialWithDetails, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
		with module occurrence
			select (
				delete BioMaterial filter .code = <str>$0
			) {
        **,
				required has_sequences := (
					exists ([is ExternalBioMat].sequences ?? [is InternalBioMat].specimens.sequences)
				),
				required is_homogenous := (
					assert_exists(
						([is ExternalBioMat].is_homogenous ?? [is InternalBioMat].is_homogenous),
						message := "Failed to determined BioMaterial homogeneity"
					)
				),
				required is_congruent := (
					assert_exists(
						([is ExternalBioMat].is_congruent ?? [is InternalBioMat].is_congruent),
						message := "Failed to determined BioMaterial congruence"
					)
				),
				seq_consensus := (
					[is ExternalBioMat].seq_consensus ?? [is InternalBioMat].seq_consensus
				) { * },
				event := .sampling.event { *, site: {name, code} },
				identification: { **, identified_by: { * } },
        external:= [is occurrence::ExternalBioMat]{
					original_source,
          original_link,
          in_collection,
          item_vouchers,
          quantity,
          content_description
        }
      }
		`,
		&deleted, code)
	return
}

type BioMaterialInput struct {
	OccurrenceInnerInput `edgedb:"$inline" json:",inline"`
	Code                 models.OptionalInput[string] `edgedb:"code" json:"code,omitempty"`
	IsType               models.OptionalInput[bool]   `edgedb:"is_type" json:"is_type,omitempty"`
}

func (i *BioMaterialInput) SetCode(code string) {
	i.Code.SetValue(code)
}

func (i *BioMaterialInput) UseSamplingCode(samplingCode string) {
	i.SetCode(i.OccurrenceInnerInput.Code(samplingCode))
}

func (i BioMaterialInput) GetCode(db edgedb.Executor) (string, error) {
	if i.Code.IsSet {
		return i.Code.Value, nil
	} else {
		return i.GenerateCode(db)
	}
}

type BioMaterialUpdate struct {
	OccurrenceUpdate `edgedb:"$inline" json:",inline"`
	Code             models.OptionalInput[string] `edgedb:"code" json:"code,omitempty"`
	IsType           models.OptionalInput[bool]   `edgedb:"is_type" json:"is_type,omitempty"`
}

type InternalBioMatInput struct {
	BioMaterialInput `edgedb:"$inline" json:",inline"`
	// TODO: Internal-specific fields
}

func (i InternalBioMatInput) Save(e edgedb.Executor, samplingID edgedb.UUID) (created BioMaterialWithDetails, err error) {
	data, _ := json.Marshal(i)
	code, err := i.GetCode(e)
	if err != nil {
		return
	}
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$1,
			identification := data['identification'],
			select (insert occurrence::InternalBioMat {
				code := <str>$2,
				identification := (
					insert occurrence::Identification {
						taxon := taxonomy::taxonByName(<str>identification['taxon']),
						identified_by := people::personByAlias(<str>identification['identified_by']),
						identified_on := date::from_json_with_precision(identification['identified_on']),
					}
				),
				sampling := assert_exists((
						select (<events::Sampling><uuid>$0)
					),
					message := "Failed to find sampling with ID: " ++ <str><uuid>$0
				),
				is_type := <bool>json_get(data, 'is_type') ?? false,
				published_in := (
					with pubs := json_array_unpack(json_get(data, 'published_in'))
					select assert_distinct(
						(for p in pubs union (
							select references::Article {
								@original_source := <bool>json_get(p, 'original')
							} filter .code = <str>p['code']
						)),
						message := "Duplicate publication references: " ++ to_str(pubs)
					)
				)
			}) {
				*,
				sampling: {
					*,
					target_taxa: { * },
					fixatives: { * },
					methods: { * },
					habitats: { * },
					samples: { **, identification: { ** } },
					occurring_taxa: { * }
				},
				published_in: { *, @original_source },
				event := .sampling.event { *, site: {name, code} },
				identification: { **, identified_by: { * } },
				meta: { * }
			}
		`, &created, samplingID, data, code)
	return
}

type ExternalBioMatOccurrenceInput struct {
	Sampling            edgedb.UUID `edgedb:"sampling" json:"sampling"`
	ExternalBioMatInput `edgedb:"$inline" json:",inline"`
}

func (i ExternalBioMatOccurrenceInput) Save(e edgedb.Executor) (created BioMaterialWithDetails, err error) {
	return i.ExternalBioMatInput.Save(e, i.Sampling)
}

type ExternalBioMatInput struct {
	BioMaterialInput   `edgedb:"$inline" json:",inline"`
	OriginalSource     models.OptionalInput[string] `edgedb:"original_source" json:"original_source,omitempty"`
	OriginalLink       models.OptionalInput[string] `edgedb:"original_link" json:"original_link,omitempty"`
	Quantity           specimen.Quantity            `edgedb:"quantity" json:"quantity"`
	ContentDescription models.OptionalInput[string] `edgedb:"content_description" json:"content_description,omitempty"`
	Collection         models.OptionalInput[string] `edgedb:"in_collection" json:"collection,omitempty"`
	Item               []string                     `edgedb:"item_vouchers" json:"vouchers,omitempty"`
	Comments           models.OptionalInput[string] `edgedb:"comments" json:"comments,omitempty"`
}

func (i ExternalBioMatInput) Save(e edgedb.Executor, samplingID edgedb.UUID) (created BioMaterialWithDetails, err error) {
	data, _ := json.Marshal(i)
	code, err := i.GetCode(e)
	if err != nil {
		return
	}
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$1,
			identification := data['identification'],
			select (insert occurrence::ExternalBioMat {
				code := <str>$2,
				original_source := (
					with src := <str>json_get(data, 'original_source')
					select (if exists src then default::get_vocabulary(src)[is references::DataSource] else <references::DataSource>{})
				),
				original_link := <str>json_get(data, 'original_link'),
				quantity := <occurrence::QuantityType>json_get(data, 'quantity'),
				content_description := <str>json_get(data, 'content_description'),
				in_collection := <str>json_get(data, 'collection'),
				item_vouchers := <str>json_array_unpack(json_get(data, 'item_vouchers')),
				comments := <str>json_get(data, 'comments'),
				published_in := (
					with pubs := json_array_unpack(json_get(data, 'published_in'))
					select assert_distinct(
						(for p in pubs union (
							select references::Article {
								@original_source := <bool>json_get(p, 'original')
							} filter .code = <str>p['code']
						)),
						message := "Duplicate publication references: " ++ to_str(pubs)
					)
				),
        identification := (
          insert occurrence::Identification {
            taxon := taxonomy::taxonByName(<str>identification['taxon']),
            identified_by := people::personByAlias(<str>identification['identified_by']),
            identified_on := date::from_json_with_precision(identification['identified_on']),
          }
        ),
        sampling := assert_exists(
					(select (<events::Sampling><uuid>$0)),
					message := "Failed to find sampling with ID: " ++ <str><uuid>$0
				),
        is_type := <bool>json_get(data, 'is_type') ?? false,
			}) {
        [is occurrence::BioMaterial].**,
				is_homogenous,
				is_congruent,
				seq_consensus: { * },
				event := .sampling.event { *, site: {name, code} },
				identification: { ** },
        external := [is occurrence::ExternalBioMat]{
          original_link,
          in_collection,
          item_vouchers,
          quantity,
          content_description
        }
      }
		`, &created, samplingID, data, code)
	return
}

type ExternalBioMatUpdate struct {
	BioMaterialUpdate  `edgedb:"$inline" json:",inline"`
	OriginalSource     models.OptionalNull[string]                                `edgedb:"original_source" json:"original_source,omitempty"`
	OriginalLink       models.OptionalNull[string]                                `edgedb:"original_link" json:"original_link,omitempty"`
	OriginalTaxon      models.OptionalNull[string]                                `edgedb:"original_taxon" json:"original_taxon,omitempty"`
	Quantity           models.OptionalInput[specimen.Quantity]                    `edgedb:"quantity" json:"quantity,omitempty"`
	ContentDescription models.OptionalNull[string]                                `edgedb:"content_description" json:"content_description,omitempty"`
	Collection         models.OptionalNull[string]                                `edgedb:"in_collection" json:"collection,omitempty"`
	Item               models.OptionalInput[[]string]                             `edgedb:"item_vouchers" json:"vouchers,omitempty"`
	Comments           models.OptionalNull[string]                                `edgedb:"comments" json:"comments,omitempty"`
	PublishedIn        models.OptionalNull[[]references.OccurrenceReferenceInput] `edgedb:"published_in" json:"published_in"`
}

func (u ExternalBioMatUpdate) Save(e edgedb.Executor, code string) (updated BioMaterialWithDetails, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
      with item := <json>$1,
      select (update occurrence::ExternalBioMat filter .code = <str>$0 set {
        %s
      }) {
        **,
				seq_consensus: { * },
				event := .sampling.event { *, site: {name, code} },
				identification: { **, identified_by: { * } },
        external := [is occurrence::ExternalBioMat]{
          original_link,
          in_collection,
          item_vouchers,
          quantity,
          content_description
        }
      }
    `,
		Mappings: map[string]string{
			"code": "<str>item['code']", // if not explicitly provided, updated code is autogenerated
			"original_source": `#edgeql
				default::get_vocabulary(<str>item['original_source'])
			`,
			"original_link":       "<str>item['original_link']",
			"original_taxon":      "<str>item['original_taxon']",
			"quantity":            "<occurrence::QuantityType>item['quantity']",
			"content_description": "<str>item['content_description']",
			"in_collection":       "<str>item['collection']",
			"item_vouchers":       "<str>json_array_unpack(item['item_vouchers'])",
			"comments":            "<str>item['comments']",
			"is_type":             "<bool>item['is_type']",
			"identification":      u.Identification.Value.UpdateQuery(".identification"),
			"published_in": `#edgeql
					with pubs := json_array_unpack(json_get(item, 'published_in'))
					select assert_distinct(
						(for p in pubs union (
							select references::Article {
								@original_source := <bool>json_get(p, 'original')
							} filter .code = <str>p['code']
						)),
						message := "Duplicate publication references: " ++ to_str(pubs)
					)
			`,
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, code, data)
	return
}
