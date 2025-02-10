package occurrence

import (
	"context"
	"encoding/json"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/location"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/validations"

	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
)

type LatLongCoords struct {
	Latitude  float32 `edgedb:"latitude" json:"latitude" minimum:"-90" maximum:"90" example:"39.1137"`
	Longitude float32 `edgedb:"longitude" json:"longitude" minimum:"-180" maximum:"180" example:"9.5064"`
}

func (c LatLongCoords) LatLong() (float32, float32) {
	return c.Latitude, c.Longitude
}

type Coordinates struct {
	Precision     location.CoordinatesPrecision `edgedb:"precision" json:"precision" doc:"Where the coordinates point to"`
	LatLongCoords `edgedb:"$inline" json:",inline"`
}

type SiteInput struct {
	Name                string                       `json:"name" minLength:"4"`
	Code                string                       `json:"code" pattern:"[A-Z0-9]+" patternDescription:"alphanum" minLength:"4" maxLength:"10" example:"SITE89" doc:"A short unique uppercase alphanumeric identifier"`
	Description         models.OptionalInput[string] `json:"description,omitempty"`
	Coordinates         Coordinates                  `json:"coordinates" doc:"Site coordinates in decimal degrees"`
	Altitude            models.OptionalInput[int32]  `json:"altitude,omitempty" doc:"Site altitude in meters"`
	Locality            models.OptionalInput[string] `json:"locality,omitempty" doc:"Nearest populated place"`
	UserDefinedLocality bool                         `edgedb:"user_defined_locality" json:"user_defined_locality" doc:"Signals if locality was manually entered by user, and automatically inferred from coordinates"`
	CountryCode         models.OptionalInput[string] `json:"country_code,omitempty" format:"country-code" pattern:"[A-Z]{2}" example:"FR"`
}

func (c SiteInput) LatLong() (float32, float32) {
	return c.Coordinates.LatLong()
}

// Validate checks if the site code is unique
func (i *SiteInput) Validate(edb edgedb.Executor) validations.ValidationErrors {
	codeChecker := db.DBProperty{Object: "location::Site", Property: "code"}
	var errs validations.ValidationErrors
	if _, codeExists := codeChecker.Exists(edb, i.Code); codeExists {
		errs = append(errs, &huma.ErrorDetail{Message: "Code already exists", Value: i.Code, Location: "code"})
	}
	return errs
}

type SiteItem struct {
	ID                  edgedb.UUID                       `edgedb:"id" json:"id" format:"uuid"`
	Name                string                            `edgedb:"name" json:"name" minLength:"4"`
	Code                string                            `edgedb:"code" json:"code" minLength:"4" maxLength:"8"`
	Description         edgedb.OptionalStr                `edgedb:"description" json:"description,omitempty"`
	Coordinates         Coordinates                       `edgedb:"coordinates" json:"coordinates"`
	Altitude            edgedb.OptionalInt32              `edgedb:"altitude" json:"altitude,omitempty"`
	Locality            edgedb.OptionalStr                `edgedb:"locality" json:"locality,omitempty"`
	Country             models.Optional[location.Country] `edgedb:"country" json:"country,omitempty"`
	AccessPoint         edgedb.OptionalStr                `edgedb:"access_point" json:"access_point,omitempty"`
	UserDefinedLocality bool                              `edgedb:"user_defined_locality" json:"user_defined_locality"`
}

type Site struct {
	SiteItem `edgedb:"$inline" json:",inline"`
	Datasets []DatasetInner `edgedb:"datasets" json:"datasets,omitempty"`
	Events   []Event        `edgedb:"events" json:"events,omitempty"`
	Meta     people.Meta    `edgedb:"meta" json:"meta"`
}

func ListSites(db edgedb.Executor) ([]Site, error) {
	var sites []Site
	err := db.Query(context.Background(),
		`#edgeql
			select location::Site {
				*,
				datasets: { * },
				meta: { * },
				country: { * },
				events: { * } order by .performed_on.date desc
			}
		`, &sites)
	return sites, err
}

func GetSite(db edgedb.Executor, identifier string) (Site, error) {
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

func (i SiteInput) Save(db edgedb.Executor) (*Site, error) {
	var created Site
	data, _ := json.Marshal(i)
	query := `#edgeql
		with module location, data := <json>$0
		select (insert_site(data)) { *, country: { * }, meta: { * }, datasets: { * } }
	`
	err := db.QuerySingle(context.Background(), query, &created, data)
	return &created, err
}

type SiteUpdate struct {
	Name                models.OptionalInput[string]      `edgedb:"name" json:"name,omitempty" minLength:"4"`
	Code                models.OptionalInput[string]      `edgedb:"code" json:"code,omitempty" pattern:"[A-Z0-9]+" patternDescription:"alphanum" minLength:"4" maxLength:"10" example:"SITE89" doc:"A short unique uppercase alphanumeric identifier"`
	Description         models.OptionalNull[string]       `edgedb:"description" json:"description,omitempty"`
	Coordinates         models.OptionalInput[Coordinates] `edgedb:"coordinates" json:"coordinates,omitempty" doc:"Site coordinates in decimal degrees"`
	Altitude            models.OptionalNull[int32]        `edgedb:"altitude" json:"altitude,omitempty" doc:"Site altitude in meters"`
	Locality            models.OptionalNull[string]       `edgedb:"locality" json:"locality,omitempty" doc:"Nearest populated place"`
	UserDefinedLocality models.OptionalInput[bool]        `edgedb:"user_defined_locality" json:"user_defined_locality" doc:"Signals whether locality was manually entered by user, and automatically inferred from coordinates"`
	CountryCode         models.OptionalNull[string]       `edgedb:"country" json:"country_code,omitempty" format:"country-code" pattern:"[A-Z]{2}" example:"FR"`
}

func (u SiteUpdate) Save(e edgedb.Executor, code string) (updated Site, err error) {
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
