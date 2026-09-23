package models

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/types"
)

type BaseOccurrenceInput struct {
	TypeStatus         Optional[OccurrenceTypeStatus] `json:"type_status,omitzero"`
	Identification     IdentificationInput            `json:"identification"`
	Quantity           Optional[QuantityInput]        `json:"quantity,omitzero"`
	ContentDescription Optional[string]               `json:"content_description,omitzero"`
	Sources            []string                       `json:"sources,omitempty"`
	Comments           Optional[string]               `json:"comments,omitzero"`
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
		TypeStatus:                  (i.TypeStatus.ToPtr()),
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

// pattern to match collection input in the format: "collection_name" or "collection_name[voucher1;voucher2]".
var collectionInputPattern = regexp.MustCompile(
	`^([^\[\]]+)(?:\[([^\[\]]+(?:;[^\[\]]+)*)\])?$`,
)

// UnmarshalCSV implements the csvutil.Unmarshaler interface for CollectionInput.
func (i *CollectionInput) UnmarshalCSV(v []byte) error {
	if len(v) == 0 {
		return nil
	}

	matches := collectionInputPattern.FindStringSubmatch(string(v))
	if matches == nil {
		return fmt.Errorf(
			"invalid collection format %q: expected collection_name or collection_name[voucher1;voucher2]",
			string(v),
		)
	}

	vouchers := []string{}
	if matches[2] != "" {
		vouchers = strings.Split(strings.ReplaceAll(matches[2], " ", ""), ";")
	}

	*i = CollectionInput{
		Name:     matches[1],
		Vouchers: vouchers,
	}

	return nil
}

func (i CollectionInput) ToDBParams(occurrenceID types.ULID) biomedb.AddOccurrenceCollectionParams {
	return biomedb.AddOccurrenceCollectionParams{
		OccurrenceID: occurrenceID,
		Name:         i.Name,
		Vouchers:     i.Vouchers,
	}
}

type CollectionArrayInput []CollectionInput

func (i *CollectionArrayInput) UnmarshalCSV(v []byte) error {
	if len(v) == 0 {
		return nil
	}
	// split by '|' to get individual collections
	collectionsStr := strings.Split(string(v), "|")
	collections := make([]CollectionInput, len(collectionsStr))
	for i, c := range collectionsStr {
		if err := collections[i].UnmarshalCSV([]byte(c)); err != nil {
			return err
		}
	}
	*i = collections
	return nil
}

func (i CollectionArrayInput) CollectionNames() []string {
	names := make([]string, len(i))
	for i, c := range i {
		names[i] = c.Name
	}
	return names
}

func (i CollectionArrayInput) Vouchers() [][]string {
	vouchers := make([][]string, len(i))
	for i, c := range i {
		vouchers[i] = c.Vouchers
	}
	return vouchers
}
