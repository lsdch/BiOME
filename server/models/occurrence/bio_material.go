package occurrence

import (
	"context"
	"darco/proto/db"
	"darco/proto/models"
	"darco/proto/models/people"
	"darco/proto/models/references"
	"darco/proto/models/sequences"
	"darco/proto/models/specimen"
	"darco/proto/models/taxonomy"
	"encoding/json"
	"time"

	"github.com/edgedb/edgedb-go"
)

type SpecimenVoucher struct {
	Collection edgedb.OptionalStr `edgedb:"in_collection" json:"collection"`
	Item       []string           `edgedb:"item_vouchers" json:"vouchers"`
}

type ExternalBioMatSpecific struct {
	// ID                 edgedb.UUID        `edgedb:"id" json:"id" format:"uuid"`
	OriginalLink       edgedb.OptionalStr      `edgedb:"original_link" json:"original_link,omitempty"`
	OriginalTaxon      edgedb.OptionalStr      `edgedb:"original_taxon" json:"original_taxon,omitempty"`
	Quantity           specimen.Quantity       `edgedb:"quantity" json:"quantity"`
	ContentDescription edgedb.OptionalStr      `edgedb:"content_description" json:"content_description,omitempty"`
	Archive            SpecimenVoucher         `edgedb:"$inline" json:"archive"`
	Comments           edgedb.OptionalStr      `edgedb:"comments" json:"comments"`
	Content            []ExternalBioMatContent `edgedb:"content" json:"content,omitempty"`
}

type ExternalBioMatSequence struct {
	ID             edgedb.UUID `edgedb:"id" json:"id" format:"uuid"`
	SequenceInner  `edgedb:"$inline" json:",inline"`
	Category       OccurrenceCategory       `edgedb:"category" json:"category"`
	Origin         sequences.ExtSeqOrigin   `edgedb:"origin" json:"origin"`
	ReferencedIn   []sequences.SeqReference `edgedb:"referenced_in" json:"referenced_in"`
	Identification Identification           `edgedb:"identification" json:"identification"`
	Comments       edgedb.OptionalStr       `edgedb:"comments" json:"comments"`
	PublishedIn    []references.Article     `edgedb:"published_in" json:"published_in"`
	// SourceSample            `edgedb:"source_sample" json:"source_sample"`
	AccessionNumber    edgedb.OptionalStr `edgedb:"accession_number" json:"accession_number"`
	SpecimenIdentifier string             `edgedb:"specimen_identifier" json:"specimen_identifier"`
	OriginalTaxon      edgedb.OptionalStr `edgedb:"original_taxon" json:"original_taxon"`
}

type ExternalBioMatContent struct {
	Specimen  string                   `edgedb:"specimen" json:"specimen"`
	Sequences []ExternalBioMatSequence `edgedb:"sequences" json:"sequences"`
}

type CodeHistory struct {
	Code string    `edgedb:"code" json:"code"`
	Time time.Time `edgedb:"time" json:"time"`
}

type GenericBioMaterial[SamplingType any] struct {
	GenericOccurrence[SamplingType] `edgedb:"$inline" json:",inline"`
	Code                            string                                  `edgedb:"code" json:"code"`
	CodeHistory                     []CodeHistory                           `edgedb:"code_history" json:"code_history,omitempty"`
	Category                        OccurrenceCategory                      `edgedb:"category" json:"category"`
	IsType                          bool                                    `edgedb:"is_type" json:"is_type"`
	HasSequences                    bool                                    `edgedb:"has_sequences" json:"has_sequences"`
	IsHomogenous                    bool                                    `edgedb:"is_homogenous" json:"is_homogenous"`
	IsCongruent                     bool                                    `edgedb:"is_congruent" json:"is_congruent"`
	SequenceConsensus               models.Optional[taxonomy.Taxon]         `edgedb:"sequence_consensus" json:"sequence_consensus,omitempty"`
	References                      []references.Article                    `edgedb:"published_in" json:"published_in"`
	External                        models.Optional[ExternalBioMatSpecific] `edgedb:"external" json:"external,omitempty"`
	Meta                            people.Meta                             `edgedb:"meta" json:"meta"`
}

type BioMaterial GenericBioMaterial[SamplingInner]

type BioMaterialWithDetails struct {
	GenericBioMaterial[Sampling] `edgedb:"$inline" json:",inline"`
	Event                        EventInner `edgedb:"event" json:"event"`
}

