package location

import (
	"context"
	_ "embed"

	"github.com/edgedb/edgedb-go"
)

//go:embed queries/setup_countries.edgeql
var setupCountriesCmd string

//go:embed data/countries.json
var seed []byte

func SetupCountries(db *edgedb.Client) error {
	return db.Execute(context.Background(), setupCountriesCmd, seed)
}

type Country struct {
	ID   edgedb.UUID `json:"id" edgedb:"id" format:"uuid"`
	Name string      `json:"name" edgedb:"name" example:"Germany" binding:"required"`
	Code string      `json:"code" edgedb:"code" example:"DE" binding:"required,country_code=iso3166_1_alpha2"`
}

func ListCountries(db edgedb.Executor) ([]Country, error) {
	var countries []Country
	query := `select location::Country { * };`
	err := db.Query(context.Background(), query, &countries)
	return countries, err
}
