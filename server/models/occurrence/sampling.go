package occurrence

import (
	"context"
	"darco/proto/db"
	"darco/proto/models"
	"darco/proto/models/people"
	"darco/proto/models/taxonomy"
	"darco/proto/models/vocabulary"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type SamplingMethod struct {
	vocabulary.Vocabulary `edgedb:"$inline" json:",inline"`
	Meta                  people.Meta `edgedb:"meta" json:"meta"`
}

func ListSamplingMethods(db edgedb.Executor) ([]SamplingMethod, error) {
	var items = []SamplingMethod{}
	err := db.Query(context.Background(),
		`select events::SamplingMethod { ** } order by .label`,
		&items)
	return items, err
}

type SamplingMethodInput vocabulary.VocabularyInput

func (i SamplingMethodInput) Save(db edgedb.Executor) (created SamplingMethod, err error) {
	data, _ := json.Marshal(i)
	err = db.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0
			select (insert events::SamplingMethod {
				label := <str>data['label'],
				code := <str>data['code'],
				description := <str>json_get(data, 'description'),
			}) { ** }
		`, &created, data)
	return
}

type SamplingMethodUpdate vocabulary.VocabularyUpdate

func (u SamplingMethodUpdate) Save(e edgedb.Executor, code string) (updated SamplingMethod, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$1,
			select (update events::SamplingMethod filter .code = <str>$0 set {
				%s
			}) { ** }
		`,
		Mappings: vocabulary.VocabularyUpdate(u).FieldMappingsWith("item"),
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, code, data)
	return
}

func DeleteSamplingMethod(db edgedb.Executor, code string) (deleted SamplingMethod, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
				delete events::SamplingMethod filter .code = <str>$0
			) { ** }
		`,
		&deleted, code)
	return
}

type SamplingTarget struct {
	Kind       SamplingTargetKind `edgedb:"sampling_target" json:"kind"`
	TargetTaxa []taxonomy.Taxon   `edgedb:"target_taxa" json:"target_taxa,omitempty"`
}

type SamplingInner struct {
	ID           edgedb.UUID           `edgedb:"id" json:"id" format:"uuid"`
	Number       int64                 `edgedb:"number" json:"-"`
	Code         string                `edgedb:"code" json:"code"`
	Target       SamplingTarget        `edgedb:"$inline" json:"target"`
	Duration     edgedb.OptionalInt32  `edgedb:"sampling_duration" json:"duration,omitempty" doc:"Sampling duration in minutes"`
	Methods      []SamplingMethod      `edgedb:"methods" json:"methods"`
	Fixatives    []vocabulary.Fixative `edgedb:"fixatives" json:"fixatives"`
	Habitats     []Habitat             `edgedb:"habitats" json:"habitats"`
	AccessPoints []string              `edgedb:"access_points" json:"access_points"`
	Comments     edgedb.OptionalStr    `edgedb:"comments" json:"comments,omitempty"`
}

type Sampling struct {
	SamplingInner `edgedb:"$inline" json:",inline"`
	Samples       []BioMaterial    `edgedb:"samples" json:"samples"`
	OccurringTaxa []taxonomy.Taxon `edgedb:"occurring_taxa" json:"occurring_taxa"`
	Meta          people.Meta      `edgedb:"meta" json:"meta"`
}

type SamplingInput struct {
	EventID      edgedb.UUID        `json:"event_id"`
	TargetKind   SamplingTargetKind `json:"target_kind"`
	TargetTaxa   []string           `json:"target_taxa,omitempty"`
	Methods      []string           `json:"methods,omitempty"`
	Fixatives    []string           `json:"fixatives,omitempty"`
	Duration     *int32             `json:"duration,omitempty" doc:"Sampling duration in minutes"`
	Comments     *string            `json:"comments,omitempty"`
	Habitats     []string           `json:"habitats,omitempty"`
	AccessPoints []string           `json:"access_points,omitempty"`
}

func (i SamplingInput) Save(e edgedb.Executor) (created Sampling, err error) {
	data, _ := json.Marshal(i)
	logrus.Debugf("data: %s", string(data))
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with module events,
			data := <json>$0,
			select (insert events::Sampling {
				event := (select Event filter .id = (<uuid>data['event_id'])),
				methods := (
					select SamplingMethod
					filter .code in <str>json_array_unpack(json_get(data, 'methods'))
				),
				fixatives := (
					select samples::Fixative
					filter .code in <str>json_array_unpack(json_get(data, 'fixatives'))
				),
				sampling_target := <SamplingTarget>(data['target_kind']),
				target_taxa := (
					select taxonomy::Taxon
					filter .name in <str>json_array_unpack(json_get(data, 'target_taxa'))
				),
				sampling_duration := <int32>json_get(data, 'duration'),
				comments := <str>json_get(data, 'comments'),
				habitats := (
					select sampling::Habitat
					filter .label in <str>json_array_unpack(json_get(data, 'habitats'))
				),
				access_points := (<str>json_array_unpack(json_get(data, 'access_points')))
			}) {
				*,
				habitats: { * },
				target_taxa: { * },
				fixatives: { * },
				methods: { * },
				meta: { * }
			}
		`, &created, data)
	return
}

