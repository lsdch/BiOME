package models

import (
	"encoding/json"
	"strings"

	"github.com/caltechlibrary/crossrefapi"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/sirupsen/logrus"
)

type DOI string

func normalizeDOI(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "https://doi.org/")
	s = strings.TrimPrefix(s, "http://doi.org/")
	s = strings.TrimPrefix(s, "doi:")
	s = strings.TrimPrefix(s, "doi/")
	s = strings.TrimPrefix(s, "doi.org/")
	return s
}

func (d *DOI) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	s = normalizeDOI(s)
	*d = DOI(s)
	return nil
}

func ParseDOI(s string) DOI {
	s = normalizeDOI(s)
	return DOI(s)
}

func (d DOI) String() string {
	return string(d)
}

type Article struct {
	ID       uuid.UUID        `json:"id"`
	Authors  []string         `json:"authors"`
	Year     int32            `json:"year"`
	Title    Optional[string] `json:"title,omitempty"`
	Journal  Optional[string] `json:"journal,omitempty"`
	Verbatim Optional[string] `json:"verbatim,omitempty"`
	DOI      Optional[string] `json:"doi,omitempty"`
	Comments Optional[string] `json:"comments,omitempty"`
}

func ArticleFromDB(a biomedb.Article) Article {
	return Article{
		ID:       a.ID,
		Authors:  a.Authors,
		Year:     a.Year,
		Title:    NewOptionalFromPtr(a.Title),
		Journal:  NewOptionalFromPtr(a.Journal),
		Verbatim: NewOptionalFromPtr(a.Verbatim),
		DOI:      NewOptionalFromPtr(a.Doi),
		Comments: NewOptionalFromPtr(a.Comments),
	}
}

type CreateArticleParams struct {
	Authors  []string         `json:"authors"`
	Year     int32            `json:"year"`
	Title    Optional[string] `json:"title,omitempty"`
	Journal  Optional[string] `json:"journal,omitempty"`
	Verbatim Optional[string] `json:"verbatim,omitempty"`
	DOI      Optional[string] `json:"doi,omitempty"`
	Comments Optional[string] `json:"comments,omitempty"`
}

func ArticleParamsFromCrossref(cr *crossrefapi.Works) CreateArticleParams {
	logrus.Debugf("Creating ArticleInput from Crossref data: %+v", cr)
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
	return CreateArticleParams{
		Authors: authors,
		Year:    year,
		Title:   NewOptional(cr.Message.Title[0]),
		Journal: NewOptional(cr.Message.ContainerTitle[0]),
		DOI:     NewOptional(cr.Message.DOI),
	}
}

func (p CreateArticleParams) ToDBParams() biomedb.CreateArticleParams {
	return biomedb.CreateArticleParams{
		Authors:  p.Authors,
		Year:     p.Year,
		Title:    p.Title.ToPtr(),
		Journal:  p.Journal.ToPtr(),
		Verbatim: p.Verbatim.ToPtr(),
		Doi:      p.DOI.ToPtr(),
		Comments: p.Comments.ToPtr(),
	}
}

type UpdateArticleParams struct {
	Code     Optional[string]     `json:"code,omitempty"`
	Authors  Optional[[]string]   `json:"authors,omitempty"`
	Year     Optional[int32]      `json:"year,omitempty" minimum:"1500"`
	Title    OptionalNull[string] `json:"title,omitempty"`
	Journal  OptionalNull[string] `json:"journal,omitempty"`
	Verbatim OptionalNull[string] `json:"verbatim,omitempty"`
	Comments OptionalNull[string] `json:"comments,omitempty"`
	DOI      OptionalNull[string] `json:"doi,omitempty"`
}

func (p UpdateArticleParams) ToDBParams(id uuid.UUID) biomedb.UpdateArticleByIDParams {
	return biomedb.UpdateArticleByIDParams{
		ID:          id,
		Authors:     p.Authors.Value,
		Year:        p.Year.ToPtr(),
		SetTitle:    p.Title.IsSet,
		Title:       p.Title.ToPtr(),
		SetJournal:  p.Journal.IsSet,
		Journal:     p.Journal.ToPtr(),
		SetVerbatim: p.Verbatim.IsSet,
		Verbatim:    p.Verbatim.ToPtr(),
		SetComments: p.Comments.IsSet,
		Comments:    p.Comments.ToPtr(),
		SetDoi:      p.DOI.IsSet,
		Doi:         p.DOI.ToPtr(),
	}
}
