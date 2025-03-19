package occurrence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/taxonomy"
)

type DateWithPrecision struct {
	Date      time.Time     `gel:"date" json:"date,omitempty"`
	Precision DatePrecision `gel:"precision" json:"precision"`
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
	Name string `gel:"name" json:"name"`
	Code string `gel:"code" json:"code"`
}

type EventInner struct {
	ID          geltypes.UUID        `gel:"id" json:"id" format:"uuid"`
	Site        SiteInfo             `gel:"site" json:"site"`
	Code        string               `gel:"code" json:"code"`
	PerformedOn DateWithPrecision    `gel:"performed_on" json:"performed_on"`
	Comments    geltypes.OptionalStr `gel:"comments" json:"comments,omitempty"`
}

type Event struct {
	EventInner          `gel:"$inline" json:",inline"`
	PerformedBy         []people.PersonUser  `gel:"performed_by" json:"performed_by" minLength:"1"`
	Programs            []ProgramInner       `gel:"programs" json:"programs,omitempty"`
	AbioticMeasurements []AbioticMeasurement `gel:"abiotic_measurements" json:"abiotic_measurements,omitempty"`
	Samplings           []Sampling           `gel:"samplings" json:"samplings,omitempty"`
	Spottings           []taxonomy.Taxon     `gel:"spottings" json:"spottings,omitempty"`
	Meta                people.Meta          `gel:"meta" json:"meta"`
}

func (e *Event) AddSampling(db geltypes.Executor, sampling SamplingInput) error {
	created, err := sampling.Save(db, e.ID)
	if err != nil {
		return err
	}
	e.Samplings = append(e.Samplings, created)
	return nil
}

func (e *Event) AddSpottings(db geltypes.Executor, taxa SpottingUpdate) error {
	spottings, err := taxa.Save(db, e.ID)
	if err != nil {
		return err
	}
	e.Spottings = spottings
	return nil
}

// AddAbioticMeasurement adds an abiotic measurement to the event.
// If a value for a given parameter already exists, it will be overwritten.
func (e *Event) AddAbioticMeasurement(db geltypes.Executor, measurements AbioticMeasurementInput) error {
	created, err := measurements.Save(db, e.ID)
	if err != nil {
		return err
	}
	e.AbioticMeasurements = append(e.AbioticMeasurements, created)
	return nil
}

var listEventsQuery = `#edgeql
	select events::Event {
		id,
		site: {name, code},
		programs: { * },
		performed_by: { * },
		performed_on,
		comments,
		spottings: { * },
		abiotic_measurements: { *, param: { * }  },
		samplings: {
			*,
			target_taxa: { * },
			fixatives: { * },
			methods: { * },
			habitats: { * }
		},
		meta: { * }
	}
`

func ListEvents(db geltypes.Executor) ([]Event, error) {
	var items = []Event{}
	err := db.Query(context.Background(),
		listEventsQuery+" order by .performed_on.date desc",
		&items)
	return items, err
}

func ListSiteEvents(e geltypes.Executor, siteCode string) ([]Event, error) {
	var items = []Event{}
	query := fmt.Sprintf(`#edgeql
		%s
		filter .site = assert_exists((select location::Site filter .code = <str>$0))
		order by .performed_on.date desc
	`, listEventsQuery)
	err := e.Query(context.Background(), query, &items, siteCode)
	// Capture assert exists error and return a NoDataError
	if db.IsCardinalityViolation(err) {
		err = db.NewNoDataError("Site not found")
	}
	return items, err
}

type EventInput struct {
	PerformedBy []string                       `json:"performed_by" minLength:"1"`
	PerformedOn DateWithPrecisionInput         `json:"performed_on"`
	Programs    models.OptionalInput[[]string] `json:"programs,omitempty"`
}

func (i EventInput) Save(e geltypes.Executor, site_code string) (created Event, err error) {

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
				spottings: { * },
				abiotic_measurements: { *, param: { * }  },
				samplings: { *, target_taxa: { * }, fixatives: { * }, methods: { * }, habitats: { * } },
				meta: { * }
			}
		`, &created, site_code, data)
	return
}

type EventUpdate struct {
	PerformedBy models.OptionalInput[[]string]               `gel:"performed_by" json:"performed_by,omitempty"`
	PerformedOn models.OptionalInput[DateWithPrecisionInput] `gel:"performed_on" json:"performed_on,omitempty"`
	Programs    models.OptionalNull[[]string]                `gel:"programs" json:"programs,omitempty"`
	Comments    models.OptionalNull[string]                  `gel:"comments" json:"comments,omitempty"`
	Spottings   models.OptionalNull[[]string]                `gel:"spottings" json:"spottings,omitempty"`
}

func (u EventUpdate) Save(e geltypes.Executor, id geltypes.UUID) (updated Event, err error) {
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
				spottings: { * },
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
			"comments": `data['comments']`,
			"spottings": `#edgeql
				(
					select taxonomy::Taxon
					filter .name in <str>json_array_unpack(data['spottings'])
				)`,
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, id, data)
	return
}

func DeleteEvent(db geltypes.Executor, id geltypes.UUID) (deleted Event, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
				delete events::Event filter .id = <uuid>$0
			) {
				site: {name, code},
				programs: { * },
				performed_by: { * },
				spottings: { * },
				abiotic_measurements: { *, param: { * }  },
				samplings: { *, target_taxa: { * }, fixatives: { * }, methods: { * }, habitats: { * } }
			};
		`,
		&deleted, id)
	return
}
