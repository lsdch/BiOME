package references

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/caltechlibrary/crossrefapi"
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/sirupsen/logrus"
)

type DOI string

func (d *DOI) UnmarshalEdgeDBStr(data []byte) error {
	return d.UnmarshalJSON(data)
}

func (d *DOI) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	s = strings.TrimPrefix(s, "https://doi.org/")
	s = strings.TrimPrefix(s, "http://doi.org/")
	s = strings.TrimPrefix(s, "doi:")
	s = strings.TrimPrefix(s, "doi/")
	*d = DOI(s)
	return nil
}

func ParseDOI(s string) DOI {
	s = strings.ToLower(s)
	s = strings.TrimPrefix(s, "https://doi.org/")
	s = strings.TrimPrefix(s, "http://doi.org/")
	s = strings.TrimPrefix(s, "doi:")
	s = strings.TrimPrefix(s, "doi/")
	return DOI(s)
}

func (d DOI) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(d), nil
}

func (d DOI) String() string {
	return string(d)
}

type Article struct {
	ID       geltypes.UUID        `gel:"id" json:"id" format:"uuid"`
	Code     string               `gel:"code" json:"code"`
	Authors  []string             `gel:"authors" json:"authors"`
	Year     int32                `gel:"year" json:"year"`
	Title    geltypes.OptionalStr `gel:"title" json:"title,omitempty"`
	Journal  geltypes.OptionalStr `gel:"journal" json:"journal,omitempty"`
	Verbatim geltypes.OptionalStr `gel:"verbatim" json:"verbatim,omitempty"`
	DOI      geltypes.OptionalStr `gel:"doi" json:"doi,omitempty"`
	Comments geltypes.OptionalStr `gel:"comments" json:"comments,omitempty"`
	Meta     people.Meta          `gel:"meta" json:"meta"`
}

func ListArticles(db geltypes.Executor) ([]Article, error) {
	var items = []Article{}
	err := db.Query(context.Background(),
		`#edgeql
			select references::Article { ** } order by .authors[0] asc then .year desc;
		`,
		&items)
	return items, err
}

func DeleteArticle(db geltypes.Executor, code string) (deleted Article, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
				delete references::Article filter .code = <str>$0
		 	) { ** };
		`,
		&deleted, code)
	return
}

type ArticleInput struct {
	Code             models.OptionalInput[string] `gel:"code" json:"code,omitempty"`
	CodeDiscriminant int                          `gel:"-" json:"-"`
	Authors          []string                     `gel:"authors" json:"authors"`
	Year             int32                        `gel:"year" json:"year" minimum:"1500"`
	Title            models.OptionalInput[string] `gel:"title" json:"title,omitempty"`
	Journal          models.OptionalInput[string] `gel:"journal" json:"journal,omitempty"`
	Verbatim         models.OptionalInput[string] `gel:"verbatim" json:"verbatim,omitempty"`
	Comments         models.OptionalInput[string] `gel:"comments" json:"comments,omitempty"`
	DOI              models.OptionalInput[string] `gel:"doi" json:"doi,omitempty"`
}

func (i *ArticleInput) GenerateCode() {
	code := ""
	code += strings.Trim(strings.Split(i.Authors[0], " ")[0], ", ")
	if len(i.Authors) == 2 {
		code += "_" + strings.Trim(strings.Split(i.Authors[1], " ")[0], ", ")
	} else if len(i.Authors) > 2 {
		code += "_et_al"
	}
	code += "_" + fmt.Sprintf("%d", i.Year)
	if i.CodeDiscriminant > 0 {
		code += string('a' + rune(i.CodeDiscriminant-1))
	}
	i.Code.SetValue(code)
	logrus.Debugf("Generated code '%s' for article: %+v", code, i)
}

func (i *ArticleInput) CheckCodeUniqueness(db geltypes.Executor) (isUnique bool, err error) {
	var existing bool
	err = db.QuerySingle(context.Background(),
		`#edgeql
			with code := <str>$0
			select exists(
				select references::Article filter .code ilike code
			)
		`,
		&existing, i.Code.Value,
	)
	return !existing, err
}

func (i *ArticleInput) EnsureUniqueCode(db geltypes.Executor) error {
	if !i.Code.IsSet {
		i.GenerateCode()
	}
	for {
		isUnique, err := i.CheckCodeUniqueness(db)
		if err != nil {
			return err
		}
		if isUnique {
			return nil
		}
		logrus.Warnf("Generated code '%s' is not unique, incrementing discriminant and regenerating", i.Code.Value)
		i.CodeDiscriminant++
		i.GenerateCode()
	}
}

func (i ArticleInput) Save(e geltypes.Executor) (created Article, err error) {
	for j, authors := range i.Authors {
		i.Authors[j] = strings.TrimSpace(authors)
	}

	if err = i.EnsureUniqueCode(e); err != nil {
		return
	}

	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0,
			select (insert references::Article {
				code := <str>data['code'],
				authors := <array<str>>data['authors'],
				year := <int32>data['year'],
				title := <str>json_get(data, "title"),
				journal := <str>json_get(data, "journal"),
				verbatim := <str>json_get(data, "verbatim"),
				comments := <str>json_get(data, "comments"),
				doi := <str>json_get(data, "doi"),
			}) { ** }
		`, &created, data)
	if err == nil {
		return
	}
	return
}

func ArticleInputFromCrossref(cr *crossrefapi.Works) ArticleInput {
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
	return ArticleInput{
		Authors: authors,
		Year:    year,
		Title:   models.NewOptionalInput(cr.Message.Title[0]),
		Journal: models.NewOptionalInput(cr.Message.ContainerTitle[0]),
		DOI:     models.NewOptionalInput(cr.Message.DOI),
	}
}

type ArticleUpdate struct {
	Code     models.OptionalInput[string]   `gel:"code" json:"code,omitempty"`
	Authors  models.OptionalInput[[]string] `gel:"authors" json:"authors,omitempty"`
	Year     models.OptionalInput[int32]    `gel:"year" json:"year,omitempty" minimum:"1500"`
	Title    models.OptionalNull[string]    `gel:"title" json:"title,omitempty"`
	Journal  models.OptionalNull[string]    `gel:"journal" json:"journal,omitempty"`
	Verbatim models.OptionalNull[string]    `gel:"verbatim" json:"verbatim,omitempty"`
	Comments models.OptionalNull[string]    `gel:"comments" json:"comments,omitempty"`
	DOI      models.OptionalNull[string]    `gel:"doi" json:"doi,omitempty"`
}

func (u ArticleUpdate) Save(e geltypes.Executor, code string) (updated Article, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$1,
			select (update references::Article filter .code = <str>$0 set {
				%s
			}) { ** }
		`,
		Mappings: map[string]string{
			"code":     "<str>item['code']",
			"authors":  "<array<str>>item['authors']",
			"year":     "<int32>item['year']",
			"title":    "<str>item['title']",
			"journal":  "<str>item['journal']",
			"comments": "<str>item['comments']",
			"verbatim": "<str>item['verbatim']",
			"doi":      "<str>item['doi']",
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, code, data)
	return
}
