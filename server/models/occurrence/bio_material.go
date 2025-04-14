package occurrence

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/references"
	"github.com/lsdch/biome/models/sequences"
	"github.com/lsdch/biome/models/specimen"
	"github.com/lsdch/biome/models/taxonomy"
	"github.com/sirupsen/logrus"
)

type SpecimenVoucher struct {
	Collection geltypes.OptionalStr `gel:"in_collection" json:"collection"`
	Item       []string             `gel:"item_vouchers" json:"vouchers"`
}

type ExternalBioMatSpecific struct {
	// ID                 geltypes.UUID        `gel:"id" json:"id" format:"uuid"`
	OriginalSource     models.Optional[references.DataSource] `gel:"original_source" json:"original_source,omitempty"`
	OriginalLink       geltypes.OptionalStr                   `gel:"original_link" json:"original_link,omitempty"`
	OriginalTaxon      geltypes.OptionalStr                   `gel:"original_taxon" json:"original_taxon,omitempty"`
	Quantity           specimen.Quantity                      `gel:"quantity" json:"quantity"`
	ContentDescription geltypes.OptionalStr                   `gel:"content_description" json:"content_description,omitempty"`
	Archive            SpecimenVoucher                        `gel:"$inline" json:"archive"`
	Comments           geltypes.OptionalStr                   `gel:"comments" json:"comments"`
	Content            []ExternalBioMatContent                `gel:"content" json:"content,omitempty"`
}

// ExternalBioMatSequence represents a sequence of an external biomaterial.
// It is intended to be embedded in external bio material occurrence details.
type ExternalBioMatSequence struct {
	ID             geltypes.UUID `gel:"id" json:"id" format:"uuid"`
	SequenceInner  `gel:"$inline" json:",inline"`
	Category       OccurrenceCategory       `gel:"category" json:"category"`
	Origin         sequences.ExtSeqOrigin   `gel:"origin" json:"origin"`
	ReferencedIn   []sequences.SeqReference `gel:"referenced_in" json:"referenced_in,omitempty"`
	Identification Identification           `gel:"identification" json:"identification"`
	Comments       geltypes.OptionalStr     `gel:"comments" json:"comments"`
	PublishedIn    []references.Article     `gel:"published_in" json:"published_in,omitempty"`
	// SourceSample            `gel:"source_sample" json:"source_sample"`
	AccessionNumber    geltypes.OptionalStr `gel:"accession_number" json:"accession_number"`
	SpecimenIdentifier string               `gel:"specimen_identifier" json:"specimen_identifier"`
	OriginalTaxon      geltypes.OptionalStr `gel:"original_taxon" json:"original_taxon"`
}

type ExternalBioMatContent struct {
	Specimen  string                   `gel:"specimen" json:"specimen"`
	Sequences []ExternalBioMatSequence `gel:"sequences" json:"sequences"`
}

type BioMaterialInner struct {
	ID             geltypes.UUID `gel:"id" json:"id" format:"uuid"`
	CodeIdentifier `gel:"$inline" json:",inline"`
	Category       OccurrenceCategory   `gel:"category" json:"category"`
	IsType         bool                 `gel:"is_type" json:"is_type"`
	Comments       geltypes.OptionalStr `gel:"comments" json:"comments,omitempty"`
}

