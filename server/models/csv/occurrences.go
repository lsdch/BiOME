package csvmodels

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/types"
)

const STRING_LIST_DEFAULT_SEPARATOR = "|"

type StringListInput struct {
	Values []string `csv:"values"`
	Sep    string
}

func (s *StringListInput) SetSeparator(sep string) {
	s.Sep = sep
}

func (s *StringListInput) UnmarshalCSV(data []byte) error {
	str := strings.TrimSpace(string(data))
	if str == "" {
		s.Values = []string{}
		return nil
	}
	sep := s.Sep
	if sep == "" {
		sep = STRING_LIST_DEFAULT_SEPARATOR
	}
	s.Values = strings.Split(str, sep)
	for i := range s.Values {
		s.Values[i] = strings.TrimSpace(s.Values[i])
	}
	return nil
}

type OccurrenceImportRow struct {
	rowNumber int32      `csv:"-"`
	ID        types.ULID `csv:"-"`

	// Sampling

	SamplingID        *string                       `csv:"sampling_id,omitempty"`
	SiteCode          *string                       `csv:"site_code,omitempty"`
	SiteName          *string                       `csv:"site_name,omitempty"`
	SiteLocality      *string                       `csv:"locality,omitempty"`
	SiteCountryCode   *string                       `csv:"country,omitempty" validate:"omitempty,iso3166_1_alpha3"`
	Coordinates       CoordinatesWithPrecisionInput `csv:",inline" validate:"required"`
	Altitude          *int32                        `csv:"altitude,omitempty"`
	EventDate         models.EventDateInput         `csv:"event_date,omitempty"`
	PerformedBy       StringListInput               `csv:"sampling_participants,omitempty"`
	Duration          *int32                        `csv:"sampling_duration,omitempty" validate:"omitempty,gt=0"`
	AccessPoints      StringListInput               `csv:"access_points,omitempty" sep:","`
	Habitat           StringListInput               `csv:"habitat,omitempty"`
	SamplingTargets   StringListInput               `csv:"sampling_targets,omitempty"`
	SamplingMethods   StringListInput               `csv:"sampling_methods,omitempty"`
	SamplingFixatives StringListInput               `csv:"sampling_fixatives,omitempty"`
	SamplingComments  *string                       `csv:"sampling_comments,omitempty"`

	// Occurrence

	OccurrenceCode     *string                      `csv:"occurrence_code,omitempty"`
	TypeStatus         *models.OccurrenceTypeStatus `csv:"type_status,omitempty"`
	OccurrenceComments *string                      `csv:"occurrence_comments,omitempty"`

	// Identification

	TaxonName              string                `csv:"taxon_name" validate:"required"`
	TaxonRank              *models.TaxonRank     `csv:"taxon_rank,omitempty"`
	TaxonAuthorship        *string               `csv:"taxon_authorship,omitempty"`
	VerbatimIdentification *string               `csv:"verbatim_identification,omitempty"`
	IdentifiedBy           StringListInput       `csv:"identified_by,omitempty"`
	IdentificationDate     models.EventDateInput `csv:"identification_date,omitempty"`
	IdentificationConfer   bool                  `csv:"identification_confer,omitempty"`
	IdentificationAddendum *string               `csv:"identification_addendum,omitempty"`

	// Content

	ContentDescription *string              `csv:"content_description,omitempty"`
	Quantity           models.QuantityInput `csv:"specimen_quantity,omitempty"`

	// References

	Collections models.CollectionArrayInput `csv:"collections,omitempty"`
	Sources     StringListInput             `csv:"sources,omitempty"`

	// Publication fields
	Publication PublicationImportRow `csv:"pub_,inline"`
}

func (r *OccurrenceImportRow) WithTaxonDefinition(taxonDefinition models.TaxonDefinition) {
	r.TaxonRank = taxonDefinition.Rank.ToPtr()
	r.TaxonAuthorship = taxonDefinition.Authorship.ToPtr()
}

func (r OccurrenceImportRow) HasPublication() bool {
	return r.Publication.DOI != nil || r.Publication.Verbatim != nil
}

