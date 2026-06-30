package models

import (
	"time"

	. "github.com/go-jet/jet/v2/postgres"
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

type Occurrence struct {
	BaseOccurrence
	Sampling
}

func OccurrenceFromDB(o biomedb.Occurrence, taxon biomedb.Taxon, s biomedb.Sampling, c biomedb.Country) Occurrence {
	return Occurrence{
		BaseOccurrence: BaseOccurrenceFromDB(o, taxon),
		Sampling:       NewSamplingFromDB(s, c),
	}
}

type OccurrenceWithDetails struct {
	BaseOccurrence
	SamplingWithDetails
	OccurrenceMetadata
}

func (o Occurrence) WithDetails(samplingMetadata SamplingMetadata, occurrenceMetadata OccurrenceMetadata) OccurrenceWithDetails {
	return OccurrenceWithDetails{
		BaseOccurrence:      o.BaseOccurrence,
		SamplingWithDetails: o.Sampling.WithDetails(samplingMetadata),
		OccurrenceMetadata:  occurrenceMetadata,
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

type Collection struct {
	Name     string   `json:"name"`
	Vouchers []string `json:"vouchers,omitempty"`
}

func CollectionFromDB(c biomedb.OccurrenceCollection) Collection {
	return Collection{
		Name:     c.Name,
		Vouchers: c.Vouchers,
	}
}
