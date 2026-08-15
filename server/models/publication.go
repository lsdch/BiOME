package models

import (
	"strings"

	"github.com/caltechlibrary/crossrefapi"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/services/crossref"
	"github.com/lsdch/biome/types"
	"github.com/sirupsen/logrus"
)

type Publication struct {
	ID       uuid.UUID           `json:"id"`
	Authors  []string            `json:"authors,omitempty"`
	Year     Optional[int32]     `json:"year,omitempty"`
	Title    Optional[string]    `json:"title,omitempty"`
	Journal  Optional[string]    `json:"journal,omitempty"`
	Verbatim string              `json:"verbatim,omitempty"`
	DOI      Optional[types.DOI] `json:"doi,omitempty"`
	Comments Optional[string]    `json:"comments,omitempty"`
}

func PublicationFromDB(a biomedb.Publication) Publication {
	return Publication{
		ID:       a.ID,
		Authors:  a.Authors,
		Year:     NewOptionalFromPtr(a.Year),
		Title:    NewOptionalFromPtr(a.Title),
		Journal:  NewOptionalFromPtr(a.Journal),
		Verbatim: a.Verbatim,
		DOI:      NewOptionalFromPtr(a.DOI),
		Comments: NewOptionalFromPtr(a.Comments),
	}
}

type CreatePublicationParams struct {
	Authors  []string            `json:"authors,omitempty"`
	Year     Optional[int32]     `json:"year,omitempty"`
	Title    Optional[string]    `json:"title,omitempty"`
	Journal  Optional[string]    `json:"journal,omitempty"`
	Verbatim string              `json:"verbatim,omitempty"`
	DOI      Optional[types.DOI] `json:"doi,omitempty"`
}

func PublicationParamsFromCrossref(cr *crossrefapi.Works) CreatePublicationParams {
	logrus.Debugf("Creating PublicationInput from Crossref data: %+v", cr)
	authors := []string{}
	for _, author := range cr.Message.Author {
		nameParts := []string{}
		if author.Family != "" {
			nameParts = append(nameParts, author.Family)
		}
		if author.Given != "" {
			nameParts = append(nameParts, author.Given)
		}
		authors = append(authors, strings.Join(nameParts, ", "))
	}
	year := int32(0)
	if cr.Message.PublishedPrint != nil && len(cr.Message.PublishedPrint.DateParts) > 0 && len(cr.Message.PublishedPrint.DateParts[0]) > 0 {
		year = int32(cr.Message.PublishedPrint.DateParts[0][0])
	} else if cr.Message.PublishedOnline != nil && len(cr.Message.PublishedOnline.DateParts) > 0 && len(cr.Message.PublishedOnline.DateParts[0]) > 0 {
		year = int32(cr.Message.PublishedOnline.DateParts[0][0])
	}
	return CreatePublicationParams{
		Authors: authors,
		Year:    NewOptional(year),
		Title:   NewOptional(cr.Message.Title[0]),
		Journal: NewOptional(cr.Message.ContainerTitle[0]),
		DOI:     NewOptional(types.DOI(cr.Message.DOI)),
	}
}

func (p CreatePublicationParams) ToDBParams() biomedb.CreatePublicationParams {
	return biomedb.CreatePublicationParams{
		Authors:  p.Authors,
		Year:     p.Year.ToPtr(),
		Title:    p.Title.ToPtr(),
		Journal:  p.Journal.ToPtr(),
		Verbatim: p.Verbatim,
		DOI:      p.DOI.ToPtr(),
	}
}

type UpdatePublicationParams struct {
	Code     Optional[string]        `json:"code,omitempty"`
	Authors  Optional[[]string]      `json:"authors,omitempty"`
	Year     Optional[int32]         `json:"year,omitempty" minimum:"1500"`
	Title    OptionalNull[string]    `json:"title,omitempty"`
	Journal  OptionalNull[string]    `json:"journal,omitempty"`
	Verbatim OptionalNull[string]    `json:"verbatim,omitempty"`
	Comments OptionalNull[string]    `json:"comments,omitempty"`
	DOI      OptionalNull[types.DOI] `json:"doi,omitempty"`
}

