package occurrence

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/people"
	"darco/proto/models/references"
	"darco/proto/models/specimen"

	"github.com/edgedb/edgedb-go"
)

type SpecimenVoucher struct {
	Collection edgedb.OptionalStr `edgedb:"in_collection" json:"collection"`
	Item       []string           `edgedb:"item_vouchers" json:"vouchers"`
}

type ExternalBioMatSpecific struct {
	// ID                 edgedb.UUID        `edgedb:"id" json:"id" format:"uuid"`
	OriginalLink       edgedb.OptionalStr `edgedb:"original_link" json:"original_link,omitempty"`
	OriginalTaxon      edgedb.OptionalStr `edgedb:"original_taxon" json:"original_taxon,omitempty"`
	Quantity           specimen.Quantity  `edgedb:"quantity" json:"quantity"`
	ContentDescription edgedb.OptionalStr `edgedb:"content_description" json:"content_description,omitempty"`
	Archive            SpecimenVoucher    `edgedb:"$inline" json:"archive"`
	Comments           edgedb.OptionalStr `edgedb:"comments" json:"comments"`
}

type BioMaterialType string

//generate:enum
const (
	Internal BioMaterialType = "Internal"
	External BioMaterialType = "External"
)

type BioMaterial struct {
	Occurrence `edgedb:"$inline" json:",inline"`
	Code       string                                  `edgedb:"code" json:"code"`
	References []references.Article                    `edgedb:"published_in" json:"reference,omitempty"`
	Type       BioMaterialType                         `edgedb:"type" json:"type"`
	Meta       people.Meta                             `edgedb:"meta" json:"meta"`
	External   models.Optional[ExternalBioMatSpecific] `edgedb:"external" json:"external,omitempty"`
}

type BioMaterialWithSite struct {
	BioMaterial `edgedb:"$inline" json:",inline"`
	Event       EventInner `edgedb:"event" json:"event"`
}

func ListBioMaterials(db edgedb.Executor) ([]BioMaterialWithSite, error) {
	var items = []BioMaterialWithSite{}
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