type GenericBioMaterial[SamplingType any] struct {
	GenericOccurrence[SamplingType] `gel:"$inline" json:",inline"`
	BioMaterialInner                `gel:"$inline" json:",inline"`
	HasSequences                    bool                                    `gel:"has_sequences" json:"has_sequences"`
	IsHomogenous                    bool                                    `gel:"is_homogenous" json:"is_homogenous"`
	IsCongruent                     bool                                    `gel:"is_congruent" json:"is_congruent"`
	SequenceConsensus               models.Optional[taxonomy.Taxon]         `gel:"seq_consensus" json:"seq_consensus,omitempty"`
	External                        models.Optional[ExternalBioMatSpecific] `gel:"external" json:"external,omitempty"`
	Meta                            people.Meta                             `gel:"meta" json:"meta"`
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
	GenericBioMaterial[Sampling] `gel:"$inline" json:",inline"`
	Event                        EventInner `gel:"event" json:"event"`
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

func GetBioMaterial(db geltypes.Executor, code string) (biomat BioMaterialWithDetails, err error) {
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
			event := .sampling.event { *, site: { *, country: { * } } },
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

type ListBioMaterialOptions struct {
	models.Pagination `json:",inline"`
	models.Sorting    `json:",inline"`
	models.Filter     `json:",inline"`
	Category          models.OptionalInput[OccurrenceCategory] `query:"category" json:"category,omitzero"`
	Taxon             models.OptionalInput[string]             `query:"taxon" json:"taxon,omitzero"`
	HasSequences      models.OptionalInput[bool]               `query:"has_sequences" json:"has_sequences,omitzero"`
	IsType            models.OptionalInput[bool]               `query:"is_type" json:"is_type,omitzero"`
}

func (o ListBioMaterialOptions) Options() ListBioMaterialOptions {
	return o
}

func ListBioMaterials(db geltypes.Executor, opts ListBioMaterialOptions) ([]BioMaterialWithDetails, error) {
	params, _ := json.Marshal(opts)
	logrus.Debugf("Params: %s", string(params))
	var items = []BioMaterialWithDetails{}
	err := db.Query(context.Background(),
		`#edgeql
			with module occurrence,
				params := <json>$0,
				search_term := <str>json_get(params, 'search'),
				category := <OccurrenceCategory>json_get(params, 'category'),
				taxon := <str>json_get(params, 'taxon'),
				has_sequences := <bool>json_get(params, 'has_sequences'),
				is_type := <bool>json_get(params, 'is_type'),
				owner := <uuid>json_get(params, 'owner'),
			select BioMaterialWithType {
        **,
				event := .sampling.event { *, site: { *, country: { * } } },
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
			filter (
				(.code ilike '%' ++ search_term ++ '%' if exists search_term else true) and
				(.category = category if exists category else true) and
				(.identification.taxon.name ilike '%' ++ taxon ++ '%' if exists taxon else true) and
				(.has_sequences = has_sequences if exists has_sequences else true) and
				(.is_type = is_type if exists is_type else true) and
				(.meta.created_by_user.id = owner if exists owner else true)
			)
			order by .identification.identified_on.date desc
			offset <optional int64>json_get(params, 'offset')
			limit <optional int64>json_get(params, 'limit');
		`,
		&items, params)
	return items, err
}

func DeleteBioMaterial(db geltypes.Executor, code string) (deleted BioMaterialWithDetails, err error) {
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
				event := .sampling.event { *, site: { *, country: { * } } },
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
	OccurrenceInnerInput `gel:"$inline" json:",inline"`
	Code                 models.OptionalInput[string] `gel:"code" json:"code,omitempty" doc:"Unique code identifier for the bio material. Generated from taxon and sampling if not provided." example:"Genus_sp[SITE|2001-01]"`
	IsType               models.OptionalInput[bool]   `gel:"is_type" json:"is_type,omitempty" doc:"Flag indicating if the bio material is a type specimen, i.e. the reference specimen used to describe a new species."`
}

func (i *BioMaterialInput) SetCode(code string) {
	i.Code.SetValue(code)
}

type BioMaterialUpdate struct {
	OccurrenceUpdate `gel:"$inline" json:",inline"`
	Code             models.OptionalInput[string] `gel:"code" json:"code,omitempty"`
	IsType           models.OptionalInput[bool]   `gel:"is_type" json:"is_type,omitempty"`
}

type InternalBioMatInput struct {
	BioMaterialInput `gel:"$inline" json:",inline"`
	// TODO: Internal-specific fields
}

func (i *InternalBioMatInput) WithCreatedMetadata(c CreatedMetadata) InternalBioMatInput {
	i.OccurrenceInnerInput.WithCreatedMetadata(c)
	return *i
}

func (i InternalBioMatInput) Save(e geltypes.Executor, samplingID geltypes.UUID) (created BioMaterialWithDetails, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with
				sampling := (
					assert_exists(
						(select (<events::Sampling><uuid>$0)),
						message := "Failed to find sampling with ID: " ++ <str><uuid>$0
					)
				),
				data := <json>$1,
				identification := data['identification'],
				taxon := taxonomy::taxonByName(<str>identification['taxon']),
			select (insert occurrence::InternalBioMat {
				code := <str>json_get(data, 'code') ?? occurrence::biomat_code(taxon, sampling),
				identification := (
					insert occurrence::Identification {
						taxon := taxon,
						identified_by := people::personByAlias(<str>identification['identified_by']),
						identified_on := date::from_json_with_precision(identification['identified_on']),
					}
				),
				sampling := sampling,
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
				event := .sampling.event { *, site: { *, country: { * } } },
				identification: { **, identified_by: { * } },
				meta: { * }
			}
		`, &created, samplingID, data)
	return
}

type ExternalBioMatOccurrenceInput struct {
	Sampling            geltypes.UUID `gel:"sampling" json:"sampling"`
	ExternalBioMatInput `gel:"$inline" json:",inline"`
}

func (i ExternalBioMatOccurrenceInput) Save(e geltypes.Executor) (created BioMaterialWithDetails, err error) {
	return i.ExternalBioMatInput.Save(e, i.Sampling)
}

type ExternalBioMatInput struct {
	BioMaterialInput   `gel:"$inline" json:",inline"`
	OriginalSource     models.OptionalInput[string] `gel:"original_source" json:"original_source,omitempty"`
	OriginalLink       models.OptionalInput[string] `gel:"original_link" json:"original_link,omitempty"`
	Quantity           specimen.Quantity            `gel:"quantity" json:"quantity"`
	ContentDescription models.OptionalInput[string] `gel:"content_description" json:"content_description,omitempty" doc:"Description of the content of the bio material" example:"2 females, 1 juvenile male"`
	Collection         models.OptionalInput[string] `gel:"in_collection" json:"collection,omitempty"`
	Item               []string                     `gel:"item_vouchers" json:"vouchers,omitempty"`
	Comments           models.OptionalInput[string] `gel:"comments" json:"comments,omitempty"`
}

func (bm *ExternalBioMatInput) WithCreatedMetadata(c CreatedMetadata) ExternalBioMatInput {
	bm.BioMaterialInput.WithCreatedMetadata(c)
	if dataSource, ok := bm.OriginalSource.Get(); ok {
		if s, ok := c.DataSources[dataSource]; ok {
			bm.OriginalSource = (&bm.OriginalSource).SetValue(s)
		}
	}
	return *bm
}

func (i ExternalBioMatInput) Save(e geltypes.Executor, samplingID geltypes.UUID) (created BioMaterialWithDetails, err error) {
	data, _ := json.Marshal(i)
	logrus.Infof("Creating ExternalBioMat with args: %s", string(data))
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with
				sampling := (assert_exists(
					(select (<events::Sampling><uuid>$0)),
					message := "Failed to find sampling with ID: " ++ <str><uuid>$0
				)),
				data := <json>$1,
				identification := data['identification'],
				taxon := taxonomy::taxonByName(<str>identification['taxon']),
			select (insert occurrence::ExternalBioMat {
				code := <str>json_get(data, 'code') ?? occurrence::biomat_code(taxon, sampling),
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
            taxon := taxon,
            identified_by := people::personByAlias(<str>identification['identified_by']),
            identified_on := date::from_json_with_precision(identification['identified_on']),
          }
        ),
        sampling := sampling,
        is_type := <bool>json_get(data, 'is_type') ?? false,
			}) {
        [is occurrence::BioMaterial].**,
				is_homogenous,
				is_congruent,
				seq_consensus: { * },
				event := .sampling.event { *, site: { *, country: { * } } },
				identification: { ** },
        external := [is occurrence::ExternalBioMat]{
          original_link,
          in_collection,
          item_vouchers,
          quantity,
          content_description
        }
      }
		`, &created, samplingID, data)
	return
}

type ExternalBioMatUpdate struct {
	BioMaterialUpdate  `gel:"$inline" json:",inline"`
	OriginalSource     models.OptionalNull[string]                                `gel:"original_source" json:"original_source,omitempty"`
	OriginalLink       models.OptionalNull[string]                                `gel:"original_link" json:"original_link,omitempty"`
	OriginalTaxon      models.OptionalNull[string]                                `gel:"original_taxon" json:"original_taxon,omitempty"`
	Quantity           models.OptionalInput[specimen.Quantity]                    `gel:"quantity" json:"quantity,omitempty"`
	ContentDescription models.OptionalNull[string]                                `gel:"content_description" json:"content_description,omitempty"`
	Collection         models.OptionalNull[string]                                `gel:"in_collection" json:"collection,omitempty"`
	Item               models.OptionalInput[[]string]                             `gel:"item_vouchers" json:"vouchers,omitempty"`
	Comments           models.OptionalNull[string]                                `gel:"comments" json:"comments,omitempty"`
	PublishedIn        models.OptionalNull[[]references.OccurrenceReferenceInput] `gel:"published_in" json:"published_in"`
}

func (u ExternalBioMatUpdate) Save(e geltypes.Executor, code string) (updated BioMaterialWithDetails, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
      with item := <json>$1,
      select (update occurrence::ExternalBioMat filter .code = <str>$0 set {
        %s
      }) {
        **,
				seq_consensus: { * },
				event := .sampling.event { *, site: { *, country: { * } } },
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
