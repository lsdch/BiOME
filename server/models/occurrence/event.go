package occurrence

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/people"
	"encoding/json"
	"time"

	"github.com/edgedb/edgedb-go"
)

type DateWithPrecision struct {
	Date      time.Time     `edgedb:"date" json:"date"`
	Precision DatePrecision `edgedb:"precision" json:"precision"`
}

type Event struct {
	ID                  edgedb.UUID               `edgedb:"id" json:"id" format:"uuid"`
	SiteCode            string                    `edgedb:"site_code" json:"site_code"`
	PerformedBy         []people.PersonUser       `edgedb:"performed_by" json:"performed_by" minLength:"1"`
	PerformedOn         DateWithPrecision         `edgedb:"performed_on" json:"performed_on"`
	Programs            []ProgramInner            `edgedb:"programs" json:"programs,omitempty"`
	AbioticMeasurements []AbioticMeasurement      `edgedb:"abiotic_measurements" json:"abiotic_measurements"`
	Samplings           []Sampling                `edgedb:"samplings" json:"samplings"`
	Spotting            models.Optional[Spotting] `edgedb:"spotting" json:"spotting,omitempty"`
	Meta                people.Meta               `edgedb:"meta" json:"meta"`
}

func ListEvents(db edgedb.Executor) ([]Event, error) {
	var items = []Event{}
	err := db.Query(context.Background(),
		`select events::Event { **, site_code := .site.code };`,
		&items)
	return items, err
}

type EventInput struct {
	SiteCode    string                         `json:"site_code"`
	PerformedBy []string                       `json:"performed_by" minLength:"1"`
	PerformedOn DateWithPrecision              `json:"performed_on"`
	Programs    models.OptionalInput[[]string] `json:"programs,omitempty"`
}

func (i EventInput) Create(e edgedb.Executor) (created Event, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
		with data = <json>$0,
		select (insert events::Event {
			site := (
				select location::Site filter .code = <str>data['site_code']
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
		})`, &created, data)
	return
}
