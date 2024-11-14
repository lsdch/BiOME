package occurrence

import (
	"darco/proto/models"
	"darco/proto/models/people"
	"darco/proto/models/references"
	"darco/proto/models/specimen"

	"github.com/edgedb/edgedb-go"
)

type SpecimenVoucher struct {
	Collection string `edgedb:"collection" json:"collection"`
	Item       string `edgedb:"item" json:"item"`
}

type OccurrenceReport struct {
	Occurrence         `edgedb:"$inline" json:",inline"`
	ReportedBy         people.OptionalPerson               `edgedb:"reported_by" json:"reported_by"`
	Reference          models.Optional[references.Article] `edgedb:"reference" json:"reference"`
	OriginalLink       edgedb.OptionalStr                  `edgedb:"original_link" json:"original_link"`
	OriginalTaxon      edgedb.OptionalStr                  `edgedb:"original_taxon" json:"original_taxon"`
	Quantity           specimen.Quantity                   `edgedb:"quantity" json:"quantity"`
	ContentDescription edgedb.OptionalStr                  `edgedb:"content_description" json:"content_description"`
	Voucher            models.Optional[SpecimenVoucher]    `edgedb:"voucher" json:"voucher"`
}
