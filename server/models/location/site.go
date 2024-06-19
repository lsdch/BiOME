package location

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/people"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
)

type Coordinates struct {
	Precision CoordinatePrecision `edgedb:"precision" json:"precision"`
	Latitude  float32             `edgedb:"latitude" json:"latitude" minimum:"-90" maximum:"90"`
	Longitude float32             `edgedb:"longitude" json:"longitude" minimum:"-180" maximum:"180"`
}

type SiteInput struct {
	Name         string                       `json:"name" minLength:"4"`
	Code         string                       `json:"code" minLength:"4" maxLength:"8"`
	Description  models.OptionalInput[string] `json:"description,omitempty"`
	Coordinates  Coordinates                  `json:"coordinates"`
	Altitude     models.OptionalInput[int32]  `json:"altitude,omitempty"`
	Region       models.OptionalInput[string] `json:"region,omitempty"`
	Municipality models.OptionalInput[string] `json:"municipality,omitempty"`
	CountryCode  string                       `json:"country_code"`
}

type SiteItem struct {
	ID           edgedb.UUID          `edgedb:"id" json:"id" format:"uuid"`
	Name         string               `edgedb:"name" json:"name" minLength:"4"`
	Code         string               `edgedb:"code" json:"code" minLength:"4" maxLength:"8"`
	Description  edgedb.OptionalStr   `edgedb:"description" json:"description"`
	Coordinates  Coordinates          `edgedb:"coordinates" json:"coordinates,omitempty"`
	Altitude     edgedb.OptionalInt32 `edgedb:"altitude" json:"altitude,omitempty"`
	Region       edgedb.OptionalStr   `edgedb:"region" json:"region,omitempty"`
	Municipality edgedb.OptionalStr   `edgedb:"municipality" json:"municipality,omitempty"`
	Country      Country              `edgedb:"country" json:"country"`
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
				precision := <CoordinateMaxPrecision>coords['precision'],
				latitude := <float32>coords['latitude'],
				longitude := <float32>coords['longitude']
			),
			region := <str>json_get(data, 'region'),
			municipality := <str>json_get(data, 'municipality'),
			country := (assert_exists(select Country filter .code = <str>data['country_code'])),
			altitude := <int32>json_get(data, 'altitude')
		}) { **, country: { * } }`,
		&created, data)
	return &created, err
}
