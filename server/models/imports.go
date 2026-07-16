package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
)

type TaxonMatchSource = biomedb.TaxonMatchSource
type TaxonMatchType = biomedb.TaxonMatchType

type ResolutionStatus = biomedb.ResolutionStatus

type TaxonGBIFStatus = biomedb.TaxonGBIFStatus

type ImportWorkflow struct {
	ImportID    uuid.UUID           `json:"import_id"`
	Label       string              `json:"label"`
	Description Optional[string]    `json:"description,omitempty"`
	CreatedBy   uuid.UUID           `json:"created_by"`
	CreatedAt   time.Time           `json:"created_at"`
	CompletedAt Optional[time.Time] `json:"completed_at,omitempty"`
}

func ImportWorkflowFromDB(res biomedb.ImportWorkflow) ImportWorkflow {
	return ImportWorkflow{
		ImportID:    res.ImportID,
		Label:       res.Label,
		Description: NewOptionalFromPtr(res.Description),
		CreatedBy:   res.CreatedBy,
		CreatedAt:   res.CreatedAt,
		CompletedAt: NewOptionalFromTimestamp(res.CompletedAt),
	}
}

type ImportWorkflowInput struct {
	Label       string           `json:"label" form:"label" required:"true"`
	Description Optional[string] `json:"description,omitempty" form:"description"`
	AssembledBy []string         `json:"assembled_by,omitempty" form:"assembled_by"`
}

func (i ImportWorkflowInput) ToParams(userID uuid.UUID) biomedb.InitImportWorkflowParams {
	return biomedb.InitImportWorkflowParams{
		Label:       i.Label,
		Description: i.Description.ToPtr(),
		AssembledBy: i.AssembledBy,
		CreatedBy:   userID,
	}
}

type MaterializationReadyCheck struct {
	Taxonomy bool `json:"taxonomy"`
	Methods  bool `json:"methods"`
}

func (c MaterializationReadyCheck) IsReady() bool {
	return c.Taxonomy && c.Methods
}

func MaterializationReadyCheckFromDB(res biomedb.CheckReadyToMaterializeRow) MaterializationReadyCheck {
	return MaterializationReadyCheck{
		Taxonomy: res.Taxonomy,
		Methods:  res.Methods,
	}
}

type TaxonResolution struct {
	ID              uuid.UUID                  `json:"id"`
	ImportID        uuid.UUID                  `json:"import_id"`
	InputName       string                     `json:"input_name"`
	InputAuthorship Optional[string]           `json:"input_authorship,omitempty"`
	InputRank       Optional[string]           `json:"input_rank,omitempty"`
	ScientificName  string                     `json:"scientific_name"`
	ResolvedTo      Optional[uuid.UUID]        `json:"resolved_to,omitempty"`
	Status          Optional[ResolutionStatus] `json:"status,omitempty"`
	GBIFStatus      Optional[TaxonGBIFStatus]  `json:"gbif_status,omitempty"`
}

func TaxonResolutionFromDB(res biomedb.TaxonResolution) TaxonResolution {
	return TaxonResolution{
		ID:              res.ID,
		ImportID:        res.ImportID,
		InputName:       res.InputName,
		InputAuthorship: NewOptionalFromPtr(res.InputAuthorship),
		InputRank:       NewOptionalFromPtr(res.InputRank),
		ScientificName:  res.ScientificName,
		ResolvedTo:      NewOptionalFromUUID(res.ResolvedTo),
		Status:          NewOptionalFromPtr(res.Status),
		GBIFStatus:      NewOptionalFromPtr(res.GBIFStatus),
	}
}

func TaxonResolutionFromDBSlice(res []biomedb.TaxonResolution) []TaxonResolution {
	result := make([]TaxonResolution, len(res))
	for i, r := range res {
		result[i] = TaxonResolutionFromDB(r)
	}
	return result
}

type TaxonResolutionWithCandidates struct {
	TaxonResolution
	Candidates []TaxonCandidate `json:"candidates"`
}

type ResolveTaxonInput struct {
	ResolutionID uuid.UUID `json:"resolution_id"`
	CandidateID  uuid.UUID `json:"candidate_id"`
}

