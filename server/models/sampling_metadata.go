package models

import (
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
