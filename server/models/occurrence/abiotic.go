package occurrence

import (
	"context"
	"darco/proto/db"
	"darco/proto/models"
	"darco/proto/models/people"
	"darco/proto/models/vocabulary"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
)

type AbioticParameter struct {
	vocabulary.Vocabulary `edgedb:"$inline" json:",inline"`
	Unit                  string      `edgedb:"unit" json:"unit"`
	Meta                  people.Meta `edgedb:"meta" json:"meta"`
}

func ListAbioticParameters(db edgedb.Executor) ([]AbioticParameter, error) {
	var items = []AbioticParameter{}
	err := db.Query(context.Background(),
		`select events::AbioticParameter { ** };`,
		&items)
	return items, err
}

type AbioticParameterInput struct {
	vocabulary.VocabularyInput `json:",inline"`
	Unit                       string `json:"unit"`
}

func (i AbioticParameterInput) Save(e edgedb.Executor) (created AbioticParameter, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data = <json>$0,
			select (insert events::AbioticParameter {
				label := <str>data['label'],
				code := <str>data['code'],
				description := <str>json_get(data, 'description'),
				unit := <str>data['unit']
			}) { ** }
		`, &created, data)
	return
}

type AbioticParameterUpdate struct {
	vocabulary.VocabularyUpdate `edgedb:"$inline" json:",inline"`
	Unit                        models.OptionalInput[string] `edgedb:"unit" json:"unit"`
}

func (u AbioticParameterUpdate) Save(e edgedb.Executor, code string) (updated AbioticParameter, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$1,
			select (update events::AbioticParameter filter .code = <str>$0 set {
				%s
			}) { ** }
		`,
		Mappings: u.FieldMappingsWith("item", map[string]string{
			"unit": "<str>item['unit']",
		}),
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, code, data)
	return
}

func DeleteAbioticParameter(db edgedb.Executor, label string) (deleted AbioticParameter, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
				delete events::AbioticParameter filter .label = <str>$0
			) { ** }
		`,
		&deleted, label)
	return
}

type AbioticMeasurement struct {
	ID    edgedb.UUID      `edgedb:"id" json:"id" format:"uuid"`
	Param AbioticParameter `edgedb:"param" json:"param"`
	Value float32          `edgedb:"value" json:"value"`
}

type AbioticMeasurementInput struct {
	Param string  `json:"param"` // Parameter code
	Value float32 `json:"value"`
}

func DeleteAbioticMeasurement(db edgedb.Executor, id edgedb.UUID) (deleted AbioticMeasurement, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
				delete events::AbioticMeasurement filter .id = <uuid>$0
			) { ** }
		`,
		&deleted, id)
	return
}
