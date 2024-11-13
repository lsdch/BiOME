package occurrence

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/people"
	"darco/proto/models/taxonomy"
	"darco/proto/models/vocabulary"
	"encoding/json"

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
	Kind       SamplingTargetKind `edgedb:"sampling_target" json:"kind"`
	TargetTaxa []taxonomy.Taxon   `edgedb:"target_taxa" json:"target_taxa"`
}

type Sampling struct {
	ID           edgedb.UUID           `edgedb:"id" json:"id" format:"uuid"`
	Methods      []SamplingMethod      `edgedb:"methods" json:"methods,omitempty"`
	Fixatives    []vocabulary.Fixative `edgedb:"fixatives" json:"fixatives"`
	Target       SamplingTarget        `edgedb:"$inline" json:"target"`
	Duration     edgedb.OptionalInt32  `edgedb:"sampling_duration" json:"duration,omitempty" doc:"Sampling duration in minutes"`
	Habitats     []Habitat             `edgedb:"habitats" json:"habitats"`
	AccessPoints []string              `edgedb:"access_points" json:"access_points"`
	Comments     edgedb.OptionalStr    `edgedb:"comments" json:"comments,omitempty"`
}

type SamplingInput struct {
	Methods   models.OptionalInput[[]string] `edgedb:"methods" json:"methods,omitempty"`
	Fixatives models.OptionalInput[[]string] `edgedb:"fixatives" json:"fixatives"`
	Target    SamplingTarget                 `edgedb:"target" json:"target"`
	Duration  models.OptionalInput[int32]    `edgedb:"duration" json:"duration,omitempty" doc:"Sampling duration in minutes"`
	Comments  models.OptionalInput[string]   `edgedb:"comments" json:"comments"`
}
