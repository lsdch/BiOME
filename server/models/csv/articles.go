package csvmodels

type ArticleInput struct {
	Authors  StringListInput `csv:"authors,omitempty"`
	Year     *int32          `csv:"year,omitempty"`
	Title    *string         `csv:"title,omitempty"`
	Journal  *string         `csv:"journal,omitempty"`
	Verbatim *string         `csv:"verbatim,omitempty"`
	DOI      *string         `csv:"DOI,omitempty"`
	Comments *string         `csv:"comments,omitempty"`
}
