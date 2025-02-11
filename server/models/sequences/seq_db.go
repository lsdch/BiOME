package sequences

import (
	"context"
	"encoding/json"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/vocabulary"

	"github.com/edgedb/edgedb-go"
)

type DataSource struct {
	vocabulary.Vocabulary `edgedb:"$inline" json:",inline"`
	LinkTemplate          edgedb.OptionalStr `edgedb:"link_template" json:"link_template,omitempty"`
	URL                   edgedb.OptionalStr `edgedb:"url" json:"url,omitempty"`
	Meta                  people.Meta        `edgedb:"meta" json:"meta"`
}

func ListDataSources(db edgedb.Executor) ([]DataSource, error) {
	var items = []DataSource{}
	err := db.Query(context.Background(),
		`#edgeql
			select references::DataSource { ** };
		`,
		&items)
	return items, err
}

func DeleteDataSources(db edgedb.Executor, code string) (deleted DataSource, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
			  delete references::DataSource filter .code = <str>$0
		 	) { ** };
		`,
		&deleted, code)
	return
}

type DataSourceInput struct {
	vocabulary.VocabularyInput `edgedb:"$inline" json:",inline"`
	LinkTemplate               models.OptionalInput[string] `edgedb:"link_template" json:"link_template,omitempty"`
	URL                        models.OptionalInput[string] `edgedb:"url" json:"url,omitempty"`
}

func (i DataSourceInput) Save(e edgedb.Executor) (created DataSource, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0,
			select (insert references::DataSource {
				label := <str>data['label'],
				code := <str>data['code'],
				description := <str>json_get(data, 'description'),
				link_template := <str>json_get(data, 'link_template')
				url := <str>json_get(data, 'url')
			}) { ** }
		`, &created, data)
	return
}

type DataSourceUpdate struct {
	vocabulary.VocabularyUpdate `edgedb:"$inline" json:",inline"`
	LinkTemplate                models.OptionalNull[string] `edgedb:"link_template" json:"link_template,omitempty"`
	URL                         models.OptionalNull[string] `edgedb:"url" json:"url,omitempty"`
}

func (u DataSourceUpdate) Save(e edgedb.Executor, code string) (updated DataSource, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$1,
			select (update references::DataSource filter .code = <str>$0 set {
				%s
			}) { ** }
		`,
		Mappings: u.FieldMappingsWith("item", map[string]string{
			"link_template": "<str>item['link_template']",
		}),
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, code, data)
	return
}
