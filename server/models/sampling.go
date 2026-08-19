package models

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/uber/h3-go/v4"
)

type ErrUnknownCode struct {
	Entity string
	Code   string
}

func (e ErrUnknownCode) Error() string {
	return fmt.Sprintf("unknown %s code: %s", e.Entity, e.Code)
}

func ErrUnknownSamplingMethodCode(code string) ErrUnknownCode {
	return ErrUnknownCode{Entity: "sampling method", Code: code}
}

func ErrUnknownFixativeCode(code string) ErrUnknownCode {
	return ErrUnknownCode{Entity: "fixative", Code: code}
}

type Site struct {
	Name     Optional[string]  `json:"name,omitzero"`
	Code     Optional[string]  `json:"code,omitzero"`
	Locality Optional[string]  `json:"locality,omitzero"`
	Country  Optional[Country] `json:"country,omitzero"`
}

func OptionalDateWithPrecisionFromDB(date pgtype.Date, precision *biomedb.EventDatePrecision) Optional[DateWithPrecision] {
	d := Optional[DateWithPrecision]{}
	if date.Valid && precision.Valid() {
		d = NewOptional(DateWithPrecision{Date: date.Time, Precision: *precision})
	}
	return d
}

type Sampling struct {
	ID           uuid.UUID                   `json:"id"`
	Site         Site                        `json:"site"`
	Coordinates  CoordinatesWithPrecision    `json:"coordinates"`
	Altitude     Optional[int32]             `json:"altitude,omitzero"`
	PerformedOn  Optional[DateWithPrecision] `json:"performed_on,omitzero"`
	PerformedBy  []string                    `json:"performed_by,omitempty"`
	Duration     Optional[int32]             `json:"duration,omitzero"`
	AccessPoints []string                    `json:"access_points,omitempty"`
	H3Cell       h3.Cell                     `json:"h3_cell"`
	Comments     Optional[string]            `json:"comments,omitzero"`
}

func (s Sampling) Code() string {
	codePart := s.Site.Code.GetWithDefault(s.Coordinates.ToCode())
	datePart := "NA"
	if s.PerformedOn.IsSet {
		datePart = s.PerformedOn.Value.String()
	}
	return fmt.Sprintf("%s|%s", codePart, datePart)
}

func CountryFromSamplingDB(s biomedb.SamplingsWithCountry) Optional[Country] {
	if s.CountryCode != nil {
		return NewOptional(Country{
			Code:         *s.CountryCode,
			Name:         *s.CountryName,
			Continent:    *s.CountryContinent,
			Subcontinent: *s.CountrySubcontinent,
		})
	}
	return Optional[Country]{}
}

func NewSamplingFromDB(s biomedb.SamplingsWithCountry) Sampling {

	return Sampling{
		ID: s.ID,
		Site: Site{
			Name:     NewOptionalFromPtr(s.SiteName),
			Code:     NewOptionalFromPtr(s.SiteCode),
			Locality: NewOptionalFromPtr(s.SiteLocality),
			Country:  CountryFromSamplingDB(s),
		},
		Coordinates: CoordinatesWithPrecision{
			Coordinates: Coordinates{
				Latitude:  s.Latitude,
				Longitude: s.Longitude,
			},
			Precision: NewOptionalFromPtr(s.CoordinatesPrecision),
		},
		Altitude:     NewOptionalFromPtr(s.Altitude),
		PerformedOn:  OptionalDateWithPrecisionFromDB(s.EventDate, s.EventDatePrecision),
		PerformedBy:  s.PerformedBy,
		Duration:     NewOptionalFromPtr(s.Duration),
		AccessPoints: s.AccessPoints,
		H3Cell:       h3.Cell(s.H3Index),
		Comments:     NewOptionalFromPtr(s.Comments),
	}
}

func (s Sampling) WithDetails(metadata SamplingMetadata) SamplingWithDetails {
	return SamplingWithDetails{
		Sampling:         s,
		SamplingMetadata: metadata,
	}
}

func (s Sampling) WithDistance(distanceMeters int32) SamplingWithDistance {
	return SamplingWithDistance{
		Sampling:       s,
		DistanceMeters: distanceMeters,
	}
}

func (s Sampling) WithOccurrences(occurrences []BaseOccurrence) SamplingWithOccurrences {
	return SamplingWithOccurrences{
		Sampling:    s,
		Occurrences: occurrences,
	}
}

type SamplingInput struct {
	Site         SiteInput              `json:"site"`
	PerformedOn  DateWithPrecisionInput `json:"performed_on"`
	PerformedBy  []string               `json:"performed_by,omitempty"`
	Duration     Optional[int32]        `json:"duration,omitzero"`
	AccessPoints []string               `json:"access_points,omitempty"`
	Methods      []string               `json:"methods,omitempty"`
	Fixatives    []string               `json:"fixatives,omitempty"`
	TargetTaxa   []uuid.UUID            `json:"target_taxa,omitempty"`
}

