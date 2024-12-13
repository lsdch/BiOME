package occurrence

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/people"
	"darco/proto/models/references"
	"darco/proto/models/sequences"
	"darco/proto/models/specimen"
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
	Sequence       `edgedb:"$inline" json:",inline"`
	Origin         sequences.ExtSeqOrigin `edgedb:"origin" json:"origin"`
	Category       sequences.SeqDB        `edgedb:"type" json:"category"`
	Identification Identification         `edgedb:"identification" json:"identification"`
	Comments       edgedb.OptionalStr     `edgedb:"comments" json:"comments"`
	References     []references.Article   `edgedb:"references" json:"references"`
	// SourceSample            `edgedb:"source_sample" json:"source_sample"`
	AccessionNumber    edgedb.OptionalStr `edgedb:"accession_number" json:"accession_number"`
	SpecimenIdentifier string             `edgedb:"specimen_identifier" json:"specimen_identifier"`
	OriginalTaxon      edgedb.OptionalStr `edgedb:"original_taxon" json:"original_taxon"`
}

type ExternalBioMatContent struct {
	Specimen  string                   `edgedb:"specimen" json:"specimen"`
	Sequences []ExternalBioMatSequence `edgedb:"sequences" json:"sequences"`
}

type BioMaterialCategory string

//generate:enum
const (
	Internal BioMaterialCategory = "Internal"
	External BioMaterialCategory = "External"
)

type CodeHistory struct {
	Code string    `edgedb:"code" json:"code"`
	Time time.Time `edgedb:"time" json:"time"`
}

type GenericBioMaterial[SamplingType any] struct {
	GenericOccurrence[SamplingType] `edgedb:"$inline" json:",inline"`
	Code                            string                                  `edgedb:"code" json:"code"`
	CodeHistory                     []CodeHistory                           `edgedb:"code_history" json:"code_history,omitempty"`
	Category                        BioMaterialCategory                     `edgedb:"category" json:"category"`
	IsType                          bool                                    `edgedb:"is_type" json:"is_type"`
	IsHomogenous                    bool                                    `edgedb:"is_homogenous" json:"is_homogenous"`
	IsCongruent                     bool                                    `edgedb:"is_congruent" json:"is_congruent"`
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
						sequences := .elements { * }
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
	Code            string                     `edgedb:"code" json:"code"`
	Category        BioMaterialCategory        `edgedb:"category" json:"category"`
	IsType          models.OptionalInput[bool] `edgedb:"is_type" json:"is_type" default:"false"`
	// Category dependent field
}

type InternalBioMatInput struct {
	BioMaterialInput `edgedb:"$inline" json:",inline"`
	// TODO: Internal-specific fields
}

type ExternalBioMatInput struct {
	BioMaterialInput   `edgedb:"$inline" json:",inline"`
	OriginalLink       models.OptionalInput[string] `edgedb:"original_link" json:"original_link,omitempty"`
	OriginalTaxon      models.OptionalInput[string] `edgedb:"original_taxon" json:"original_taxon,omitempty"`
	Quantity           specimen.Quantity            `edgedb:"quantity" json:"quantity"`
	ContentDescription models.OptionalInput[string] `edgedb:"content_description" json:"content_description,omitempty"`
	Collection         models.OptionalInput[string] `edgedb:"in_collection" json:"collection,omitempty"`
	Item               []string                     `edgedb:"item_vouchers" json:"vouchers,omitempty"`
	Comments           models.OptionalInput[string] `edgedb:"comments" json:"comments,omitempty"`
}
