package csvmodels

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/types"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type StringListInput struct {
	Values []string `csv:"values"`
}

func (s *StringListInput) UnmarshalCSV(data []byte) error {
	str := strings.TrimSpace(string(data))
	if str == "" {
		s.Values = []string{}
		return nil
	}
	s.Values = strings.Split(str, ";")
	for i := range s.Values {
		s.Values[i] = strings.TrimSpace(s.Values[i])
	}
	return nil
}

type OccurrenceImportRow struct {
	RowNumber int32 `csv:"row_number"`

	// Sampling

	SamplingID        *string                       `csv:"sampling_id,omitempty"`
	SiteCode          *string                       `csv:"site_code,omitempty"`
	SiteName          *string                       `csv:"site_name,omitempty"`
	SiteLocality      *string                       `csv:"locality,omitempty"`
	SiteCountryCode   string                        `csv:"country" validate:"required,iso3166_1_alpha3"`
	Coordinates       CoordinatesWithPrecisionInput `csv:",inline"`
	Altitude          *int32                        `csv:"altitude,omitempty"`
	EventDate         models.EventDateInput         `csv:"event_date,omitempty"`
	PerformedBy       StringListInput               `csv:"sampling_participants,omitempty"`
	Duration          *int32                        `csv:"sampling_duration,omitempty" validate:"omitempty,gt=0"`
	AccessPoints      StringListInput               `csv:"access_points,omitempty"`
	Habitats          StringListInput               `csv:"habitats,omitempty"`
	SamplingTargets   StringListInput               `csv:"sampling_targets,omitempty"`
	SamplingMethods   StringListInput               `csv:"sampling_methods,omitempty"`
	SamplingFixatives StringListInput               `csv:"sampling_fixatives,omitempty"`
	SamplingComments  *string                       `csv:"sampling_comments,omitempty"`

	// Occurrence

	OccurrenceCode     *string `csv:"occurrence_code,omitempty"`
	TypeStatus         *string `csv:"type_status,omitempty"`
	OccurrenceComments *string `csv:"comments,omitempty"`

	// Identification

	TaxonName              string                `csv:"taxon_name"`
	TaxonRank              *string               `csv:"taxon_rank,omitempty"`
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

	Sources StringListInput `csv:"sources,omitempty"`

	// Publication fields

	PubAuthors  StringListInput `csv:"pub_authors,omitempty"`
	PubYear     *int32          `csv:"pub_year,omitempty"`
	PubTitle    *string         `csv:"pub_title,omitempty"`
	PubJournal  *string         `csv:"pub_journal,omitempty"`
	PubVerbatim *string         `csv:"pub_verbatim,omitempty"`
	PubDOI      *string         `csv:"pub_DOI,omitempty"`
}

func (r *OccurrenceImportRow) Validate() error {
	return validate.Struct(r)
}

func (r *OccurrenceImportRow) SamplingHash() string {

	// If SamplingID is provided, use it directly as the hash
	if r.SamplingID != nil && *r.SamplingID != "" {
		return fmt.Sprintf("id:%s", *r.SamplingID)
	}

	datePart := r.EventDate.String()
	// If the date is not set, use the row number as a fallback to ensure uniqueness
	if datePart == "" {
		datePart = string(r.RowNumber)
	}

	coordsPart := fmt.Sprintf("%f|%f|%d", r.Coordinates.Latitude, r.Coordinates.Longitude, r.Coordinates.PrecisionM)

	return strings.Join([]string{
		models.ValueOrZero(r.SiteCode),
		models.ValueOrZero(r.SiteName),
		models.ValueOrZero(r.SiteLocality),
		r.SiteCountryCode,
		coordsPart,
		datePart,
		fmt.Sprintf("%d", models.ValueOrZero(r.Duration)),
		strings.Join(r.PerformedBy.Values, ";"),
		strings.Join(r.AccessPoints.Values, ";"),
		strings.Join(r.Habitats.Values, ";"),
		strings.Join(r.SamplingTargets.Values, ";"),
		strings.Join(r.SamplingMethods.Values, ";"),
		strings.Join(r.SamplingFixatives.Values, ";"),
	}, "|")
}

func (r *OccurrenceImportRow) ToStaging(importID uuid.UUID) biomedb.CopyImportStagingParams {

	var typeStatus *models.OccurrenceTypeStatus = nil
	if r.TypeStatus != nil {
		typeStatus = new(models.OccurrenceTypeStatus)
		*typeStatus = models.OccurrenceTypeStatus(strings.ToUpper(*r.TypeStatus))
	}

	return biomedb.CopyImportStagingParams{
		RowNumber: r.RowNumber,
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
		Habitats:             r.Habitats.Values,
		SamplingTargets:      r.SamplingTargets.Values,
		SamplingFixatives:    r.SamplingFixatives.Values,
		SamplingMethods:      r.SamplingMethods.Values,

		// Occurrence fields

		OccurrenceCode:     r.OccurrenceCode,
		TypeStatus:         typeStatus,
		OccurrenceComments: r.OccurrenceComments,

		// Identification fields

		TaxonName:                   r.TaxonName,
		TaxonRank:                   r.TaxonRank,
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
