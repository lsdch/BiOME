package models

import (
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/db/biomedb/biomedb/public/table"
	"github.com/oklog/ulid/v2"
)

type OccurrenceSortKey string

const (
	OccurrenceSortKeyCode      OccurrenceSortKey = "code"
	OccurrenceSortKeySiteName  OccurrenceSortKey = "site_name"
	OccurrenceSortKeySiteCode  OccurrenceSortKey = "site_code"
	OccurrenceSortKeyEventDate OccurrenceSortKey = "event_date"
	OccurrenceSortKeyTaxonName OccurrenceSortKey = "taxon_name"
	OccurrenceSortKeyCreatedAt OccurrenceSortKey = "created_at"
	OccurrenceSortKeyUpdatedAt OccurrenceSortKey = "updated_at"
)

func (k OccurrenceSortKey) Column() Column {
	switch k {
	case OccurrenceSortKeyCode:
		return table.Occurrences.Code
	case OccurrenceSortKeySiteName:
		return table.Samplings.SiteName
	case OccurrenceSortKeySiteCode:
		return table.Samplings.SiteCode
	case OccurrenceSortKeyEventDate:
		return table.Samplings.EventDate
	case OccurrenceSortKeyTaxonName:
		return table.Taxa.Name
	case OccurrenceSortKeyCreatedAt:
		return table.Occurrences.CreatedAt
	case OccurrenceSortKeyUpdatedAt:
		return table.Occurrences.UpdatedAt
	default:
		return table.Occurrences.Code
	}
}

type OccurrencePaginationParams struct {
	Pagination                `json:"pagination"`
	SortBy[OccurrenceSortKey] `json:"sort_by"`
}

type OccurrenceTypeStatus = biomedb.OccurrenceTypeStatus

type Identification struct {
	IdentifiedBy []string                    `json:"identified_by,omitempty"`
	IdentifiedOn Optional[DateWithPrecision] `json:"identified_on,omitempty"`
	Confer       bool                        `json:"confer"`
	Addendum     Optional[string]            `json:"addendum,omitempty"`
	Taxon        Taxon                       `json:"taxon"`
	Verbatim     Optional[string]            `json:"verbatim,omitempty"`
}

type IdentificationInput struct {
	IdentifiedBy []string               `json:"identified_by,omitempty"`
	IdentifiedOn DateWithPrecisionInput `json:"identified_on"`
	Confer       bool                   `json:"confer"`
	Addendum     Optional[string]       `json:"addendum,omitempty"`
	TaxonID      uuid.UUID              `json:"taxon_id"`
	Verbatim     Optional[string]       `json:"verbatim,omitempty"`
}

type BaseOccurrence struct {
	ID                 ulid.ULID                      `json:"id"`
	Code               string                         `json:"code"`
	TypeStatus         Optional[OccurrenceTypeStatus] `json:"type_status,omitempty"`
	Identification     Identification                 `json:"identification"`
	Quantity           Optional[OccurrenceQuantity]   `json:"quantity,omitempty"`
	ContentDescription Optional[string]               `json:"content_description,omitempty"`
	Sources            []string                       `json:"sources,omitempty"`
	Comments           Optional[string]               `json:"comments,omitempty"`
}

func BaseOccurrenceFromDB(o biomedb.Occurrence, taxon biomedb.Taxon) BaseOccurrence {
	identificationDate := OptionalDateWithPrecisionFromDB(o.IdentificationDate, o.IdentificationDatePrecision)

	return BaseOccurrence{
		ID:         o.ID,
		Code:       o.Code,
		TypeStatus: NewOptionalFromPtr(o.TypeStatus),
		Identification: Identification{
			Taxon:        *TaxonFromDB(&taxon),
			IdentifiedBy: o.IdentifiedBy,
			IdentifiedOn: identificationDate,
			Confer:       o.IdentificationConfer,
			Addendum:     NewOptionalFromPtr(o.IdentificationAddendum),
			Verbatim:     NewOptionalFromPtr(o.VerbatimIdentification),
		},
		Quantity:           NewOptionalOccurrenceQuantity(o.QuantityExact, o.QuantityLower, o.QuantityUpper),
		ContentDescription: NewOptionalFromPtr(o.ContentDescription),
		Sources:            o.Sources,
		Comments:           NewOptionalFromPtr(o.Comments),
	}
}

type BaseOccurrenceWithSamplingID struct {
	BaseOccurrence
	SamplingID uuid.UUID `json:"sampling_id"`
}

