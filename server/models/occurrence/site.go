package occurrence

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/dataset"
	"github.com/lsdch/biome/models/location"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/validations"

	"github.com/danielgtaylor/huma/v2"
)

type LatLongCoords struct {
	Latitude  float32 `gel:"latitude" json:"latitude" minimum:"-90" maximum:"90" example:"43.5684"`
	Longitude float32 `gel:"longitude" json:"longitude" minimum:"-180" maximum:"180" example:"3.5678"`
}

func (c LatLongCoords) LatLong() (float32, float32) {
	return c.Latitude, c.Longitude
}

func (c LatLongCoords) FindCountry(db geltypes.Executor) (country location.Country, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select location::position_to_country(<float32>$0, <float32>$1) { * }
		`,
		&country, c.Latitude, c.Longitude)
	return
}

func (c LatLongCoords) SitesProximity(db geltypes.Executor, distance float32) ([]SiteItem, error) {
	var sites []SiteItem
	err := db.Query(context.Background(),
		`#edgeql
			select location::sites_proximity(<float32>$0, <float32>$1, <float32>$2) { * }
		`,
		&sites, c.Latitude, c.Longitude, distance)
	return sites, err
}

type Coordinates struct {
	Precision     location.CoordinatesPrecision `gel:"precision" json:"precision" doc:"Where the coordinates point to"`
	LatLongCoords `gel:"$inline" json:",inline"`
}

type SiteInput struct {
	Name                string                       `json:"name" minLength:"4" doc:"A short descriptive name"`
	Code                string                       `json:"code" pattern:"[A-Z0-9]+" patternDescription:"alphanum" minLength:"4" maxLength:"10" example:"SITE89" doc:"A short unique uppercase alphanumeric identifier"`
	Description         models.OptionalInput[string] `json:"description,omitempty"`
	Coordinates         Coordinates                  `json:"coordinates" doc:"Site coordinates in decimal degrees"`
	Altitude            models.OptionalInput[int32]  `json:"altitude,omitempty" doc:"Site altitude in meters"`
	Locality            models.OptionalInput[string] `json:"locality,omitempty" doc:"Nearest populated place"`
	UserDefinedLocality bool                         `json:"user_defined_locality,omitempty" doc:"Signals if locality was manually entered by user, and automatically inferred from coordinates"`
	CountryCode         models.OptionalInput[string] `json:"country_code,omitempty" format:"country-code" pattern:"[A-Z]{3}" example:"FRA" doc:"ISO 3166-1 alpha-3 country code"`
}

func (c SiteInput) LatLong() (float32, float32) {
	return c.Coordinates.LatLong()
}

// Validate checks if the site code is unique
func (i *SiteInput) Validate(edb geltypes.Executor) validations.ValidationErrors {
	codeChecker := db.DBProperty{Object: "location::Site", Property: "code"}
	var errs validations.ValidationErrors
	if _, codeExists := codeChecker.Exists(edb, i.Code); codeExists {
		errs = append(errs, &huma.ErrorDetail{Message: "Code already exists", Value: i.Code, Location: "code"})
	}
	return errs
}

type SiteItem struct {
	ID                  geltypes.UUID                     `gel:"id" json:"id" format:"uuid"`
	Name                string                            `gel:"name" json:"name" minLength:"4"`
	Code                string                            `gel:"code" json:"code" minLength:"4" maxLength:"8"`
	Description         geltypes.OptionalStr              `gel:"description" json:"description,omitempty"`
	Coordinates         Coordinates                       `gel:"coordinates" json:"coordinates"`
	Altitude            geltypes.OptionalInt32            `gel:"altitude" json:"altitude,omitempty"`
	Locality            geltypes.OptionalStr              `gel:"locality" json:"locality,omitempty"`
	Country             models.Optional[location.Country] `gel:"country" json:"country,omitempty"`
	AccessPoint         geltypes.OptionalStr              `gel:"access_point" json:"access_point,omitempty"`
	UserDefinedLocality bool                              `gel:"user_defined_locality" json:"user_defined_locality"`
}

type Site struct {
	SiteItem `gel:"$inline" json:",inline"`
	Datasets []dataset.DatasetInner `gel:"datasets" json:"datasets,omitempty"`
	Events   []Event                `gel:"events" json:"events,omitempty"`
	Meta     people.Meta            `gel:"meta" json:"meta"`
}

type ListSitesOptions struct {
	Datasets  []string `json:"datasets,omitempty" query:"datasets"`
	Countries []string `json:"countries,omitempty" query:"countries"`
	// Sampled bool `json:"sampled,omitempty" query:"sampled"`
}

func (o ListSitesOptions) Options() ListSitesOptions {
	return o
}

func ListSites(db geltypes.Executor, options ListSitesOptions) ([]Site, error) {
	var sites []Site
	opts, _ := json.Marshal(options)
	err := db.Query(context.Background(),
		`#edgeql
			with opts := <json>$0,
      countries := <str>json_array_unpack(json_get(opts, 'countries'))
			select location::Site {
				*,
				datasets: { * },
				meta: { * },
				country: { * },
				events: { * } order by .performed_on.date desc
			}
			filter (not exists countries) or (.country.code in countries)
    `,
		&sites, opts)
	return sites, err
}

func GetSite(db geltypes.Executor, identifier string) (Site, error) {
	var site Site
	err := db.QuerySingle(context.Background(),
		`#edgeql
			select location::Site {
				*,
				country: { * },
				datasets: { * },
				meta: { * },
				events: { *,
					site: {name, code},
					programs: { * },
					performed_by: { * },
					spottings: { * },
					abiotic_measurements: { *, param: { * }  },
					samplings: {
						*,
						target_taxa: { * },
						fixatives: { * },
						methods: { * },
						habitats: { * },
						samples: {
							**,
							identification: { **, identified_by: { * } }
						},
						occurring_taxa: { * }
					},
					meta: { * }
				} order by .performed_on.date desc
			} filter .code = <str>$0
		`,
		&site, identifier)
	return site, err
}

func (i SiteInput) Save(db geltypes.Executor) (*Site, error) {
	var created Site
	data, _ := json.Marshal(i)
	query := `#edgeql
		with module location, data := <json>$0
		select insert_site(data) { *, country: { * }, meta: { * }, datasets: { * } }
	`
	err := db.QuerySingle(context.Background(), query, &created, data)
	return &created, err
}

type SiteUpdate struct {
	Name                models.OptionalInput[string]      `gel:"name" json:"name,omitempty" minLength:"4"`
	Code                models.OptionalInput[string]      `gel:"code" json:"code,omitempty" pattern:"[A-Z0-9]+" patternDescription:"alphanum" minLength:"4" maxLength:"10" example:"SITE89" doc:"A short unique uppercase alphanumeric identifier"`
	Description         models.OptionalNull[string]       `gel:"description" json:"description,omitempty"`
	Coordinates         models.OptionalInput[Coordinates] `gel:"coordinates" json:"coordinates,omitempty" doc:"Site coordinates in decimal degrees"`
	Altitude            models.OptionalNull[int32]        `gel:"altitude" json:"altitude,omitempty" doc:"Site altitude in meters"`
	Locality            models.OptionalNull[string]       `gel:"locality" json:"locality,omitempty" doc:"Nearest populated place"`
	UserDefinedLocality models.OptionalInput[bool]        `gel:"user_defined_locality" json:"user_defined_locality" doc:"Signals whether locality was manually entered by user, and automatically inferred from coordinates"`
	CountryCode         models.OptionalNull[string]       `gel:"country" json:"country_code,omitempty" format:"country-code" pattern:"[A-Z]{2}" example:"FR"`
}

func (u SiteUpdate) Save(e geltypes.Executor, code string) (updated Site, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$1,
			select (update location::Site filter .code = <str>$0 set {
				%s
			}) { *, datasets: { * }, meta: { * }, country: { * } }
		`,
		Mappings: map[string]string{
			"name":        "<str>item['name']",
			"code":        "<str>item['code']",
			"description": "<str>item['description']",
			"coordinates": `#edgeql
				location::coords_from_json(item['coordinates']),
			`,
			"altitude": "<int32>item['altitude']",
			"locality": "<str>item['locality']",
			"country": `#edgeql
				location::find_country(item['country_code'])
			`,
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, code, data)
	updated.Meta.Save(e)
	return
}
