package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
)

type TaxonMatchSource = biomedb.TaxonMatchSource
type TaxonMatchType = biomedb.TaxonMatchType

type GBIFImportStatus = biomedb.GBIFImportStatus

type ResolutionStatus = biomedb.ResolutionStatus

type TaxonGBIFStatus = biomedb.TaxonGBIFStatus

type GBIFImportState struct {
	Status    GBIFImportStatus    `json:"status"`
	ClaimedAt Optional[time.Time] `json:"claimed_at,omitempty"`
}

func GBIFImportStateFromDB(res biomedb.GetGBIFImportStatusRow) GBIFImportState {
	return GBIFImportState{
		Status:    res.GBIFStatus,
		ClaimedAt: NewOptionalFromTimestamp(res.GBIFClaimedAt),
	}
}

type TaxonResolution struct {
	ImportHash string                     `json:"import_hash"`
	InputName  string                     `json:"input_name"`
	Source     Optional[TaxonMatchSource] `json:"source,omitempty"`
	GBIFID     Optional[int32]            `json:"gbif_id,omitempty"`
	TaxonID    Optional[uuid.UUID]        `json:"taxon_id,omitempty"`
	StagingID  Optional[uuid.UUID]        `json:"staging_id,omitempty"`
	Status     Optional[ResolutionStatus] `json:"status,omitempty"`
	GBIFStatus Optional[TaxonGBIFStatus]  `json:"gbif_status,omitempty"`
}

func TaxonResolutionFromDB(res biomedb.TaxonResolution) TaxonResolution {
	return TaxonResolution{
		ImportHash: res.ImportHash,
		InputName:  res.InputName,
		Source:     NewOptionalFromPtr(res.Source),
		GBIFID:     NewOptionalFromPtr(res.GBIFID),
		TaxonID:    NewOptionalFromUUID(res.TaxonID),
		StagingID:  NewOptionalFromUUID(res.StagingID),
		Status:     NewOptionalFromPtr(res.Status),
		GBIFStatus: NewOptionalFromPtr(res.GBIFStatus),
	}
}

func TaxonResolutionFromDBSlice(res []biomedb.TaxonResolution) []TaxonResolution {
	result := make([]TaxonResolution, len(res))
	for i, r := range res {
		result[i] = TaxonResolutionFromDB(r)
	}
	return result
}

type TaxonResolutionState struct {
	Resolution []TaxonResolution           `json:"resolution"`
	Candidates map[string][]TaxonCandidate `json:"candidates"`
}

type TaxonCandidate struct {
	Name       string              `json:"name"`
	Authorship Optional[string]    `json:"authorship,omitempty"`
	Rank       TaxonRank           `json:"rank"`
	Status     TaxonStatus         `json:"status"`
	Source     TaxonMatchSource    `json:"source"`
	MatchType  TaxonMatchType      `json:"match_type"`
	Score      Optional[float64]   `json:"score,omitempty"`
	TaxonID    Optional[uuid.UUID] `json:"taxon_id,omitempty"`
	GBIF_ID    Optional[int32]     `json:"gbif_id,omitempty"`
}

func TaxonCandidateFromDB(candidate biomedb.ListAllTaxonCandidatesRow) TaxonCandidate {
	return TaxonCandidate{
		Name:       candidate.TaxonName,
		Authorship: NewOptionalFromPtr(candidate.TaxonAuthorship),
		Rank:       candidate.ResolvedRank,
		Status:     candidate.ResolvedStatus,
		Source:     candidate.Source,
		MatchType:  candidate.MatchType,
		Score:      NewOptionalFromPtr(candidate.Score),
		TaxonID:    NewOptionalFromUUID(candidate.ResolvedTaxonID),
		GBIF_ID:    NewOptionalFromPtr(candidate.ResolvedGBIFID),
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

func (p TaxonStagingParams) ToParams(importHash string) biomedb.InsertTaxonStagingParams {
	return biomedb.InsertTaxonStagingParams{
		ImportHash:      importHash,
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
