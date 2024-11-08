package occurrence

import (
	"context"
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

func (i AbioticParameterInput) Create(e edgedb.Executor) (created AbioticParameter, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`with data = <json>$0,
		select (insert events::AbioticParameter {
			 label := <str>data['label'],
			 code := <str>data['code'],
			 description := <str>json_get(data, 'description'),
			 unit := <str>data['unit']
		})`, &created, data)
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
