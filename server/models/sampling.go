package models

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/uber/h3-go/v4"
)

type Site struct {
	Name        Optional[string]         `json:"name,omitempty"`
	Code        Optional[string]         `json:"code,omitempty"`
	Locality    Optional[string]         `json:"locality,omitempty"`
	Country     Country                  `json:"country"`
	Coordinates CoordinatesWithPrecision `json:"coordinates"`
	Altitude    Optional[int32]          `json:"altitude,omitempty"`
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
	PerformedOn  Optional[DateWithPrecision] `json:"performed_on,omitempty"`
	PerformedBy  []string                    `json:"performed_by,omitempty"`
	Duration     Optional[int32]             `json:"duration,omitempty"`
	AccessPoints []string                    `json:"access_points,omitempty"`
	H3Cell       h3.Cell                     `json:"h3_cell"`
}

func NewSamplingFromDB(s biomedb.Sampling, c biomedb.Country) Sampling {
	return Sampling{
		ID: s.ID,
		Site: Site{
			Name:     NewOptionalFromPtr(s.SiteName),
			Code:     NewOptionalFromPtr(s.SiteCode),
			Locality: NewOptionalFromPtr(s.SiteLocality),
			Country:  Country(c),
			Coordinates: CoordinatesWithPrecision{
				Coordinates: Coordinates{
					Latitude:  s.Latitude,
					Longitude: s.Longitude,
				},
				Precision: NewOptionalFromPtr(s.CoordinatesPrecision),
			},
			Altitude: NewOptionalFromPtr(s.Altitude),
		},
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

type samplingVocab struct {
	ID          uuid.UUID        `json:"id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Description Optional[string] `json:"description,omitempty"`
}

type SamplingMethod = samplingVocab

func SamplingMethodFromDB(m biomedb.SamplingMethod) SamplingMethod {
	return SamplingMethod{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: NewOptionalFromPtr(m.Description),
	}
}

type Fixative = samplingVocab

func FixativeFromDB(f biomedb.Fixative) Fixative {
	return Fixative{
		ID:          f.ID,
		Code:        f.Code,
		Name:        f.Name,
		Description: NewOptionalFromPtr(f.Description),
	}
}

type SamplingMetadata struct {
	SamplingMethods []SamplingMethod       `json:"methods,omitempty"`
	Fixatives       []Fixative             `json:"fixatives,omitempty"`
	TargetTaxa      []Taxon                `json:"target_taxa,omitempty"`
	Habitats        []HabitatWithGroupName `json:"habitats,omitempty"`
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
	Name        Optional[string]     `json:"name"`
	Code        Optional[string]     `json:"code"`
	Description OptionalNull[string] `json:"description"`
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

type SamplingFixativeUpdateParams SamplingVocabUpdateParams

func (s *SamplingFixativeUpdateParams) ToParams(oldCode string) biomedb.UpdateFixativeParams {
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

func (i SamplingMethodResolutionInput) ToParams(importHash string) biomedb.ResolveMethodParams {
	return biomedb.ResolveMethodParams{
		ImportHash:       importHash,
		InputText:        i.InputText,
		ResolvedMethodID: UUIDOpt(i.ResolvedMethodId),
		Status:           i.Status,
	}
}
