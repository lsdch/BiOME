package models

import (
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/db/biomedb/biomedb/public/table"
	"github.com/lsdch/biome/types"
)

type OccurrenceSortKey string

//generate:enum
const (
	OccurrenceSortKeyCode         OccurrenceSortKey = "code"
	OccurrenceSortKeySiteName     OccurrenceSortKey = "site_name"
	OccurrenceSortKeySiteCode     OccurrenceSortKey = "site_code"
	OccurrenceSortKeyEventDate    OccurrenceSortKey = "event_date"
	OccurrenceSortKeyTaxonName    OccurrenceSortKey = "taxon_name"
	OccurrenceSortKeyIdentifiedOn OccurrenceSortKey = "identified_on"
	OccurrenceSortKeyCreatedAt    OccurrenceSortKey = "created_at"
	OccurrenceSortKeyUpdatedAt    OccurrenceSortKey = "updated_at"
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
	case OccurrenceSortKeyIdentifiedOn:
		return table.Occurrences.IdentificationDate
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
	ID                 types.ULID                     `json:"id"`
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
		TypeStatus: NewOptionalFromPtr((*OccurrenceTypeStatus)(o.TypeStatus)),
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

func OccurrenceFromDB(o biomedb.Occurrence, taxon biomedb.Taxon, s biomedb.SamplingsWithCountry) Occurrence {
	return Occurrence{
		BaseOccurrence: BaseOccurrenceFromDB(o, taxon),
		Sampling:       NewSamplingFromDB(s),
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
		Sampling:               sampling,
		OccurrenceWithMetadata: o,
	}
}

type OccurrenceWithDetails struct {
	Sampling SamplingWithDetails `json:"sampling"`
	OccurrenceWithMetadata
}

func (o Occurrence) WithDetails(samplingMetadata SamplingMetadata, occurrenceMetadata OccurrenceMetadata) OccurrenceWithDetails {
	return OccurrenceWithDetails{
		Sampling:               o.Sampling.WithDetails(samplingMetadata),
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
	References  []Publication         `json:"references,omitempty"`
	Datasets    []Dataset             `json:"datasets,omitempty"`
	Collections []Collection          `json:"collections,omitempty"`
	ImportBatch Optional[ImportBatch] `json:"import_batch,omitempty"`
}

func NewOccurrenceMetadata(
	codeHistory []CodeHistoryEntry,
	references []Publication,
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

// OccurrenceOverviewItem is a representation of the occurrences count for one taxon
type OccurrenceOverviewItem struct {
	ID          uuid.UUID           `json:"id"`
	ParentID    Optional[uuid.UUID] `json:"parent_id,omitempty"`
	Name        string              `json:"name"`
	Authorship  Optional[string]    `json:"authorship,omitempty"`
	ParentName  Optional[string]    `json:"parent_name,omitempty"`
	Occurrences int32               `json:"occurrences"`
	Rank        TaxonRank           `json:"rank"`
}

func OccurrenceOverviewItemFromDB(t biomedb.OccurrencesByTaxaOverviewRow) OccurrenceOverviewItem {
	return OccurrenceOverviewItem{
		ID:          t.ID,
		ParentID:    NewOptionalFromUUID(t.ParentID),
		Name:        t.Name,
		Authorship:  NewOptionalFromPtr(t.Authorship),
		ParentName:  NewOptionalFromPtr(t.ParentName),
		Occurrences: t.OccurrencesCount,
		Rank:        TaxonRank(t.Rank),
	}
}
