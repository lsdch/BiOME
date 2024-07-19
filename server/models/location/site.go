package location

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/people"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
)

type Coordinates struct {
	Precision CoordinatesPrecision `edgedb:"precision" json:"precision" doc:"Where the coordinates point to"`
	Latitude  float32              `edgedb:"latitude" json:"latitude" minimum:"-90" maximum:"90" example:"39.1137"`
	Longitude float32              `edgedb:"longitude" json:"longitude" minimum:"-180" maximum:"180" example:"9.5064"`
}

type SiteInput struct {
	Name        string                       `json:"name" minLength:"4"`
	Code        string                       `json:"code" pattern:"[A-Z0-9]+" patternDescription:"alphanum" minLength:"4" maxLength:"10" example:"SITE89" doc:"A short unique uppercase alphanumeric code to identify the site"`
	Description models.OptionalInput[string] `json:"description,omitempty"`
	Coordinates Coordinates                  `json:"coordinates" doc:"Site coordinates in decimal degrees"`
	Altitude    models.OptionalInput[int32]  `json:"altitude,omitempty" doc:"Site altitude in meters"`
	Locality    models.OptionalInput[string] `json:"locality,omitempty" doc:"Nearest populated place"`
	CountryCode string                       `json:"country_code" format:"country-code" pattern:"[A-Z]{2}" example:"FR"`
	AccessPoint models.OptionalInput[string] `json:"access_point,omitempty"`
}

type SiteItem struct {
	ID          edgedb.UUID          `edgedb:"id" json:"id" format:"uuid"`
	Name        string               `edgedb:"name" json:"name" minLength:"4"`
	Code        string               `edgedb:"code" json:"code" minLength:"4" maxLength:"8"`
	Description edgedb.OptionalStr   `edgedb:"description" json:"description"`
	Coordinates Coordinates          `edgedb:"coordinates" json:"coordinates"`
	Altitude    edgedb.OptionalInt32 `edgedb:"altitude" json:"altitude,omitempty"`
	Locality    edgedb.OptionalStr   `edgedb:"locality" json:"locality,omitempty"`
	Country     Country              `edgedb:"country" json:"country"`
	AccessPoint edgedb.OptionalStr   `edgedb:"access_point" json:"access_point,omitempty"`
}

type Site struct {
	SiteItem `edgedb:"$inline" json:",inline"`
	Datasets []SiteDatasetInner `edgedb:"datasets" json:"datasets"`
	Meta     people.Meta        `edgedb:"meta" json:"meta"`
}

func ListSites(db edgedb.Executor) ([]Site, error) {
	var sites []Site
	err := db.Query(context.Background(), `select location::Site { ** }`, &sites)
	return sites, err
}

func GetSite(db edgedb.Executor, identifier string) (Site, error) {
	var site Site
	err := db.QuerySingle(context.Background(), `select location::Site { ** } filter .code = $0`, &site, identifier)
	return site, err
}

func (i *SiteInput) Create(db edgedb.Executor) (*Site, error) {
	var created Site
	data, _ := json.Marshal(i)
	err := db.QuerySingle(context.Background(),
		`with module location,
			data := <json>$0,
			coords := data['coordinates'],
		select ( insert Site {
			name := <str>data['name'],
			code := <str>data['code'],
			description := <str>json_get(data, 'description'),
			coordinates := (
				precision := <CoordinatesPrecision>coords['precision'],
				latitude := <float32>coords['latitude'],
				longitude := <float32>coords['longitude']
			),
			locality := <str>json_get(data, 'locality'),
			country := (assert_exists(select Country filter .code = <str>data['country_code'])),
			altitude := <int32>json_get(data, 'altitude')
		}) { **, country: { * } }`,
		&created, data)
	return &created, err
}

func ListAccessPoints(db edgedb.Executor) ([]string, error) {
	var accessPoints []string
	err := db.Query(context.Background(),
		`select distinct location::Site.access_point order by .`,
		&accessPoints,
	)
	return accessPoints, err
}
