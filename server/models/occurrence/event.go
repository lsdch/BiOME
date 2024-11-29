package occurrence

import (
	"context"
	"darco/proto/db"
	"darco/proto/models"
	"darco/proto/models/people"
	"encoding/json"
	"time"

	"github.com/edgedb/edgedb-go"
)

type DateWithPrecision struct {
	Date      time.Time     `edgedb:"date" json:"date,omitempty"`
	Precision DatePrecision `edgedb:"precision" json:"precision"`
}

type SiteInfo struct {
	Name string `edgedb:"name" json:"name"`
	Code string `edgedb:"code" json:"code"`
}

type Event struct {
	ID                  edgedb.UUID          `edgedb:"id" json:"id" format:"uuid"`
	Site                SiteInfo             `edgedb:"site" json:"site"`
	PerformedBy         []people.PersonUser  `edgedb:"performed_by" json:"performed_by" minLength:"1"`
	PerformedOn         DateWithPrecision    `edgedb:"performed_on" json:"performed_on"`
	Programs            []ProgramInner       `edgedb:"programs" json:"programs,omitempty"`
	AbioticMeasurements []AbioticMeasurement `edgedb:"abiotic_measurements" json:"abiotic_measurements"`
	Samplings           []Sampling           `edgedb:"samplings" json:"samplings"`
	Spotting            Spotting             `edgedb:"spotting" json:"spotting"`
	Meta                people.Meta          `edgedb:"meta" json:"meta"`
}

func (e *Event) AddSampling(db edgedb.Executor, sampling SamplingInput) error {
	sampling.EventID = e.ID
	created, err := sampling.Save(db)
	if err != nil {
		return err
	}
	e.Samplings = append(e.Samplings, created)
	return nil
}

func ListEvents(db edgedb.Executor) ([]Event, error) {
	var items = []Event{}
	err := db.Query(context.Background(),
		`#edgeql
			select events::Event {
				site: {name, code},
				programs: { * },
				performed_by: { * },
				spotting: { *, target_taxa: { * } },
				abiotic_measurements: { *, param: { * }  },
				samplings: { *, target_taxa: { * }, fixatives: { * }, methods: { * }, habitats: { * } }
			}
		`,
		&items)
	return items, err
}

type EventInput struct {
	PerformedBy []string                       `json:"performed_by" minLength:"1"`
	PerformedOn DateWithPrecision              `json:"performed_on"`
	Programs    models.OptionalInput[[]string] `json:"programs,omitempty"`
}

func (i EventInput) Save(e edgedb.Executor, site_code string) (created Event, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data = <json>$1,
			select (insert events::Event {
				site := (
					select location::Site filter .code = <str>$0
				),
				performed_by := (
					select people::Person filter .alias = <str>data['performed_by']
				),
				performed_on := (
					date := <datetime>data['performed_on']['date'],
					precision := <date::DatePrecision>data['performed_on']['precision']
				),
				programs := (
					select events::Program filter .code in json_array_unpack(<array<str>>json_get(data, 'programs'))
				)
			}) {
				site: {name, code},
				programs: { * },
				performed_by: { * },
				spotting: { *, target_taxa: { * } },
				abiotic_measurements: { *, param: { * }  },
				samplings: { *, target_taxa: { * }, fixatives: { * }, methods: { * }, habitats: { * } }
			}
		`, &created, site_code, data)
	return
}

type EventUpdate struct {
	PerformedBy models.OptionalInput[[]string]          `json:"performed_by,omitempty"`
	PerformedOn models.OptionalInput[DateWithPrecision] `json:"performed_on"`
	Programs    models.OptionalNull[[]string]           `json:"programs,omitempty"`
}

func (u EventUpdate) Save(e edgedb.Executor, id edgedb.UUID) (updated Event, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with data := <json>$1,
			select (update events::Event filter .id = <uuid>$0 set {
				%s
			}) {
				site: {name, code},
				programs: { * },
				performed_by: { * },
				spotting: { *, target_taxa: { * } },
				abiotic_measurements: { *, param: { * }  },
				samplings: { *, target_taxa: { * }, fixatives: { * }, methods: { * }, habitats: { * } }
			}
		`,
		Mappings: map[string]string{
			"perform_by": `#edgeql
				(
					select people::Person
					filter .alias in <str>json_array_unpack(data['performed_by'])
				)`,
			"performed_on": `#edgeql
				(
					date := <datetime>data['performed_on']['date'],
					precision := <date::DatePrecision>data['performed_on']['precision']
				)`,
			"programs": `#edgeql
				(
					select events::Program
					filter .code in <str>json_array_unpack(json_get(data, 'programs'))
				)`,
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, id, data)
	return
}

func DeleteEvent(db edgedb.Executor, id edgedb.UUID) (deleted Event, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
				delete events::Event filter .id = <uuid>$0
			) {
				site: {name, code},
				programs: { * },
				performed_by: { * },
				spotting: { *, target_taxa: { * } },
				abiotic_measurements: { *, param: { * }  },
				samplings: { *, target_taxa: { * }, fixatives: { * }, methods: { * }, habitats: { * } }
			};
		`,
		&deleted, id)
	return
}