type SamplingUpdate struct {
	TargetKind   models.OptionalInput[SamplingTargetKind] `edgedb:"target_kind" json:"target_kind,omitempty"`
	TargetTaxa   models.OptionalNull[[]string]            `edgedb:"target_taxa" json:"target_taxa,omitempty"`
	Methods      models.OptionalNull[[]string]            `edgedb:"methods" json:"methods,omitempty"`
	Fixatives    models.OptionalNull[[]string]            `edgedb:"fixatives" json:"fixatives,omitempty"`
	Duration     models.OptionalNull[int32]               `edgedb:"duration" json:"duration,omitempty" doc:"Sampling duration in minutes"`
	Comments     models.OptionalNull[string]              `edgedb:"comments" json:"comments,omitempty"`
	Habitats     models.OptionalNull[[]string]            `edgedb:"habitats" json:"habitats,omitempty"`
	AccessPoints models.OptionalNull[[]string]            `edgedb:"access_points" json:"access_points,omitempty"`
}

func (u SamplingUpdate) Save(e edgedb.Executor, id edgedb.UUID) (updated Sampling, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with data := <json>$1,
			select (update events::Sampling filter .id = <uuid>$0 set {
				%s
			}) {
				*,
				habitats: { * },
				target_taxa: { * },
				fixatives: { * },
				methods: { * },
				meta: { * }
			}
		`,
		Mappings: map[string]string{
			"sampling_target": "<events::SamplingTarget>item['target_kind']",
			"target_taxa": `#edgeql
				(
					select taxonomy::Taxon
					filter .name in <str>json_array_unpack(json_get(data, 'target_taxa'))
				)`,
			"methods": `#edgeql
				(
					select events::SamplingMethod
					filter .label in <str>json_array_unpack(json_get(data, 'methods'))
				)`,
			"fixatives": `#edgeql
				(
					select samples::Fixative
					filter .label in <str>json_array_unpack(json_get(data, 'fixatives'))
				)`,
			"sampling_duration": "<int32>json_get(data, 'duration')",
			"comments":          "<str>json_get(data, 'comments')",
			"habitats": `#edgeql
				(
					select sampling::Habitat
					filter .label in <str>json_array_unpack(json_get(data, 'habitats'))
				)`,
			"access_points": "<str>json_array_unpack(json_get(data, 'access_points'))",
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, id, data)
	return
}

func ListAccessPoints(db edgedb.Executor) ([]string, error) {
	var accessPoints []string
	err := db.Query(context.Background(),
		`#edgeql
			with a := (select distinct events::Sampling.access_points)
			select a order by a asc
		`,
		&accessPoints,
	)
	return accessPoints, err
}

func DeleteSampling(db edgedb.Executor, id edgedb.UUID) (deleted Sampling, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
			 	delete events::Sampling filter .id = <uuid>$0
		 	) {
			 	*,
				habitats: { * },
				target_taxa: { * },
				fixatives: { * },
				methods: { * },
				meta: { * }
			}
		`,
		&deleted, id)
	return
}
