package csvmodels

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/bibliography"
	"github.com/lsdch/biome/types"
)

type PublicationImportRow struct {
	rowNumber int32      `csv:"row_number" validate:"required,gt=0"`
	Authors   *string    `csv:"authors,omitempty"`
	Year      *int32     `csv:"year,omitempty"`
	Title     *string    `csv:"title,omitempty"`
	Journal   *string    `csv:"journal,omitempty"`
	Verbatim  *string    `csv:"verbatim,omitempty"`
	DOI       *types.DOI `csv:"doi,omitempty"`
}

func (r *PublicationImportRow) SetRowNumber(rowNumber int32) {
	r.rowNumber = rowNumber
}

func (r PublicationImportRow) RowNumber() int32 {
	return r.rowNumber
}

func (r PublicationImportRow) Validate(v *validator.Validate) error {
	return v.Struct(r)
}
func (r PublicationImportRow) String() string {
	if r.DOI != nil {
		return r.DOI.String()
	}
	if r.Verbatim != nil {
		return *r.Verbatim
	}
	return ""
}

type PublicationResolutionInput struct {
	Authors    *string    `csv:"authors,omitempty"`
	Year       *int32     `csv:"year,omitempty"`
	Title      *string    `csv:"title,omitempty"`
	Journal    *string    `csv:"journal,omitempty"`
	Verbatim   *string    `csv:"verbatim,omitempty"`
	DOI        *types.DOI `csv:"DOI,omitempty"`
	RowNumbers []int32    `csv:"row_numbers,omitempty"`
}

func (i PublicationResolutionInput) ToParams(importID uuid.UUID) biomedb.InitBibliographyResolutionParams {
	var (
		doi        = i.DOI
		year       = i.Year
		authorsRaw = i.Authors
		authors    []string
	)
	// If DOI is not present, attempt to extract it from the verbatim string
	if doi == nil {
		doi = bibliography.ExtractDOI(*i.Verbatim)
	}
	// Extract year and authors from verbatim if DOI is not present
	if doi == nil {
		if year == nil {
			year = bibliography.ExtractYear(*i.Verbatim)
		}
		if authorsRaw == nil {
			authorsStr := bibliography.ExtractAuthorString(*i.Verbatim)
			if authorsStr != "" {
				authorsRaw = &authorsStr
			}
		}
		if authorsRaw != nil && len(authors) == 0 {
			authors = bibliography.ParseAuthors(*authorsRaw)
		}
	}

	return biomedb.InitBibliographyResolutionParams{
		ImportID:   importID,
		DOI:        doi,
		Verbatim:   i.Verbatim,
		RowNumbers: i.RowNumbers,
		AuthorsRaw: authorsRaw,
		Authors:    authors,
		Year:       year,
		Title:      i.Title,
		Journal:    i.Journal,
	}
}
