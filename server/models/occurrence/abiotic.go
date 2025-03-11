package occurrence

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/vocabulary"
)

type AbioticParameter struct {
	vocabulary.Vocabulary `gel:"$inline" json:",inline"`
	Unit                  string      `gel:"unit" json:"unit"`
	Meta                  people.Meta `gel:"meta" json:"meta"`
}

func ListAbioticParameters(db geltypes.Executor) ([]AbioticParameter, error) {
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

func (i AbioticParameterInput) Save(e geltypes.Executor) (created AbioticParameter, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0,
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
	vocabulary.VocabularyUpdate `gel:"$inline" json:",inline"`
	Unit                        models.OptionalInput[string] `gel:"unit" json:"unit"`
}

func (u AbioticParameterUpdate) Save(e geltypes.Executor, code string) (updated AbioticParameter, err error) {
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

func DeleteAbioticParameter(db geltypes.Executor, code string) (deleted AbioticParameter, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
				delete events::AbioticParameter filter .code = <str>$0
			) { ** }
		`,
		&deleted, code)
	return
}

type AbioticMeasurement struct {
	ID    geltypes.UUID    `gel:"id" json:"id" format:"uuid"`
	Param AbioticParameter `gel:"param" json:"param"`
	Value float32          `gel:"value" json:"value"`
}

type AbioticMeasurementInput struct {
	Param string  `json:"param"` // Parameter code
	Value float32 `json:"value"`
}

// Upsert AbioticMeasurement with the given event ID
func (u AbioticMeasurementInput) Save(e geltypes.Executor, eventID geltypes.UUID) (updated AbioticMeasurement, err error) {
	data, _ := json.Marshal(u)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with item := <json>$1,
			param := (select events::AbioticParameter filter .code = <str>item['param']),
			select (
				insert events::AbioticMeasurement {
					event := (<events::Event><uuid>$0),
					param := param,
					value := <float32>item['value']
				} unless conflict on ((.event, .param)) else (
					update events::AbioticMeasurement set {
						param := param,
						value := <float32>item['value']
					}
				)
			) { param: { * }, value}
		`,
		&updated, eventID, data)
	return
}

func DeleteAbioticMeasurement(db geltypes.Executor, id geltypes.UUID) (deleted AbioticMeasurement, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
				delete events::AbioticMeasurement filter .id = <uuid>$0
			) { ** }
		`,
		&deleted, id)
	return
}
