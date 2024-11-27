package occurrence

import (
	"bytes"
	"context"
	"darco/proto/db"
	"darco/proto/models"
	"darco/proto/models/location"
	"darco/proto/models/people"
	"darco/proto/models/validations"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type Coordinates struct {
	Precision location.CoordinatesPrecision `edgedb:"precision" json:"precision" doc:"Where the coordinates point to"`
	Latitude  float32                       `edgedb:"latitude" json:"latitude" minimum:"-90" maximum:"90" example:"39.1137"`
	Longitude float32                       `edgedb:"longitude" json:"longitude" minimum:"-180" maximum:"180" example:"9.5064"`
}

type SiteInput struct {
	Name        string                       `json:"name" minLength:"4"`
	Code        string                       `json:"code" pattern:"[A-Z0-9]+" patternDescription:"alphanum" minLength:"4" maxLength:"10" example:"SITE89" doc:"A short unique uppercase alphanumeric identifier"`
	Description models.OptionalInput[string] `json:"description,omitempty"`
	Coordinates Coordinates                  `json:"coordinates" doc:"Site coordinates in decimal degrees"`
	Altitude    models.OptionalInput[int32]  `json:"altitude,omitempty" doc:"Site altitude in meters"`
	Locality    models.OptionalInput[string] `json:"locality,omitempty" doc:"Nearest populated place"`
	CountryCode string                       `json:"country_code" format:"country-code" pattern:"[A-Z]{2}" example:"FR"`
}

func (i *SiteInput) Validate(edb edgedb.Executor) validations.ValidationErrors {
	codeChecker := db.DBProperty{Object: "location::Site", Property: "code"}
	nameChecker := db.DBProperty{Object: "location::Site", Property: "name"}
	var errs validations.ValidationErrors
	if _, codeExists := codeChecker.Exists(edb, i.Code); codeExists {
		errs = append(errs, &huma.ErrorDetail{Message: "Code already exists", Value: i.Code, Location: "code"})
	}
	if _, nameExists := nameChecker.Exists(edb, i.Name); nameExists {
		errs = append(errs, &huma.ErrorDetail{Message: "Name already exists", Value: i.Name, Location: "name"})
	}
	return errs
}

type SiteItem struct {
	ID          edgedb.UUID          `edgedb:"id" json:"id" format:"uuid"`
	Name        string               `edgedb:"name" json:"name" minLength:"4"`
	Code        string               `edgedb:"code" json:"code" minLength:"4" maxLength:"8"`
	Description edgedb.OptionalStr   `edgedb:"description" json:"description"`
	Coordinates Coordinates          `edgedb:"coordinates" json:"coordinates"`
	Altitude    edgedb.OptionalInt32 `edgedb:"altitude" json:"altitude,omitempty"`
	Locality    edgedb.OptionalStr   `edgedb:"locality" json:"locality,omitempty"`
	Country     location.Country     `edgedb:"country" json:"country"`
	AccessPoint edgedb.OptionalStr   `edgedb:"access_point" json:"access_point,omitempty"`
}

type Site struct {
	SiteItem `edgedb:"$inline" json:",inline"`
	Datasets []DatasetInner `edgedb:"datasets" json:"datasets"`
	Events   []Event        `edgedb:"events" json:"events"`
	Meta     people.Meta    `edgedb:"meta" json:"meta"`
}

func ListSites(db edgedb.Executor) ([]Site, error) {
	var sites []Site
	err := db.Query(context.Background(),
		`select location::Site {
			*,
			datasets: { * },
			meta: { * },
			country: { * },
			events: { * } order by .performed_on.date desc
		}`, &sites)
	return sites, err
}

func GetSite(db edgedb.Executor, identifier string) (Site, error) {
	var site Site
	err := db.QuerySingle(context.Background(),
		`select location::Site {
			*,
			country: { * },
			datasets: { * },
			meta: { * },
			events: { *,
				site: {name, code},
				programs: { * },
				performed_by: { * },
				spotting: { *, target_taxa: { * } },
				abiotic_measurements: { *, param: { * }  },
				samplings: { *, target_taxa: { * }, fixatives: { * }, methods: { * }, habitats: { * } }
			} order by .performed_on.date desc
		} filter .code = <str>$0`,
		&site, identifier)
	return site, err
}

var siteInsertQueryTmpl = template.Must(
	template.New("siteInsertQuery").
		Parse(`insert location::Site {{ "{" }}
			name := <str>{{.Json}}['name'],
			code := <str>{{.Json}}['code'],
			description := <str>json_get({{.Json}}, 'description'),
			coordinates := (
				precision := <location::CoordinatesPrecision>{{.Json}}['coordinates']['precision'],
				latitude := <float32>{{.Json}}['coordinates']['latitude'],
				longitude := <float32>{{.Json}}['coordinates']['longitude']
			),
			locality := <str>json_get({{.Json}}, 'locality'),
			country := (
				select assert_exists(location::Country
				filter .code = <str>{{.Json}}['country_code'])
			),
			altitude := <int32>json_get({{.Json}}, 'altitude')
		{{ "}" }}
	`))

func (i *SiteInput) InsertQuery(jsonVar string) string {
	var query bytes.Buffer
	_ = siteInsertQueryTmpl.Execute(&query, struct{ Json string }{jsonVar})
	logrus.Infof("%s", query.String())
	return query.String()
}

func (i *SiteInput) Create(db edgedb.Executor) (*Site, error) {
	var created Site
	data, _ := json.Marshal(i)
	query := fmt.Sprintf(
		`with module location,
		data := <json>$0,
		coords := data['coordinates'],
	select ( %s ) { *, country: { * }, meta: { * }, datasets: { * } }`,
		i.InsertQuery("data"))
	err := db.QuerySingle(context.Background(), query, &created, data)
	return &created, err
}

type SiteUpdate struct {
	Name        models.OptionalInput[string]      `json:"name" minLength:"4"`
	Code        models.OptionalInput[string]      `json:"code" pattern:"[A-Z0-9]+" patternDescription:"alphanum" minLength:"4" maxLength:"10" example:"SITE89" doc:"A short unique uppercase alphanumeric identifier"`
	Description models.OptionalNull[string]       `json:"description,omitempty"`
	Coordinates models.OptionalInput[Coordinates] `json:"coordinates" doc:"Site coordinates in decimal degrees"`
	Altitude    models.OptionalNull[int32]        `json:"altitude,omitempty" doc:"Site altitude in meters"`
	Locality    models.OptionalNull[string]       `json:"locality,omitempty" doc:"Nearest populated place"`
	CountryCode models.OptionalInput[string]      `json:"country_code" format:"country-code" pattern:"[A-Z]{2}" example:"FR"`
}

func (u SiteUpdate) Update(e edgedb.Executor, code string) (updated Site, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `with item := <json>$1,
		select (update location::Site filter .code = <str>$0 set {
			%s
		}) { *, datasets: { * }, meta: { * }, country: { * } }`,
		Mappings: map[string]string{
			"name":        "<str>item['name']",
			"code":        "<str>item['code']",
			"description": "<str>item['description']",
			"coordinates": `(
				precision := <location::CoordinatesPrecision>item['coordinates']['precision'],
				latitude := <float32>item['coordinates']['latitude'],
				longitude := <float32>item['coordinates']['longitude']
			)`,
			"altitude": "<int32>item['altitude']",
			"locality": "<str>item['locality']",
			"country": `(
				select assert_exists(location::Country
				filter .code = <str>{{.Json}}['country_code'])
			)`,
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, code, data)
	updated.Meta.Update(e)
	return
}