func (p UpdatePublicationParams) ToDBParams(id uuid.UUID) biomedb.UpdatePublicationByIDParams {
	return biomedb.UpdatePublicationByIDParams{
		ID:          id,
		Authors:     p.Authors.Value,
		Year:        p.Year.ToPtr(),
		SetTitle:    p.Title.IsSet,
		Title:       p.Title.ToPtr(),
		SetJournal:  p.Journal.IsSet,
		Journal:     p.Journal.ToPtr(),
		SetVerbatim: p.Verbatim.IsSet,
		Verbatim:    p.Verbatim.Value,
		SetComments: p.Comments.IsSet,
		Comments:    p.Comments.ToPtr(),
		SetDOI:      p.DOI.IsSet,
		DOI:         p.DOI.ToPtr(),
	}
}

type BasePublicationCandidate struct {
	ID           uuid.UUID                          `json:"id"`
	ImportID     uuid.UUID                          `json:"import_id"`
	ResolutionID uuid.UUID                          `json:"resolution_id"`
	Source       biomedb.PublicationCandidateSource `json:"source"`
	MatchType    biomedb.PubMatchType               `json:"match_type"`
	InternalID   Optional[uuid.UUID]                `json:"internal_id,omitempty"`
	StagingID    Optional[uuid.UUID]                `json:"staging_id,omitempty"`
	Score        float32                            `json:"score"`
}

func BasePublicationCandidateFromDB(c biomedb.PublicationCandidate) BasePublicationCandidate {
	return BasePublicationCandidate{
		ID:           c.ID,
		ImportID:     c.ImportID,
		ResolutionID: c.ResolutionID,
		Source:       c.Source,
		MatchType:    c.MatchType,
		InternalID:   NewOptionalFromUUID(c.InternalID),
		StagingID:    NewOptionalFromUUID(c.StagingID),
		Score:        c.Score,
	}
}

type PublicationCandidate struct {
	BasePublicationCandidate
	DOI      Optional[types.DOI] `json:"doi,omitempty"`
	Authors  []string            `json:"authors,omitempty"`
	Year     Optional[int32]     `json:"year,omitempty"`
	Title    Optional[string]    `json:"title,omitempty"`
	Journal  Optional[string]    `json:"journal,omitempty"`
	Verbatim string              `json:"verbatim"`
}

type PublicationResolution struct {
	ID         uuid.UUID           `json:"id"`
	ImportID   uuid.UUID           `json:"import_id"`
	Status     ResolutionStatus    `json:"status"`
	ResolvedID Optional[uuid.UUID] `json:"resolved_id,omitempty"`
	DOI        Optional[types.DOI] `json:"doi,omitempty"`
	Verbatim   Optional[string]    `json:"verbatim,omitempty"`
}

func PublicationResolutionFromDB(r biomedb.PublicationResolution) PublicationResolution {
	return PublicationResolution{
		ID:         r.ID,
		ImportID:   r.ImportID,
		Status:     ResolutionStatus(r.Status),
		ResolvedID: NewOptionalFromUUID(r.ResolvedCandidateID),
		DOI:        NewOptionalFromPtr(r.DOI),
		Verbatim:   NewOptionalFromPtr(r.Verbatim),
	}
}

type PublicationResolutionWithCandidates struct {
	PublicationResolution
	Candidates []PublicationCandidate `json:"candidates"`
}

type PublicationStagingInput struct {
	Authors  []string                  `json:"authors,omitempty"`
	Year     Optional[int32]           `json:"year,omitempty"`
	Title    Optional[string]          `json:"title,omitempty"`
	Journal  Optional[string]          `json:"journal,omitempty"`
	Verbatim string                    `json:"verbatim"`
	DOI      Optional[types.DOI]       `json:"doi,omitempty"`
	Source   biomedb.PublicationSource `json:"source"`
}

func PubStagingFromCrossref(cr crossref.WorkMessage) PublicationStagingInput {
	return PublicationStagingInput{
		Authors:  cr.Authors(),
		Year:     NewOptional(cr.Year()),
		Title:    NewOptional(cr.TitleString()),
		Journal:  NewOptional(cr.Journal()),
		Verbatim: cr.Verbatim(),
		DOI:      NewOptional(types.DOI(cr.DOI)),
		Source:   biomedb.PublicationSourceCrossref,
	}
}

func (i PublicationStagingInput) ToDBParams() biomedb.StagePublicationsParams {
	return biomedb.StagePublicationsParams{
		Authors:  i.Authors,
		Year:     i.Year.ToPtr(),
		Title:    i.Title.ToPtr(),
		Journal:  i.Journal.ToPtr(),
		Verbatim: i.Verbatim,
		DOI:      i.DOI.ToPtr(),
		Source:   i.Source,
	}
}
