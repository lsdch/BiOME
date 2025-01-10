package references

import (
	"context"
	"darco/proto/db"
	"darco/proto/models"
	"darco/proto/models/people"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
)

type Article struct {
	ID             edgedb.UUID         `edgedb:"id" json:"id" format:"uuid"`
	Code           string              `edgedb:"code" json:"code"`
	Authors        []string            `edgedb:"authors" json:"authors"`
	Year           int32               `edgedb:"year" json:"year"`
	Title          edgedb.OptionalStr  `edgedb:"title" json:"title,omitempty"`
	Journal        edgedb.OptionalStr  `edgedb:"journal" json:"journal,omitempty"`
	Verbatim       edgedb.OptionalStr  `edgedb:"verbatim" json:"verbatim,omitempty"`
	DOI            edgedb.OptionalStr  `edgedb:"doi" json:"doi,omitempty"`
	Comments       edgedb.OptionalStr  `edgedb:"comments" json:"comments,omitempty"`
	OriginalSource edgedb.OptionalBool `edgedb:"original_source" json:"original_source"`
	Meta           people.Meta         `edgedb:"meta" json:"meta"`
}

func ListArticles(db edgedb.Executor) ([]Article, error) {
	var items = []Article{}
	err := db.Query(context.Background(),
		`#edgeql
			select references::Article { ** } order by .authors[0] asc then .year desc;
		`,
		&items)
	return items, err
}

func DeleteArticle(db edgedb.Executor, code string) (deleted Article, err error) {
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
	Code     models.OptionalInput[string] `edgedb:"code" json:"code,omitempty"`
	Authors  []string                     `edgedb:"authors" json:"authors"`
	Year     int32                        `edgedb:"year" json:"year" minimum:"1500"`
	Title    models.OptionalInput[string] `edgedb:"title" json:"title,omitempty"`
	Journal  models.OptionalInput[string] `edgedb:"journal" json:"journal,omitempty"`
	Verbatim models.OptionalInput[string] `edgedb:"verbatim" json:"verbatim,omitempty"`
	Comments models.OptionalInput[string] `edgedb:"comments" json:"comments,omitempty"`
	DOI      models.OptionalInput[string] `edgedb:"doi" json:"doi,omitempty"`
}

func (i ArticleInput) Save(e edgedb.Executor) (created Article, err error) {
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
	Code     models.OptionalInput[string]   `edgedb:"code" json:"code,omitempty"`
	Authors  models.OptionalInput[[]string] `edgedb:"authors" json:"authors,omitempty"`
	Year     models.OptionalInput[int32]    `edgedb:"year" json:"year,omitempty" minimum:"1500"`
	Title    models.OptionalNull[string]    `edgedb:"title" json:"title,omitempty"`
	Journal  models.OptionalNull[string]    `edgedb:"journal" json:"journal,omitempty"`
	Verbatim models.OptionalNull[string]    `edgedb:"verbatim" json:"verbatim,omitempty"`
	Comments models.OptionalNull[string]    `edgedb:"comments" json:"comments,omitempty"`
	DOI      models.OptionalNull[string]    `edgedb:"doi" json:"doi,omitempty"`
}

func (u ArticleUpdate) Save(e edgedb.Executor, code string) (updated Article, err error) {
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