func (r OccurrenceImportRow) Validate(v *validator.Validate) error {
	if r.TaxonRank != nil && *r.TaxonRank == biomedb.TaxonRankSubspecies {
		fields := strings.Fields(r.TaxonName)
		if len(fields) < 3 {
			return &CSVParseError{RowNumber: r.rowNumber, Err: errors.New("subspecies taxon name must have at least 3 fields")}
		}
	}
	return v.Struct(r)
}

func (r OccurrenceImportRow) RowNumber() int32 {
	return r.rowNumber
}

func (r *OccurrenceImportRow) SetRowNumber(rowNumber int32) *OccurrenceImportRow {
	r.rowNumber = rowNumber
	return r
}

func (r *OccurrenceImportRow) SamplingHash() string {

	// If SamplingID is provided, use it directly as the hash
	if r.SamplingID != nil && *r.SamplingID != "" {
		return fmt.Sprintf("id:%s", *r.SamplingID)
	}

	datePart := r.EventDate.String()
	// If the date is not set, use the row number as a fallback to ensure uniqueness
	if datePart == "" {
		datePart = string(r.rowNumber)
	}

	return strings.Join([]string{
		models.ValueOrZero(r.SiteCode),
		models.ValueOrZero(r.SiteName),
		models.ValueOrZero(r.SiteLocality),
		models.ValueOrZero(r.SiteCountryCode),
		r.Coordinates.String(),
		datePart,
		fmt.Sprintf("%d", models.ValueOrZero(r.Duration)),
		strings.Join(r.PerformedBy.Values, "|"),
		strings.Join(r.AccessPoints.Values, "|"),
		strings.Join(r.Habitat.Values, "|"),
		strings.Join(r.SamplingTargets.Values, "|"),
		strings.Join(r.SamplingMethods.Values, "|"),
		strings.Join(r.SamplingFixatives.Values, "|"),
	}, "|")
}

func (r *OccurrenceImportRow) ToStaging(importID uuid.UUID) biomedb.CopyImportStagingParams {

	return biomedb.CopyImportStagingParams{
		RowNumber: r.rowNumber,
		ImportID:  importID,
		ID:        types.MakeULID(),

		// Sampling fields
		SamplingHash:         r.SamplingHash(),
		SiteCode:             r.SiteCode,
		SiteName:             r.SiteName,
		SiteLocality:         r.SiteLocality,
		SiteCountryCode:      r.SiteCountryCode,
		Longitude:            r.Coordinates.Longitude,
		Latitude:             r.Coordinates.Latitude,
		CoordinatesPrecision: r.Coordinates.PrecisionM,
		SamplingComments:     r.SamplingComments,
		Altitude:             r.Altitude,
		EventDate:            r.EventDate.DatePg(),
		EventDatePrecision:   r.EventDate.PrecisionPtr(),
		Duration:             r.Duration,
		PerformedBy:          r.PerformedBy.Values,
		AccessPoints:         r.AccessPoints.Values,
		Habitats:             r.Habitat.Values,
		SamplingTargets:      r.SamplingTargets.Values,
		SamplingFixatives:    r.SamplingFixatives.Values,
		SamplingMethods:      r.SamplingMethods.Values,

		// Occurrence fields

		OccurrenceCode:     r.OccurrenceCode,
		TypeStatus:         (r.TypeStatus),
		OccurrenceComments: r.OccurrenceComments,

		// Identification fields

		TaxonName:                   r.TaxonName,
		TaxonRank:                   (r.TaxonRank),
		TaxonAuthorship:             r.TaxonAuthorship,
		VerbatimIdentification:      r.VerbatimIdentification,
		IdentifiedBy:                r.IdentifiedBy.Values,
		IdentificationDate:          r.IdentificationDate.DatePg(),
		IdentificationDatePrecision: r.IdentificationDate.PrecisionPtr(),
		IdentificationConfer:        r.IdentificationConfer,
		IdentificationAddendum:      r.IdentificationAddendum,

		// Content fields

		ContentDescription: r.ContentDescription,
		QuantityExact:      r.Quantity.Exact.ToPtr(),
		QuantityLower:      r.Quantity.Lower.ToPtr(),
		QuantityUpper:      r.Quantity.Upper.ToPtr(),
		Sources:            r.Sources.Values,
	}
}
