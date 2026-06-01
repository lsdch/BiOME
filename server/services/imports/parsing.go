package imports

import (
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/db/pgutils"
	"github.com/lsdch/biome/models"
)

type CoordinatesInput struct {
	Latitude  float32 `csv:"latitude"`
	Longitude float32 `csv:"longitude"`
}

type CoordinatesWithPrecisionInput struct {
	CoordinatesInput
	PrecisionM *int32 `csv:"coordinates_precision_m,omitempty"`
}

type EventDateInput struct {
	Date      time.Time
	Precision biomedb.EventDatePrecision
	IsSet     bool
}

func (d *EventDateInput) String() string {
	if !d.IsSet {
		return ""
	}
	switch d.Precision {
	case biomedb.EventDatePrecisionDay:
		return d.Date.Format("2006-01-02")
	case biomedb.EventDatePrecisionMonth:
		return d.Date.Format("2006-01")
	case biomedb.EventDatePrecisionYear:
		return d.Date.Format("2006")
	default:
		return d.Date.Format(time.RFC3339)
	}
}

func (d *EventDateInput) UnmarshalCSV(data []byte) error {
	str := strings.TrimSpace(string(data))
	if str == "" {
		d.IsSet = false
		return nil
	}
	date, err := time.Parse("2006-01-02", str)
	if err == nil {
		d.Date = date
		d.Precision = biomedb.EventDatePrecisionDay
		d.IsSet = true
		return nil
	}
	date, err = time.Parse("2006-01", str)
	if err == nil {
		d.Date = date
		d.Precision = biomedb.EventDatePrecisionMonth
		d.IsSet = true
		return nil
	}
	date, err = time.Parse("2006", str)
	if err == nil {
		d.Date = date
		d.Precision = biomedb.EventDatePrecisionYear
		d.IsSet = true
		return nil
	}
	return fmt.Errorf("invalid date format: %s", str)
}

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

type QuantityInput struct {
	Exact models.OptionalInput[int32]
	Lower models.OptionalInput[int32]
	Upper models.OptionalInput[int32]
}

func (q *QuantityInput) UnmarshalCSV(data []byte) error {
	str := strings.TrimSpace(string(data))
	if str == "" {
		return nil
	}
	parts := strings.Split(str, "-")
	if len(parts) == 1 {
		var exact int32
		_, err := fmt.Sscanf(strings.TrimSpace(parts[0]), "%d", &exact)
		if err != nil {
			return fmt.Errorf("invalid quantity format: %s", str)
		}
		q.Exact = models.NewOptionalInput(exact)
	} else if len(parts) == 2 {
		var lower, upper int32
		_, err := fmt.Sscanf(strings.TrimSpace(parts[0]), "%d", &lower)
		if err != nil {
			return fmt.Errorf("invalid quantity format: %s", str)
		}
		_, err = fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &upper)
		if err != nil {
			return fmt.Errorf("invalid quantity format: %s", str)
		}
		q.Lower = models.NewOptionalInput(lower)
		q.Upper = models.NewOptionalInput(upper)
	} else {
		return fmt.Errorf("invalid quantity format: %s", str)
	}
	return nil
}

type OccurrenceImportRow struct {
	RowNumber int32 `csv:"row_number"`

	// Sampling

	SamplingID        string                        `csv:"sampling_id,omitempty"`
	SiteCode          string                        `csv:"site_code,omitempty"`
	SiteName          string                        `csv:"site_name,omitempty"`
	SiteLocality      string                        `csv:"locality,omitempty"`
	SiteCountryCode   string                        `csv:"country"`
	Coordinates       CoordinatesWithPrecisionInput `csv:",inline"`
	Altitude          *int32                        `csv:"altitude,omitempty"`
	EventDate         EventDateInput                `csv:"event_date,omitempty"`
	PerformedBy       StringListInput               `csv:"performed_by,omitempty"`
	Duration          int32                         `csv:"duration,omitempty"`
	AccessPoints      StringListInput               `csv:"access_points,omitempty"`
	Habitats          StringListInput               `csv:"habitats,omitempty"`
	SamplingTargets   StringListInput               `csv:"sampling_targets,omitempty"`
	SamplingMethods   StringListInput               `csv:"sampling_methods,omitempty"`
	SamplingFixatives StringListInput               `csv:"sampling_fixatives,omitempty"`
	SamplingComments  string                        `csv:"sampling_comments,omitempty"`

	// Occurrence

	OccurrenceCode     string                        `csv:"occurrence_code,omitempty"`
	TypeStatus         *biomedb.OccurrenceTypeStatus `csv:"type_status,omitempty"`
	OccurrenceComments string                        `csv:"comments,omitempty"`

	// Identification

	TaxonName              string          `csv:"taxon_name"`
	TaxonRank              string          `csv:"taxon_rank,omitempty"`
	TaxonAuthorship        string          `csv:"taxon_authorship,omitempty"`
	VerbatimIdentification string          `csv:"verbatim_identification,omitempty"`
	IdentifiedBy           StringListInput `csv:"identified_by,omitempty"`
	IdentificationDate     EventDateInput  `csv:"identification_date,omitempty"`
	IdentificationConfer   bool            `csv:"identification_confer,omitempty"`
	IdentificationAddendum string          `csv:"identification_addendum,omitempty"`

	// Content

	ContentDescription string          `csv:"content_description,omitempty"`
	Quantity           QuantityInput   `csv:"quantity,omitempty"`
	Sources            StringListInput `csv:"sources,omitempty"`
}

