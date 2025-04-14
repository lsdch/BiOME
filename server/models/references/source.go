package references

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/vocabulary"
)

type DataSource struct {
	vocabulary.Vocabulary `gel:"$inline" json:",inline"`
	LinkTemplate          geltypes.OptionalStr `gel:"link_template" json:"link_template,omitempty"`
	URL                   geltypes.OptionalStr `gel:"url" json:"url,omitempty"`
	Meta                  people.Meta          `gel:"meta" json:"meta"`
}

func ListDataSources(db geltypes.Executor) ([]DataSource, error) {
	var items = []DataSource{}
	err := db.Query(context.Background(),
		`#edgeql
			select references::DataSource { ** };
		`,
		&items)
	return items, err
}

func DeleteDataSources(db geltypes.Executor, code string) (deleted DataSource, err error) {
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
	vocabulary.VocabularyInput `gel:"$inline" json:",inline"`
	LinkTemplate               models.OptionalInput[string] `gel:"link_template" json:"link_template,omitempty"`
	URL                        models.OptionalInput[string] `gel:"url" json:"url,omitempty"`
}

func (i DataSourceInput) Save(e geltypes.Executor) (created DataSource, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0,
			select (insert references::DataSource {
				label := <str>data['label'],
				code := <str>data['code'],
				description := <str>json_get(data, 'description'),
				link_template := <str>json_get(data, 'link_template'),
				url := <str>json_get(data, 'url')
			}) { ** }
		`, &created, data)
	return
}

type DataSourceUpdate struct {
	vocabulary.VocabularyUpdate `gel:"$inline" json:",inline"`
	LinkTemplate                models.OptionalNull[string] `gel:"link_template" json:"link_template,omitempty"`
	URL                         models.OptionalNull[string] `gel:"url" json:"url,omitempty"`
}

func (u DataSourceUpdate) Save(e geltypes.Executor, code string) (updated DataSource, err error) {
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