func (i SamplingInput) ToParams() biomedb.CreateSamplingParams {
	var (
		eventDate          pgtype.Date
		eventDatePrecision *EventDatePrecision = nil
	)
	eventDate = i.PerformedOn.Date.ToPgDate()
	eventDatePrecision = &i.PerformedOn.Precision
	return biomedb.CreateSamplingParams{
		SiteCode:             i.Site.Code.ToPtr(),
		SiteName:             i.Site.Name.ToPtr(),
		SiteLocality:         i.Site.Locality.ToPtr(),
		SiteCountryCode:      i.Site.CountryCode.ToPtr(),
		Coordinates:          i.Site.Coordinates.Coordinates,
		CoordinatesPrecision: i.Site.Coordinates.Precision.ToPtr(),
		Altitude:             i.Site.Altitude.ToPtr(),
		EventDate:            eventDate,
		EventDatePrecision:   eventDatePrecision,
		PerformedBy:          i.PerformedBy,
		Duration:             i.Duration.ToPtr(),
		AccessPoints:         i.AccessPoints,
	}
}

type SiteInput struct {
	Name        Optional[string]         `json:"name,omitzero"`
	Code        Optional[string]         `json:"code,omitzero"`
	Locality    Optional[string]         `json:"locality,omitzero"`
	CountryCode Optional[string]         `json:"country_code,omitzero"`
	Coordinates CoordinatesWithPrecision `json:"coordinates"`
	Altitude    Optional[int32]          `json:"altitude,omitzero"`
}

type ListSamplingsAtProximityInput struct {
	Latitude         float64                 `json:"latitude" query:"latitude" required:"true"`
	Longitude        float64                 `json:"longitude" query:"longitude" required:"true"`
	RadiusMeters     int32                   `json:"radius_meters" query:"radius_meters" required:"true"`
	EventDate        Optional[CompositeDate] `json:"event_date,omitzero" query:"event_date"`
	DateIntervalDays Optional[int32]         `json:"date_interval_days,omitzero" query:"date_interval_days"`
	ExcludeIds       []uuid.UUID             `json:"exclude_ids,omitempty" query:"exclude_ids"`
}

func (i ListSamplingsAtProximityInput) ToParams() biomedb.ListSamplingsAtProximityParams {
	var eventDate pgtype.Date
	if i.EventDate.IsSet {
		eventDate = i.EventDate.Value.ToPgDate()
	}
	return biomedb.ListSamplingsAtProximityParams{
		Latitude:           i.Latitude,
		Longitude:          i.Longitude,
		RadiusMeters:       i.RadiusMeters,
		EventDate:          eventDate,
		DateIntervalDays:   i.DateIntervalDays.GetWithDefault(30),
		ExcludeSamplingIds: i.ExcludeIds,
	}
}
func (i ListSamplingsAtProximityInput) ToParamsH3() biomedb.ListSamplingsH3AtProximityParams {
	// var eventDate pgtype.Date
	// if i.EventDate.IsSet {
	// 	eventDate = i.EventDate.Value.ToPgDate()
	// }
	return biomedb.ListSamplingsH3AtProximityParams{
		Latitude:     i.Latitude,
		Longitude:    i.Longitude,
		RadiusMeters: i.RadiusMeters,
		// EventDate:          eventDate,
		// DateIntervalDays:   i.DateIntervalDays.GetWithDefault(30),
		ExcludeSamplingIds: i.ExcludeIds,
	}
}

type SamplingWithDistance struct {
	Sampling
	DistanceMeters int32 `json:"distance_meters"`
}

func (s SamplingWithDistance) WithOccurrences(occurrences []BaseOccurrence) SamplingWithOccurrencesAndDistance {
	return SamplingWithOccurrencesAndDistance{
		SamplingWithOccurrences: SamplingWithOccurrences{
			Sampling:    s.Sampling,
			Occurrences: occurrences,
		},
		DistanceMeters: s.DistanceMeters,
	}
}

type SamplingWithDetails struct {
	Sampling
	SamplingMetadata
}

type SamplingWithOccurrences struct {
	Sampling
	Occurrences []BaseOccurrence `json:"occurrences,omitempty"`
}

func (s SamplingWithOccurrences) WithDistance(distanceMeters int32) SamplingWithOccurrencesAndDistance {
	return SamplingWithOccurrencesAndDistance{
		SamplingWithOccurrences: s,
		DistanceMeters:          distanceMeters,
	}
}

type SamplingWithOccurrencesAndDistance struct {
	SamplingWithOccurrences
	DistanceMeters int32 `json:"distance_meters"`
}
