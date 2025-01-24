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

type CompositeDate struct {
	Day   int32 `json:"day,omitempty" minimum:"1" maximum:"31" default:"1"`
	Month int32 `json:"month,omitempty" minimum:"1" maximum:"12" default:"1"`
	Year  int32 `json:"year,omitempty" minimum:"1500" maximum:"3000"`
}

type DateWithPrecisionInput struct {
	Date      CompositeDate `json:"date"`
	Precision DatePrecision `json:"precision"`
}

type SiteInfo struct {
	Name string `edgedb:"name" json:"name"`
	Code string `edgedb:"code" json:"code"`
}

type EventInner struct {
	ID          edgedb.UUID        `edgedb:"id" json:"id" format:"uuid"`
	Site        SiteInfo           `edgedb:"site" json:"site"`
	Code        string             `edgedb:"code" json:"code"`
	PerformedOn DateWithPrecision  `edgedb:"performed_on" json:"performed_on"`
	Comments    edgedb.OptionalStr `edgedb:"comments" json:"comments,omitempty"`
}

type Event struct {
	EventInner          `edgedb:"$inline" json:",inline"`
	PerformedBy         []people.PersonUser  `edgedb:"performed_by" json:"performed_by" minLength:"1"`
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
// AddAbioticMeasurement adds an abiotic measurement to the event.
// If a value for a given parameter already exists, it will be overwritten.
func (e *Event) AddAbioticMeasurement(db edgedb.Executor, measurements AbioticMeasurementInput) error {
	created, err := measurements.Save(db, e.ID)
	if err != nil {
		return err
	}
	e.AbioticMeasurements = append(e.AbioticMeasurements, created)
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
				meta: { * }
			}
		`,
		&items)
	return items, err
}

type EventInput struct {
	PerformedBy []string                       `json:"performed_by" minLength:"1"`
	PerformedOn DateWithPrecisionInput         `json:"performed_on"`
	Programs    models.OptionalInput[[]string] `json:"programs,omitempty"`
}

func (i EventInput) Save(e edgedb.Executor, site_code string) (created Event, err error) {

	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$1,
			select (insert events::Event {
				site := (
					select location::Site filter .code = <str>$0
				),
				performed_by := (
					select people::Person filter .alias in <str>json_array_unpack(data['performed_by'])
				),
				performed_on := date::from_json_with_precision(data['performed_on']),
				programs := (
					select events::Program filter .code in <str>json_array_unpack(json_get(data, 'programs'))
				)
			}) {
				*,
				site: {name, code},
				programs: { * },
				performed_by: { * },
				spotting: { *, target_taxa: { * } },
				abiotic_measurements: { *, param: { * }  },
				samplings: { *, target_taxa: { * }, fixatives: { * }, methods: { * }, habitats: { * } },
				meta: { * }
			}
		`, &created, site_code, data)
	return
}

type EventUpdate struct {
	PerformedBy models.OptionalInput[[]string]               `edgedb:"performed_by" json:"performed_by,omitempty"`
	PerformedOn models.OptionalInput[DateWithPrecisionInput] `edgedb:"performed_on" json:"performed_on,omitempty"`
	Programs    models.OptionalNull[[]string]                `edgedb:"programs" json:"programs,omitempty"`
}

func (u EventUpdate) Save(e edgedb.Executor, id edgedb.UUID) (updated Event, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with data := <json>$1,
			select (update events::Event filter .id = <uuid>$0 set {
				%s
			}) {
				*,
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
				date::from_json_with_precision(data['performed_on'])
			`,
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
