package location

import (
	"context"
	_ "embed"

	"github.com/edgedb/edgedb-go"
)

//go:embed data/countries.json
var seed []byte

func SetupCountries(db *edgedb.Client) error {
	return db.Execute(context.Background(),
		`#edgeql
		with module location,
		data := <json>$0
		for item in json_array_unpack(data) union (
			insert Country {
				name := <str>item['name'],
				code := <str>item['code']
			}
			unless conflict on (.code) else (
				update Country set {
					name := <str>item['name']
				}
			)
		)`, seed)
}

type Country struct {
	ID   edgedb.UUID `json:"id" edgedb:"id" format:"uuid"`
	Name string      `json:"name" edgedb:"name" example:"Germany" binding:"required"`
	Code string      `json:"code" edgedb:"code" example:"DE" binding:"required,country_code=iso3166_1_alpha2"`
}

func ListCountries(db edgedb.Executor) ([]Country, error) {
	var countries []Country
	err := db.Query(context.Background(),
		`#edgeql
			select location::Country { * } order by .name;
		`,
		&countries)
	return countries, err
}

type CountryWithSitesCount struct {
	Country    `json:",inline" edgedb:"$inline"`
	SitesCount int64 `json:"sites_count" edgedb:"sites_count"`
}

func SitesCountByCountry(db edgedb.Executor) ([]CountryWithSitesCount, error) {
	var res []CountryWithSitesCount
	err := db.Query(context.Background(),
		`#edgeql
			select location::Country { *, sites_count := count(.sites) }
		`, &res)
	return res, err
}
