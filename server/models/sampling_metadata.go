package models

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
)

type samplingVocab struct {
	ID uuid.UUID `json:"id"`
	SamplingVocabInput
}

type SamplingVocabInput struct {
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Description Optional[string] `json:"description,omitempty"`
}

type SamplingMethod samplingVocab

func SamplingMethodFromDB(m biomedb.SamplingMethod) SamplingMethod {
	return SamplingMethod{
		ID: m.ID,
		SamplingVocabInput: SamplingVocabInput{
			Code:        m.Code,
			Name:        m.Name,
			Description: NewOptionalFromPtr(m.Description),
		},
	}
}

type SamplingMethodInput SamplingVocabInput

func (i SamplingMethodInput) ToDBParams() biomedb.CreateSamplingMethodParams {
	return biomedb.CreateSamplingMethodParams{
		Code:        i.Code,
		Name:        i.Name,
		Description: i.Description.ToPtr(),
	}
}

type Fixative samplingVocab

func FixativeFromDB(f biomedb.Fixative) Fixative {
	return Fixative{
		ID: f.ID,
		SamplingVocabInput: SamplingVocabInput{
			Code:        f.Code,
			Name:        f.Name,
			Description: NewOptionalFromPtr(f.Description),
		},
	}
}

type FixativeInput SamplingVocabInput

func (i FixativeInput) ToDBParams() biomedb.CreateFixativeParams {
	return biomedb.CreateFixativeParams{
		Code:        i.Code,
		Name:        i.Name,
		Description: i.Description.ToPtr(),
	}
}

type SamplingMetadata struct {
	SamplingMethods []SamplingMethod       `json:"methods,omitempty"`
	Fixatives       []Fixative             `json:"fixatives,omitempty"`
	TargetTaxa      []Taxon                `json:"target_taxa,omitempty"`
	Habitats        []HabitatWithGroupName `json:"habitats,omitempty"`
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
	if i.Status == biomedb.VocabResolutionStatusSelected {
		if !i.ResolvedMethodId.IsSet {
			return WrapErrorPath(fmt.Errorf("resolution status is 'selected' but no method was provided"), "resolved_method_id")
		}
		if i.ResolvedMethodId.Value == uuid.Nil {
			return WrapErrorPath(fmt.Errorf("resolved_method_id cannot be nil uuid when status is 'selected'"), "resolved_method_id")
		}
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

type SamplingFixativeResolutionInput struct {
	InputText          string                        `json:"input_text"`
	ResolvedFixativeID Optional[uuid.UUID]           `json:"resolved_fixative_id,omitempty"`
	Status             biomedb.VocabResolutionStatus `json:"status"`
}

func (i SamplingFixativeResolutionInput) Validate() error {
	if i.Status == biomedb.VocabResolutionStatusSelected && !i.ResolvedFixativeID.IsSet {
		return WrapErrorPath(fmt.Errorf("resolution status is 'selected' but no fixative was provided"), "resolved_fixative_id")
	}
	return nil
}

func (i SamplingFixativeResolutionInput) ToParams(importID uuid.UUID) biomedb.ResolveSamplingFixativeParams {
	return biomedb.ResolveSamplingFixativeParams{
		ImportID:           importID,
		InputText:          i.InputText,
		ResolvedFixativeID: UUIDOpt(i.ResolvedFixativeID),
		Status:             i.Status,
	}
}
