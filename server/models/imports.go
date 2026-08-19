package models

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/app_errors"
)

type TaxonMatchSource = biomedb.TaxonMatchSource
type TaxonMatchType = biomedb.TaxonMatchType

type ResolutionStatus = biomedb.ResolutionStatus

type TaxonGBIFStatus = biomedb.TaxonGBIFStatus

type TaxonDefinition struct {
	Name       string              `json:"name"`
	Authorship Optional[string]    `json:"authorship,omitzero"`
	Rank       Optional[TaxonRank] `json:"rank,omitzero"`
}

// type ImportBatch struct {
// 	ImportID    uuid.UUID           `json:"id"`
// 	Label       string              `json:"label"`
// 	Description Optional[string]    `json:"description,omitzero"`
// 	CreatedBy   uuid.UUID           `json:"created_by"`
// 	CreatedAt   time.Time           `json:"created_at"`
// 	CompletedAt Optional[time.Time] `json:"completed_at,omitzero"`
// }

// func ImportBatchFromDB(res biomedb.ImportBatch) ImportBatch {
// 	return ImportBatch{
// 		ImportID:    res.ImportID,
// 		Label:       res.Label,
// 		Description: NewOptionalFromPtr(res.Description),
// 		CreatedBy:   res.CreatedBy,
// 		CreatedAt:   res.CreatedAt,
// 		CompletedAt: NewOptionalFromTimestamp(res.CompletedAt),
// 	}
// }

type MaterializationReadyCheck struct {
	Taxonomy     bool `json:"taxonomy"`
	Methods      bool `json:"methods"`
	Fixatives    bool `json:"fixatives"`
	Bibliography bool `json:"bibliography"`
}

func (c MaterializationReadyCheck) IsReady() bool {
	return c.Taxonomy && c.Methods && c.Bibliography && c.Fixatives
}

func MaterializationReadyCheckFromDB(res biomedb.CheckReadyToMaterializeRow) MaterializationReadyCheck {
	return MaterializationReadyCheck{
		Taxonomy:     res.Taxonomy,
		Methods:      res.Methods,
		Bibliography: res.Bibliography,
		Fixatives:    res.Fixatives,
	}
}

func (c MaterializationReadyCheck) AppError() *app_errors.AppError {
	taxonomyErrorDetail := &app_errors.AppErrorDetail{
		Message: "Taxonomy materialization is not ready. Please ensure that all taxa have been resolved and staged.",
	}
	methodsErrorDetail := &app_errors.AppErrorDetail{
		Message: "Sampling methods materialization is not ready. Please ensure that all sampling methods have been resolved and staged.",
	}
	fixativesErrorDetail := &app_errors.AppErrorDetail{
		Message: "Sampling fixatives materialization is not ready. Please ensure that all sampling fixatives have been resolved and staged.",
	}
	bibliographyErrorDetail := &app_errors.AppErrorDetail{
		Message: "Bibliography materialization is not ready. Please ensure that all publications have been resolved and staged.",
	}

	errorDetails := []*app_errors.AppErrorDetail{}
	if !c.Taxonomy {
		errorDetails = append(errorDetails, taxonomyErrorDetail)
	}
	if !c.Methods {
		errorDetails = append(errorDetails, methodsErrorDetail)
	}
	if !c.Fixatives {
		errorDetails = append(errorDetails, fixativesErrorDetail)
	}
	if !c.Bibliography {
		errorDetails = append(errorDetails, bibliographyErrorDetail)
	}
	return &app_errors.AppError{
		Code:     "materialization_not_ready",
		Category: "imports",
		ErrorModel: huma.ErrorModel{
			Status: http.StatusUnprocessableEntity,
			Title:  "Materialization not ready",
			Detail: "",
			Errors: errorDetails,
		},
	}
}

type TaxonResolution struct {
	ID                 uuid.UUID                  `json:"id"`
	ImportID           uuid.UUID                  `json:"import_id"`
	InputName          string                     `json:"input_name"`
	InputAuthorship    Optional[string]           `json:"input_authorship,omitzero"`
	InputRank          Optional[string]           `json:"input_rank,omitzero"`
	ScientificName     string                     `json:"scientific_name"`
	ResolvedTo         Optional[uuid.UUID]        `json:"resolved_to,omitzero"`
	Status             Optional[ResolutionStatus] `json:"status,omitzero"`
	GBIFStatus         Optional[TaxonGBIFStatus]  `json:"gbif_status,omitzero"`
	FromResolutionID   Optional[uuid.UUID]        `json:"from_resolution_id,omitzero"`
	FromResolutionName Optional[string]           `json:"from_resolution_name,omitzero"`
	SamplingTarget     bool                       `json:"sampling_target"`
}

