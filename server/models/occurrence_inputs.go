package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/types"
)

type BaseOccurrenceInput struct {
	TypeStatus         Optional[OccurrenceTypeStatus] `json:"type_status,omitempty"`
	Identification     IdentificationInput            `json:"identification"`
	Quantity           Optional[QuantityInput]        `json:"quantity,omitempty"`
	ContentDescription Optional[string]               `json:"content_description,omitempty"`
	Sources            []string                       `json:"sources,omitempty"`
	Comments           Optional[string]               `json:"comments,omitempty"`
}

func (i BaseOccurrenceInput) ToParams(samplingID uuid.UUID, code string) biomedb.AddOccurrenceToSamplingParams {

	var (
		quantityExact               *int32              = nil
		quantityLower               *int32              = nil
		quantityUpper               *int32              = nil
		identificationDate                              = pgtype.Date{Valid: false}
		identificationDatePrecision *EventDatePrecision = nil
	)
	if q, ok := i.Quantity.Get(); ok {
		quantityExact = q.Exact.ToPtr()
		quantityLower = q.Lower.ToPtr()
		quantityUpper = q.Upper.ToPtr()
	}
	identificationDate = i.Identification.IdentifiedOn.Date.ToPgDate()
	identificationDatePrecision = &i.Identification.IdentifiedOn.Precision

	return biomedb.AddOccurrenceToSamplingParams{
		Code:                        code,
		SamplingID:                  samplingID,
		TaxonID:                     i.Identification.TaxonID,
		TypeStatus:                  i.TypeStatus.ToPtr(),
		IdentifiedBy:                i.Identification.IdentifiedBy,
		IdentificationDate:          identificationDate,
		IdentificationDatePrecision: identificationDatePrecision,
		IdentificationConfer:        i.Identification.Confer,
		IdentificationAddendum:      i.Identification.Addendum.ToPtr(),
		VerbatimIdentification:      i.Identification.Verbatim.ToPtr(),
		ContentDescription:          i.ContentDescription.ToPtr(),
		QuantityExact:               quantityExact,
		QuantityLower:               quantityLower,
		QuantityUpper:               quantityUpper,
		Sources:                     i.Sources,
		Comments:                    i.Comments.ToPtr(),
		ID:                          types.MakeULID(),
	}
}

type OccurrenceInput struct {
	BaseOccurrenceInput
	Collections []CollectionInput `json:"collections,omitempty"`
	PublishedIn []uuid.UUID       `json:"published_in,omitempty"`
	Datasets    []types.ULID      `json:"datasets,omitempty"`
}

type FullOccurrenceInput struct {
	Occurrence OccurrenceInput `json:"occurrence"`
	Sampling   SamplingInput   `json:"sampling"`
}

type CollectionInput struct {
	Name     string   `json:"name"`
	Vouchers []string `json:"vouchers,omitempty"`
}

func (i CollectionInput) ToDBParams(occurrenceID types.ULID) biomedb.AddOccurrenceCollectionParams {
	return biomedb.AddOccurrenceCollectionParams{
		OccurrenceID: occurrenceID,
		Name:         i.Name,
		Vouchers:     i.Vouchers,
	}
}