func (r *OccurrenceImportRow) SamplingHash() string {

	// If SamplingID is provided, use it directly as the hash
	if r.SamplingID != "" {
		return fmt.Sprintf("id:%s", r.SamplingID)
	}

	datePart := r.EventDate.String()
	// If the date is not set, use the row number as a fallback to ensure uniqueness
	if datePart == "" {
		datePart = string(r.RowNumber)
	}

	coordsPart := fmt.Sprintf("%f|%f|%d", r.Coordinates.Latitude, r.Coordinates.Longitude, r.Coordinates.PrecisionM)

	return strings.Join([]string{
		r.SiteCode,
		r.SiteName,
		r.SiteLocality,
		r.SiteCountryCode,
		coordsPart,
		datePart,
		string(r.Duration),
		strings.Join(r.PerformedBy.Values, ";"),
		strings.Join(r.AccessPoints.Values, ";"),
		strings.Join(r.Habitats.Values, ";"),
		strings.Join(r.SamplingTargets.Values, ";"),
		strings.Join(r.SamplingMethods.Values, ";"),
		strings.Join(r.SamplingFixatives.Values, ";"),
	}, "|")
}

func (r *OccurrenceImportRow) ToStaging(importHash string) biomedb.CopyImportStagingParams {

	return biomedb.CopyImportStagingParams{
		RowNumber:  r.RowNumber,
		ImportHash: importHash,

		// Sampling fields

		SamplingHash:         r.SamplingHash(),
		SiteCode:             pgutils.Text(r.SiteCode),
		SiteName:             pgutils.Text(r.SiteName),
		SiteLocality:         pgutils.Text(r.SiteLocality),
		SiteCountryCode:      pgutils.Text(r.SiteCountryCode),
		Longitude:            r.Coordinates.Longitude,
		Latitude:             r.Coordinates.Latitude,
		CoordinatesPrecision: pgutils.Int4Ptr(r.Coordinates.PrecisionM),
		SamplingComments:     pgutils.Text(r.SamplingComments),
		Altitude:             pgutils.Int4Ptr(r.Altitude),
		EventDate:            pgtype.Date{Time: r.EventDate.Date, Valid: r.EventDate.IsSet},
		EventDatePrecision:   biomedb.NullEventDatePrecision{Valid: r.EventDate.IsSet, EventDatePrecision: r.EventDate.Precision},
		Duration:             pgutils.Int4(r.Duration),
		PerformedBy:          r.PerformedBy.Values,
		AccessPoints:         r.AccessPoints.Values,
		Habitats:             r.Habitats.Values,
		SamplingTargets:      r.SamplingTargets.Values,
		SamplingFixatives:    r.SamplingFixatives.Values,
		SamplingMethods:      r.SamplingMethods.Values,

		// Occurrence fields

		OccurrenceCode: pgutils.Text(r.OccurrenceCode),
		TypeStatus: biomedb.NullOccurrenceTypeStatus{
			Valid:                r.TypeStatus != nil,
			OccurrenceTypeStatus: *r.TypeStatus,
		},
		OccurrenceComments: pgutils.Text(r.OccurrenceComments),

		// Identification fields

		TaxonName:              r.TaxonName,
		TaxonRank:              pgutils.Text(r.TaxonRank),
		TaxonAuthorship:        pgutils.Text(r.TaxonAuthorship),
		VerbatimIdentification: pgutils.Text(r.VerbatimIdentification),
		IdentifiedBy:           r.IdentifiedBy.Values,
		IdentificationDate:     pgtype.Date{Time: r.IdentificationDate.Date, Valid: r.IdentificationDate.IsSet},
		IdentificationDatePrecision: biomedb.NullEventDatePrecision{
			Valid:              r.IdentificationDate.IsSet,
			EventDatePrecision: r.IdentificationDate.Precision,
		},
		IdentificationConfer:   r.IdentificationConfer,
		IdentificationAddendum: pgutils.Text(r.IdentificationAddendum),

		// Content fields

		ContentDescription: pgutils.Text(r.ContentDescription),
		QuantityExact:      pgutils.Int4Opt(r.Quantity.Exact),
		QuantityLower:      pgutils.Int4Opt(r.Quantity.Lower),
		QuantityUpper:      pgutils.Int4Opt(r.Quantity.Upper),
		Sources:            r.Sources.Values,
	}
}

type CSVParseError struct {
	RowNumber int32
	Err       error
}

func (e *CSVParseError) Error() string {
	return fmt.Sprintf("error parsing CSV at row %d: %v", e.RowNumber, e.Err)
}