func (o BaseOccurrence) WithSamplingID(samplingID uuid.UUID) BaseOccurrenceWithSamplingID {
	return BaseOccurrenceWithSamplingID{
		BaseOccurrence: o,
		SamplingID:     samplingID,
	}
}

type Occurrence struct {
	BaseOccurrence
	Sampling Sampling `json:"sampling"`
}

func OccurrenceFromDB(o biomedb.Occurrence, taxon biomedb.Taxon, s biomedb.Sampling, c biomedb.Country) Occurrence {
	return Occurrence{
		BaseOccurrence: BaseOccurrenceFromDB(o, taxon),
		Sampling:       NewSamplingFromDB(s, c),
	}
}

type OccurrenceWithMetadata struct {
	BaseOccurrence
	OccurrenceMetadata
}

func (o BaseOccurrence) WithMetadata(metadata OccurrenceMetadata) OccurrenceWithMetadata {
	return OccurrenceWithMetadata{
		BaseOccurrence:     o,
		OccurrenceMetadata: metadata,
	}
}

func (o OccurrenceWithMetadata) WithSampling(sampling SamplingWithDetails) OccurrenceWithDetails {
	return OccurrenceWithDetails{
		SamplingWithDetails:    sampling,
		OccurrenceWithMetadata: o,
	}
}

type OccurrenceWithDetails struct {
	SamplingWithDetails
	OccurrenceWithMetadata
}

func (o Occurrence) WithDetails(samplingMetadata SamplingMetadata, occurrenceMetadata OccurrenceMetadata) OccurrenceWithDetails {
	return OccurrenceWithDetails{
		SamplingWithDetails:    o.Sampling.WithDetails(samplingMetadata),
		OccurrenceWithMetadata: o.WithMetadata(occurrenceMetadata),
	}
}

type CodeHistoryEntry struct {
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
}

func CodeHistoryEntryFromDB(e biomedb.OccurrenceCodeHistory) CodeHistoryEntry {
	return CodeHistoryEntry{
		Code:      e.Code,
		CreatedAt: e.CreatedAt,
	}
}

type OccurrenceMetadata struct {
	CodeHistory []CodeHistoryEntry    `json:"code_history,omitempty"`
	References  []Article             `json:"references,omitempty"`
	Datasets    []Dataset             `json:"datasets,omitempty"`
	Collections []Collection          `json:"collections,omitempty"`
	ImportBatch Optional[ImportBatch] `json:"import_batch,omitempty"`
}

func NewOccurrenceMetadata(
	codeHistory []CodeHistoryEntry,
	references []Article,
	datasets []Dataset,
	collections []Collection,
	importBatch Optional[ImportBatch],
) OccurrenceMetadata {
	return OccurrenceMetadata{
		CodeHistory: codeHistory,
		References:  references,
		Datasets:    datasets,
		Collections: collections,
		ImportBatch: importBatch,
	}
}

type CollectionInput struct {
	Name     string   `json:"name"`
	Vouchers []string `json:"vouchers,omitempty"`
}

func (i CollectionInput) ToDBParams(occurrenceID ulid.ULID) biomedb.AddOccurrenceCollectionParams {
	return biomedb.AddOccurrenceCollectionParams{
		OccurrenceID: occurrenceID,
		Name:         i.Name,
		Vouchers:     i.Vouchers,
	}
}

type Collection struct {
	ID uuid.UUID `json:"id"`
	CollectionInput
}

func CollectionFromDB(c biomedb.OccurrenceCollection) Collection {
	return Collection{
		ID: c.CollectionID,
		CollectionInput: CollectionInput{
			Name:     c.Name,
			Vouchers: c.Vouchers,
		},
	}
}

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
		ID:                          ulid.Make(),
	}
}

type OccurrenceInput struct {
	BaseOccurrenceInput
	Collections []CollectionInput `json:"collections,omitempty"`
	PublishedIn []uuid.UUID       `json:"published_in,omitempty"`
	Datasets    []ulid.ULID       `json:"datasets,omitempty"`
}

type FullOccurrenceInput struct {
	Occurrence OccurrenceInput `json:"occurrence"`
	Sampling   SamplingInput   `json:"sampling"`
}

// OccurrenceOverviewItem is a representation of the occurrences count for one taxon
type OccurrenceOverviewItem struct {
	Name        string           `json:"name"`
	Authorship  Optional[string] `json:"authorship,omitempty"`
	ParentName  Optional[string] `json:"parent_name,omitempty"`
	Occurrences int32            `json:"occurrences"`
	Rank        TaxonRank        `json:"rank"`
}
