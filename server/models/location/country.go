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

type CountrySummary struct {
	Country          `json:",inline" gel:"$inline"`
	SitesCount       int64 `json:"sites_count" gel:"sites_count"`
	OccurrencesCount int64 `json:"occurrences_count" gel:"occurrences_count"`
}

func CountriesSummary(db geltypes.Executor) ([]CountrySummary, error) {
	var res = []CountrySummary{}
	err := db.Query(context.Background(),
		`#edgeql
			select location::Country {
				*,
				sites_count := count(.sites),
				occurrences_count := count(.sites.samplings.occurrences)
				}
			order by .name asc
		`, &res)
	return res, err
}