func TaxonResolutionFromDB(res biomedb.TaxonResolution) TaxonResolution {
	return TaxonResolution{
		ID:               res.ID,
		ImportID:         res.ImportID,
		InputName:        res.InputName,
		InputAuthorship:  NewOptionalFromPtr(res.InputAuthorship),
		InputRank:        NewOptionalFromPtr(res.InputRank),
		ScientificName:   res.ScientificName,
		ResolvedTo:       NewOptionalFromUUID(res.ResolvedCandidateID),
		Status:           NewOptionalFromPtr(res.Status),
		GBIFStatus:       NewOptionalFromPtr(res.GBIFStatus),
		FromResolutionID: NewOptionalFromUUID(res.FromResolutionID),
		SamplingTarget:   res.SamplingTarget,
	}
}

func TaxonResolutionFromDBSlice(res []biomedb.TaxonResolution) []TaxonResolution {
	result := make([]TaxonResolution, len(res))
	for i, r := range res {
		result[i] = TaxonResolutionFromDB(r)
	}
	return result
}

func (t *TaxonResolution) SetFromResolutionName(fromResolutionName *string) {
	t.FromResolutionName = NewOptionalFromPtr(fromResolutionName)
}

type TaxonResolutionWithCandidates struct {
	TaxonResolution
	Candidates []TaxonCandidate `json:"candidates"`
}

type ResolveInput struct {
	ResolutionID uuid.UUID `json:"resolution_id"`
	CandidateID  uuid.UUID `json:"candidate_id"`
}

type TaxonCandidate struct {
	ID           uuid.UUID           `json:"id"`
	ResolutionID uuid.UUID           `json:"resolution_id"`
	Name         string              `json:"name"`
	Rank         TaxonRank           `json:"rank"`
	Authorship   Optional[string]    `json:"authorship,omitzero"`
	Status       TaxonStatus         `json:"status"`
	Source       TaxonMatchSource    `json:"source"`
	MatchType    TaxonMatchType      `json:"match_type"`
	Score        Optional[float64]   `json:"score,omitzero"`
	Priority     int32               `json:"priority"`
	TaxonID      Optional[uuid.UUID] `json:"taxon_id,omitzero"`
	GBIF_ID      Optional[int32]     `json:"gbif_id,omitzero"`
}

func TaxonCandidateFromDB(candidate biomedb.ListAllTaxonCandidatesRow) TaxonCandidate {
	return TaxonCandidate{
		ID:           candidate.ID,
		ResolutionID: candidate.ResolutionID,
		Name:         candidate.ResolvedName,
		Rank:         TaxonRank(candidate.ResolvedRank),
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
	ResolutionID uuid.UUID        `json:"resolution_id"`
	Name         string           `json:"name"`
	Authorship   Optional[string] `json:"authorship,omitzero"`
	Rank         TaxonRank        `json:"rank"`
	Status       TaxonStatus      `json:"status"`
	ParentName   string           `json:"parent_name"`
}

func (p TaxonStagingParams) ToParams(importID uuid.UUID) biomedb.InsertTaxaStagingParams {
	parentRank := string(ParentRank(p.Rank))
	return biomedb.InsertTaxaStagingParams{
		ImportID:     importID,
		ResolutionID: p.ResolutionID,
		Name:         p.Name,
		Authorship:   p.Authorship.ToPtr(),
		TaxonRank:    p.Rank,
		TaxonStatus:  p.Status,
		ParentName:   p.ParentName,
		ParentRank:   &parentRank,
	}
}

type VocabResolutionStatus = biomedb.VocabResolutionStatus
type SamplingMethodResolution struct {
	ImportID         uuid.UUID             `json:"import_id"`
	InputText        string                `json:"input_text"`
	ResolvedMethodID Optional[uuid.UUID]   `json:"resolved_method_id,omitzero"`
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
	ResolvedFixativeID Optional[uuid.UUID]   `json:"resolved_fixative_id,omitzero"`
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
