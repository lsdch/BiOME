package events

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/people"
	"darco/proto/models/taxonomy"
	"darco/proto/models/vocabulary"
	"encoding/json"
	"time"

	"github.com/edgedb/edgedb-go"
)

type SamplingMethod struct {
	vocabulary.Vocabulary `edgedb:"$inline" json:",inline"`
	Meta                  people.Meta `edgedb:"meta" json:"meta"`
}

func ListSamplingMethods(db edgedb.Executor) ([]SamplingMethod, error) {
	var items = []SamplingMethod{}
	err := db.Query(context.Background(),
		`select events::SamplingMethod { ** };`,
		&items)
	return items, err
}

type SamplingMethodInput struct {
	vocabulary.VocabularyInput `json:",inline"`
}

func (i SamplingMethodInput) Create(db edgedb.Executor) (created SamplingMethod, err error) {
	data, _ := json.Marshal(i)
	err = db.QuerySingle(context.Background(),
		`select (insert events::SamplingMethod {
			label := <str>data['label'],
			code := <str>data['code'],
			description := <str>json_get(data, 'description'),
		})`, &created, data)
	return
}

type SamplingTarget struct {
	Kind       SamplingTargetKind    `edgedb:"kind" json:"kind"`
	TargetTaxa []taxonomy.TaxonInner `edgedb:"target_taxa" json:"target_taxa"`
}

type Sampling struct {
	Code       string                `edgedb:"code" json:"code"`
	Methods    []SamplingMethod      `edgedb:"methods" json:"methods,omitempty"`
	Fixatives  []vocabulary.Fixative `edgedb:"fixatives" json:"fixatives"`
	Target     SamplingTarget        `edgedb:"target" json:"target"`
	Duration   edgedb.Duration       `edgedb:"sampling_duration" json:"duration"`
	IsDonation bool                  `edgedb:"is_donation" json:"is_donation"`
	Comments   edgedb.OptionalStr    `edgedb:"comments" json:"comments"`
}

type SamplingInput struct {
	Code       models.OptionalInput[string]        `edgedb:"code" json:"code,omitempty"`
	Methods    models.OptionalInput[[]string]      `edgedb:"methods" json:"methods,omitempty"`
	Fixatives  models.OptionalInput[[]string]      `edgedb:"fixatives" json:"fixatives"`
	Target     SamplingTarget                      `edgedb:"target" json:"target"`
	Duration   models.OptionalInput[time.Duration] `edgedb:"duration" json:"duration,omitempty"`
	IsDonation bool                                `edgedb:"is_donation" json:"is_donation"`
	Comments   models.OptionalInput[string]        `edgedb:"comments" json:"comments"`
}
