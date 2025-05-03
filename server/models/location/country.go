package location

import (
	"context"
	_ "embed"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/geltypes"
)

//go:embed data/countries.json
var seed []byte

func SetupCountries(db *gel.Client) error {
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
	ID           geltypes.UUID `json:"id" gel:"id" format:"uuid"`
	Name         string        `json:"name" gel:"name" example:"Germany" binding:"required"`
	Code         string        `json:"code" gel:"code" example:"DE"`
	Continent    string        `json:"continent" gel:"continent" example:"Europe"`
	Subcontinent string        `json:"subcontinent" gel:"subcontinent" example:"Western Europe"`
}

func ListCountries(db geltypes.Executor) ([]Country, error) {
	var countries []Country
	err := db.Query(context.Background(),
		`#edgeql
			select location::Country { * } order by .name;
		`,
		&countries)
	return countries, err
}

type CountryWithSitesCount struct {
	Country    `json:",inline" gel:"$inline"`
	SitesCount int64 `json:"sites_count" gel:"sites_count"`
}

func SitesCountByCountry(db geltypes.Executor) ([]CountryWithSitesCount, error) {
	var res []CountryWithSitesCount
	err := db.Query(context.Background(),
		`#edgeql
			select location::Country { *, sites_count := count(.sites) }
			order by .sites_count desc
		`, &res)
	return res, err
}
