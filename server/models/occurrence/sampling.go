package occurrence

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/taxonomy"
	"github.com/lsdch/biome/models/vocabulary"

	"github.com/sirupsen/logrus"
)

type SamplingMethod struct {
	vocabulary.Vocabulary `gel:"$inline" json:",inline"`
	Meta                  people.Meta `gel:"meta" json:"meta"`
}

func ListSamplingMethods(db geltypes.Executor) ([]SamplingMethod, error) {
	var items = []SamplingMethod{}
	err := db.Query(context.Background(),
		`select events::SamplingMethod { ** } order by .label`,
		&items)
	return items, err
}

type SamplingMethodInput vocabulary.VocabularyInput

func (i SamplingMethodInput) Save(db geltypes.Executor) (created SamplingMethod, err error) {
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

func (u SamplingMethodUpdate) Save(e geltypes.Executor, code string) (updated SamplingMethod, err error) {
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

func DeleteSamplingMethod(db geltypes.Executor, code string) (deleted SamplingMethod, err error) {
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
	Kind       SamplingTargetKind `gel:"sampling_target" json:"kind"`
	TargetTaxa []taxonomy.Taxon   `gel:"target_taxa" json:"taxa,omitempty"`
}

type SamplingOutline struct {
	ID          geltypes.UUID     `gel:"id" json:"id" format:"uuid"`
	Number      int64             `gel:"number" json:"number" doc:"Auto-incrementing number, unique per sampling"`
	PerformedOn DateWithPrecision `gel:"performed_on" json:"performed_on"`
}
type SamplingInner struct {
	SamplingOutline   `gel:"$inline" json:",inline"`
	PerformedBy       []people.PersonUser        `gel:"performed_by" json:"performed_by,omitempty"`
	PerformedByGroups []people.OrganisationInner `gel:"performed_by_groups" json:"performed_by_groups,omitempty"`
	Target            SamplingTarget             `gel:"$inline" json:"target"`
	Duration          geltypes.OptionalInt32     `gel:"sampling_duration" json:"duration,omitempty" doc:"Sampling duration in minutes"`
	Methods           []SamplingMethod           `gel:"methods" json:"methods,omitempty"`
	Fixatives         []vocabulary.Fixative      `gel:"fixatives" json:"fixatives,omitempty"`
	Habitats          []Habitat                  `gel:"habitats" json:"habitats,omitempty"`
	AccessPoints      []string                   `gel:"access_points" json:"access_points,omitempty"`
	Comments          geltypes.OptionalStr       `gel:"comments" json:"comments,omitempty"`
}

type Sampling struct {
	SamplingInner `gel:"$inline" json:",inline"`
	Samples       []BioMaterial    `gel:"samples" json:"samples,omitempty"`
	OccurringTaxa []taxonomy.Taxon `gel:"occurring_taxa" json:"occurring_taxa,omitempty"`
	Meta          people.Meta      `gel:"meta" json:"meta"`
}

func (s Sampling) Code(siteCode string) string {
	return fmt.Sprintf("%s|%s", siteCode, s.PerformedOn.ToCode())
}

type SamplingWithOccurrences struct {
	SamplingInner `gel:"$inline" json:",inline"`
	Occurrences   []OccurrenceAtSite `gel:"occurrences" json:"occurrences,omitempty"`
	OccurringTaxa []taxonomy.Taxon   `gel:"occurring_taxa" json:"occurring_taxa,omitempty"`
	Meta          people.Meta        `gel:"meta" json:"meta"`
}

type SamplingWithSite struct {
	Sampling `gel:"$inline" json:",inline"`
	Site     SiteItem `gel:"site" json:"site"`
}

type SamplingInnerWithSite struct {
	SamplingInner `gel:"$inline" json:",inline"`
	Site          SiteItem `gel:"site" json:"site"`
}

type SamplingInputAtSite struct {
	SamplingInput `json:",inline"`
	SiteCode      string `json:"site_code"`
}

func (i SamplingInputAtSite) Save(e geltypes.Executor) (Sampling, error) {
	return i.SamplingInput.Save(e, i.SiteCode)
}

type SamplingTargetInput struct {
	Kind SamplingTargetKind `json:"kind"`
	Taxa []string           `json:"taxa,omitempty"`
}

type SamplingInput struct {
	ActionInput  `json:",inline"`
	Target       SamplingTargetInput `json:"target"`
	Methods      []string            `json:"methods,omitempty"`
	Fixatives    []string            `json:"fixatives,omitempty"`
	Duration     *int32              `json:"duration,omitempty" doc:"Sampling duration in minutes"`
	Comments     *string             `json:"comments,omitempty"`
	Habitats     []string            `json:"habitats,omitempty"`
	AccessPoints []string            `json:"access_points,omitempty"`
}

func (i SamplingInput) Save(e geltypes.Executor, siteCode string) (created Sampling, err error) {
	data, _ := json.Marshal(i)
	logrus.Debugf("Inserting sampling event: %s", string(data))
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with module events,
			data := <json>$1,
			select (insert events::Sampling {
				site := assert_exists(
					(select location::Site filter .code = <str>$0),
					message := 'Site with code ' ++ <str>$0 ++ ' does not exist'
				),
				performed_by := (
					select people::Person
					filter .alias in <str>json_array_unpack(json_get(data, 'performed_by'))
				),
				performed_by_groups := (
					select people::Organisation
					filter .code in <str>json_array_unpack(json_get(data,'performed_by_groups'))
				),
				performed_on := (
					select date::from_json_with_precision(data['performed_on'])
				),
				methods := (
					select SamplingMethod
					filter .code in <str>json_array_unpack(json_get(data, 'methods'))
				),
				fixatives := (
					select samples::Fixative
					filter .code in <str>json_array_unpack(json_get(data, 'fixatives'))
				),
				sampling_target := <SamplingTarget>(data['target']['kind']),
				target_taxa := (
					select taxonomy::Taxon
					filter .name in <str>json_array_unpack(json_get(data, 'target', 'taxa'))
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
		`, &created, siteCode, data)
	return
}

type SamplingUpdate struct {
	PerformedBy       models.OptionalNull[[]string]                `gel:"performed_by" json:"performed_by,omitempty"`
	PerformedByGroups models.OptionalNull[[]string]                `gel:"performed_by_groups" json:"performed_by_groups,omitempty"`
	PerformedOn       models.OptionalInput[DateWithPrecisionInput] `gel:"performed_on" json:"performed_on,omitempty"`
	Target            models.OptionalInput[SamplingTargetInput]    `json:"target"`
	Methods           models.OptionalNull[[]string]                `gel:"methods" json:"methods,omitempty"`
	Fixatives         models.OptionalNull[[]string]                `gel:"fixatives" json:"fixatives,omitempty"`
	Duration          models.OptionalNull[int32]                   `gel:"duration" json:"duration,omitempty" doc:"Sampling duration in minutes"`
	Comments          models.OptionalNull[string]                  `gel:"comments" json:"comments,omitempty"`
	Habitats          models.OptionalNull[[]string]                `gel:"habitats" json:"habitats,omitempty"`
	AccessPoints      models.OptionalNull[[]string]                `gel:"access_points" json:"access_points,omitempty"`
}

func (u SamplingUpdate) Save(e geltypes.Executor, id geltypes.UUID) (updated Sampling, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with data := <json>$1,
			select (update events::Sampling filter .id = <uuid>$0 set {
				%s
			}) {
				*,
				site: {*, country: { *}},
				performed_by: { * },
				performed_by_groups: { * },
				habitats: { * },
				target_taxa: { * },
				fixatives: { * },
				methods: { * },
				meta: { * }
			}
		`,
		Mappings: map[string]string{
			"performed_by": `#edgeql
				(
					select people::Person
					filter .alias in <str>json_array_unpack(data['performed_by'])
				)`,
			"performed_by_groups": `#edgeql
				(
					select people::Organisation
					filter .code in <str>json_array_unpack(data['performed_by_groups'])
				)`,
			"performed_on": `#edgeql
				date::from_json_with_precision(data['performed_on'])
			`,
			"sampling_target": "<events::SamplingTarget>item['target_kind']",
			"target_taxa": `#edgeql
				(
					select taxonomy::Taxon
					filter .name in <str>json_array_unpack(json_get(data, 'target_taxa'))
				)`,
			"methods": `#edgeql
				(
					select events::SamplingMethod
					filter .code in <str>json_array_unpack(json_get(data, 'methods'))
				)`,
			"fixatives": `#edgeql
				(
					select samples::Fixative
					filter .code in <str>json_array_unpack(json_get(data, 'fixatives'))
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

func ListAccessPoints(db geltypes.Executor) ([]string, error) {
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

func DeleteSampling(db geltypes.Executor, id geltypes.UUID) (deleted Sampling, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
			 	delete events::Sampling filter .id = <uuid>$0
		 	) {
			 	*,
				site: { *, country: { * }},
				performed_by: { * },
				performed_by_groups: { * },
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