type TaxonCandidate struct {
	ID           uuid.UUID           `json:"id"`
	ResolutionID uuid.UUID           `json:"resolution_id"`
	Name         string              `json:"name"`
	Rank         TaxonRank           `json:"rank"`
	Authorship   Optional[string]    `json:"authorship,omitempty"`
	Status       TaxonStatus         `json:"status"`
	Source       TaxonMatchSource    `json:"source"`
	MatchType    TaxonMatchType      `json:"match_type"`
	Score        Optional[float64]   `json:"score,omitempty"`
	Priority     int32               `json:"priority"`
	TaxonID      Optional[uuid.UUID] `json:"taxon_id,omitempty"`
	GBIF_ID      Optional[int32]     `json:"gbif_id,omitempty"`
}

func TaxonCandidateFromDB(candidate biomedb.ListAllTaxonCandidatesRow) TaxonCandidate {
	return TaxonCandidate{
		ID:           candidate.ID,
		ResolutionID: candidate.ResolutionID,
		Name:         candidate.ResolvedName,
		Rank:         candidate.ResolvedRank,
		Authorship:   NewOptionalFromPtr(candidate.ResolvedAuthorship),
		Status:       candidate.ResolvedStatus,
		Source:       candidate.Source,
		MatchType:    candidate.MatchType,
		Score:        NewOptionalFromPtr(candidate.Score),
		Priority:     candidate.Priority,
		TaxonID:      NewOptionalFromUUID(candidate.ResolvedTaxonID),
		GBIF_ID:      NewOptionalFromPtr(candidate.ResolvedGBIFID),
	}
}

type TaxonStagingParams struct {
	Name            string              `json:"name"`
	Authorship      Optional[string]    `json:"authorship,omitempty"`
	Rank            TaxonRank           `json:"rank"`
	Status          TaxonStatus         `json:"status"`
	ParentSource    TaxonMatchSource    `json:"parent_source"`
	ParentID        Optional[uuid.UUID] `json:"parent_taxa_id,omitempty"`
	ParentGbifID    Optional[int32]     `json:"parent_gbif_id,omitempty"`
	ParentInputName Optional[string]    `json:"parent_input_name,omitempty"`
}

func (p TaxonStagingParams) ToParams(importID uuid.UUID) biomedb.InsertTaxonStagingParams {
	return biomedb.InsertTaxonStagingParams{
		ImportID:        importID,
		Name:            p.Name,
		Authorship:      p.Authorship.ToPtr(),
		Rank:            p.Rank,
		Status:          p.Status,
		ParentSource:    p.ParentSource,
		ParentTaxaID:    UUIDOpt(p.ParentID),
		ParentGBIFID:    p.ParentGbifID.ToPtr(),
		ParentInputName: p.ParentInputName.ToPtr(),
	}
}

type VocabResolutionStatus = biomedb.VocabResolutionStatus
type SamplingMethodResolution struct {
	ImportID         uuid.UUID             `json:"import_id"`
	InputText        string                `json:"input_text"`
	ResolvedMethodID Optional[uuid.UUID]   `json:"resolved_method_id,omitempty"`
	Status           VocabResolutionStatus `json:"status"`
}

func SamplingMethodResolutionFromDB(res biomedb.SamplingMethodsResolution) SamplingMethodResolution {
	return SamplingMethodResolution{
		ImportID:         res.ImportID,
		InputText:        res.InputText,
		ResolvedMethodID: NewOptionalFromUUID(res.ResolvedMethodID),
		Status:           res.Status,
	}
}

type SamplingFixativeResolution struct {
	ImportID           uuid.UUID             `json:"import_id"`
	InputText          string                `json:"input_text"`
	ResolvedFixativeID Optional[uuid.UUID]   `json:"resolved_fixative_id,omitempty"`
	Status             VocabResolutionStatus `json:"status"`
}

func SamplingFixativeResolutionFromDB(res biomedb.SamplingFixativesResolution) SamplingFixativeResolution {
	return SamplingFixativeResolution{
		ImportID:           res.ImportID,
		InputText:          res.InputText,
		ResolvedFixativeID: NewOptionalFromUUID(res.ResolvedFixativeID),
		Status:             res.Status,
	}
}