func GetBioMaterial(db edgedb.Executor, code string) (biomat BioMaterialWithDetails, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
		select occurrence::BioMaterialWithType {
			**,
			sequence_consensus: { * },
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
			identification: { **, identified_by: { * } },
			external := [is occurrence::ExternalBioMat]{
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
			select occurrence::BioMaterialWithType {
        **,
				sequence_consensus: { * },
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
		&items)
	return items, err
}

func DeleteBioMaterial(db edgedb.Executor, code string) (deleted BioMaterial, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
				delete occurence::BioMaterialWithType filter .code = <str>$0
			) {
        **,
        external:= [is occurrence::ExternalBioMat]{
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
	OccurrenceInput `edgedb:"$inline" json:",inline"`
	Code            models.OptionalInput[string] `edgedb:"code" json:"code,omitempty"`
	IsType          models.OptionalInput[bool]   `edgedb:"is_type" json:"is_type,omitempty"`
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

type ExternalBioMatInput struct {
	BioMaterialInput   `edgedb:"$inline" json:",inline"`
	OriginalLink       models.OptionalInput[string]   `edgedb:"original_link" json:"original_link,omitempty"`
	OriginalTaxon      models.OptionalInput[string]   `edgedb:"original_taxon" json:"original_taxon,omitempty"`
	Quantity           specimen.Quantity              `edgedb:"quantity" json:"quantity"`
	ContentDescription models.OptionalInput[string]   `edgedb:"content_description" json:"content_description,omitempty"`
	Collection         models.OptionalInput[string]   `edgedb:"in_collection" json:"collection,omitempty"`
	Item               []string                       `edgedb:"item_vouchers" json:"vouchers,omitempty"`
	Comments           models.OptionalInput[string]   `edgedb:"comments" json:"comments,omitempty"`
	PublishedIn        models.OptionalInput[[]string] `edgedb:"published_in" json:"published_in"`
}

func (i ExternalBioMatInput) Save(e edgedb.Executor) (created BioMaterialWithDetails, err error) {
	data, _ := json.Marshal(i)
	code, err := i.GetCode(e)
	if err != nil {
		return
	}
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0,
			select (insert occurrence::ExternalBioMat {
				original_link := <str>json_get(data, 'original_link'),
				original_taxon := <str>json_get(data, 'original_taxon'),
				quantity := <occurrence::QuantityType>json_get(data, 'quantity'),
				content_description := <str>json_get(data, 'content_description'),
				in_collection := <str>json_get(data, 'collection'),
				item_vouchers := <str>json_array_unpack(json_get(data, 'item_vouchers')),
				comments := <str>json_get(data, 'comments'),
				code := <str>$1,
        identification := (
          insert occurrence::Identification {
            taxon := assert_exists(
              select taxonomy::Taxon
              filter .code = <str>data['identification']['taxon']
            ),
            identified_by := assert_exists(
              select people::Person filter .alias = <str>data['identification']['identified_by']
            ),
            identified_on := date::from_json_with_precision(data['identification']['identified_on']),
          }
        ),
        sampling := assert_exists(
          select (<occurrence::Sampling><uuid>data['sampling_id'])
        ),
        is_type := <bool>json_get(data, 'is_type'),
			}) {
        **,
				sequence_consensus: { * },
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
		`, &created, data, code)
	return
}

type ExternalBioMatUpdate struct {
	BioMaterialUpdate  `edgedb:"$inline" json:",inline"`
	OriginalLink       models.OptionalNull[string]             `edgedb:"original_link" json:"original_link,omitempty"`
	OriginalTaxon      models.OptionalNull[string]             `edgedb:"original_taxon" json:"original_taxon,omitempty"`
	Quantity           models.OptionalInput[specimen.Quantity] `edgedb:"quantity" json:"quantity,omitempty"`
	ContentDescription models.OptionalNull[string]             `edgedb:"content_description" json:"content_description,omitempty"`
	Collection         models.OptionalNull[string]             `edgedb:"in_collection" json:"collection,omitempty"`
	Item               models.OptionalInput[[]string]          `edgedb:"item_vouchers" json:"vouchers,omitempty"`
	Comments           models.OptionalNull[string]             `edgedb:"comments" json:"comments,omitempty"`
	PublishedIn        models.OptionalNull[[]string]           `edgedb:"published_in" json:"published_in"`
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
				sequence_consensus: { * },
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
			"code":                "<str>item['code']", // if not explicitly provided, updated code is autogenerated
			"original_link":       "<str>item['original_link']",
			"original_taxon":      "<str>item['original_taxon']",
			"quantity":            "<occurrence::QuantityType>item['quantity']",
			"content_description": "<str>item['content_description']",
			"in_collection":       "<str>item['collection']",
			"item_vouchers":       "<str>json_array_unpack(item['item_vouchers'])",
			"comments":            "<str>item['comments']",
			"is_type":             "<bool>item['is_type']",
			"identification":      u.Identification.Value.UpdateQuery(".identification"),
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, code, data)
	return
}
