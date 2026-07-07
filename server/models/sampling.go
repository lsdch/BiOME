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
	Name     Optional[string] `json:"name,omitempty"`
	Code     Optional[string] `json:"code,omitempty"`
	Locality Optional[string] `json:"locality,omitempty"`
	Country  Country          `json:"country"`
}

func OptionalDateWithPrecisionFromDB(date pgtype.Date, precision *biomedb.EventDatePrecision) Optional[DateWithPrecision] {
	d := Optional[DateWithPrecision]{}
	if date.Valid && precision.Valid() {
		d = NewOptional(DateWithPrecision{Date: date.Time, Precision: *precision})
	}
	return d
}

type Sampling struct {
	ID           uuid.UUID
	Site         Site                        `json:"site"`
	Coordinates  CoordinatesWithPrecision    `json:"coordinates"`
	Altitude     Optional[int32]             `json:"altitude,omitempty"`
	PerformedOn  Optional[DateWithPrecision] `json:"performed_on,omitempty"`
	PerformedBy  []string                    `json:"performed_by,omitempty"`
	Duration     Optional[int32]             `json:"duration,omitempty"`
	AccessPoints []string                    `json:"access_points,omitempty"`
	H3Cell       h3.Cell                     `json:"h3_cell"`
}

func (s Sampling) Code() string {
	codePart := s.Site.Code.GetWithDefault(s.Coordinates.ToCode())
	datePart := "NA"
	if s.PerformedOn.IsSet {
		datePart = s.PerformedOn.Value.ToCode()
	}
	return fmt.Sprintf("%s|%s", codePart, datePart)
}

func NewSamplingFromDB(s biomedb.Sampling, c biomedb.Country) Sampling {
	return Sampling{
		ID: s.ID,
		Site: Site{
			Name:     NewOptionalFromPtr(s.SiteName),
			Code:     NewOptionalFromPtr(s.SiteCode),
			Locality: NewOptionalFromPtr(s.SiteLocality),
			Country:  Country(c),
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

type SamplingInput struct {
	Site         SiteInput              `json:"site"`
	PerformedOn  DateWithPrecisionInput `json:"performed_on"`
	PerformedBy  []string               `json:"performed_by,omitempty"`
	Duration     Optional[int32]        `json:"duration,omitempty"`
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
	Name        Optional[string]         `json:"name,omitempty"`
	Code        Optional[string]         `json:"code,omitempty"`
	Locality    Optional[string]         `json:"locality,omitempty"`
	CountryCode Optional[string]         `json:"country_code,omitempty"`
	Coordinates CoordinatesWithPrecision `json:"coordinates"`
	Altitude    Optional[int32]          `json:"altitude,omitempty"`
}

type ListSamplingsAtProximityInput struct {
	Latitude         float32                 `json:"latitude" query:"latitude" required:"true"`
	Longitude        float32                 `json:"longitude" query:"longitude" required:"true"`
	RadiusMeters     int32                   `json:"radius_meters" query:"radius_meters" required:"true"`
	EventDate        Optional[CompositeDate] `json:"event_date,omitempty" query:"event_date"`
	DateIntervalDays Optional[int32]         `json:"date_interval_days,omitempty" query:"date_interval_days"`
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

type SamplingWithDistance struct {
	Sampling
	DistanceMeters int32 `json:"distance_meters"`
}

type SamplingWithDetails struct {
	Sampling
	SamplingMetadata
}

type SamplingWithOccurrences struct {
	Sampling
	Occurrences []BaseOccurrence `json:"occurrences,omitempty"`
}

type SamplingVocabUpdateParams struct {
	Name        Optional[string]     `json:"name,omitempty"`
	Code        Optional[string]     `json:"code,omitempty"`
	Description OptionalNull[string] `json:"description,omitempty"`
}

type SamplingMethodUpdateParams SamplingVocabUpdateParams

func (s *SamplingMethodUpdateParams) ToParams(oldCode string) biomedb.UpdateSamplingMethodParams {
	return biomedb.UpdateSamplingMethodParams{
		Code:           s.Code.ToPtr(),
		Name:           s.Name.ToPtr(),
		SetDescription: s.Description.IsSet,
		Description:    s.Description.ToPtr(),
		OldCode:        oldCode,
	}
}

type FixativeUpdateParams SamplingVocabUpdateParams

func (s *FixativeUpdateParams) ToParams(oldCode string) biomedb.UpdateFixativeParams {
	return biomedb.UpdateFixativeParams{
		Code:           s.Code.ToPtr(),
		Name:           s.Name.ToPtr(),
		SetDescription: s.Description.IsSet,
		Description:    s.Description.ToPtr(),
		OldCode:        oldCode,
	}
}

type SamplingMethodResolutionInput struct {
	InputText        string                        `json:"input_text"`
	ResolvedMethodId Optional[uuid.UUID]           `json:"resolved_method_id,omitempty"`
	Status           biomedb.VocabResolutionStatus `json:"status"`
}

func (i SamplingMethodResolutionInput) Validate() error {
	if i.Status == biomedb.VocabResolutionStatusSelected && !i.ResolvedMethodId.IsSet {
		return WrapErrorPath(fmt.Errorf("resolution status is 'selected' but no method was provided"), "resolved_method_id")
	}
	return nil
}

func (i SamplingMethodResolutionInput) ToParams(importID uuid.UUID) biomedb.ResolveMethodParams {
	return biomedb.ResolveMethodParams{
		ImportID:         importID,
		InputText:        i.InputText,
		ResolvedMethodID: UUIDOpt(i.ResolvedMethodId),
		Status:           i.Status,
	}
}
