package references

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
)

type Article struct {
	ID             geltypes.UUID         `gel:"id" json:"id" format:"uuid"`
	Code           string                `gel:"code" json:"code"`
	Authors        []string              `gel:"authors" json:"authors"`
	Year           int32                 `gel:"year" json:"year"`
	Title          geltypes.OptionalStr  `gel:"title" json:"title,omitempty"`
	Journal        geltypes.OptionalStr  `gel:"journal" json:"journal,omitempty"`
	Verbatim       geltypes.OptionalStr  `gel:"verbatim" json:"verbatim,omitempty"`
	DOI            geltypes.OptionalStr  `gel:"doi" json:"doi,omitempty"`
	Comments       geltypes.OptionalStr  `gel:"comments" json:"comments,omitempty"`
	OriginalSource geltypes.OptionalBool `gel:"original_source" json:"original_source"`
	Meta           people.Meta           `gel:"meta" json:"meta"`
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
	Code     models.OptionalInput[string] `gel:"code" json:"code,omitempty"`
	Authors  []string                     `gel:"authors" json:"authors"`
	Year     int32                        `gel:"year" json:"year" minimum:"1500"`
	Title    models.OptionalInput[string] `gel:"title" json:"title,omitempty"`
	Journal  models.OptionalInput[string] `gel:"journal" json:"journal,omitempty"`
	Verbatim models.OptionalInput[string] `gel:"verbatim" json:"verbatim,omitempty"`
	Comments models.OptionalInput[string] `gel:"comments" json:"comments,omitempty"`
	DOI      models.OptionalInput[string] `gel:"doi" json:"doi,omitempty"`
}

func (i ArticleInput) Save(e geltypes.Executor) (created Article, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0,
			select (insert references::Article {
				authors := <array<str>>data['authors'],
				year := <int32>data['year'],
				title := <str>json_get(data, "title"),
				journal := <str>json_get(data, "journal"),
				verbatim := <str>json_get(data, "verbatim"),
				comments := <str>json_get(data, "comments"),
				doi := <str>json_get(data, "doi"),
			}) { ** }
		`, &created, data)
	return
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

// OccurrenceReference is a reference to an article for an occurrence
// with optional flag indicating whether the article is the original source
type OccurrenceReference struct {
	Article  `gel:"$inline" json:",inline"`
	Original bool `gel:"original_source" json:"original,omitempty"`
}

// OccurrenceReferenceInput helps binding Article references to occurrences,
// with optional original source flag
type OccurrenceReferenceInput struct {
	ArticleCode string                     `json:"code"`
	Original    models.OptionalInput[bool] `json:"original,omitempty"`
}
