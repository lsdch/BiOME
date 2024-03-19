package country

import (
	"context"
	"embed"

	"github.com/edgedb/edgedb-go"
)

//go:embed queries/setup_countries.edgeql
var setupCountriesCmd string

//go:embed countries.json
var seed embed.FS

func Setup(db *edgedb.Client) error {
	json, err := seed.ReadFile("countries.json")
	if err != nil {
		return err
	}
	err = db.Execute(context.Background(), setupCountriesCmd, json)
	return err
}

type Country struct {
	ID           edgedb.UUID `json:"id" edgedb:"id" example:"<UUID>"`
	Name         string      `json:"name" edgedb:"name" example:"Germany" binding:"required"`
	Code         string      `json:"code" edgedb:"code" example:"DE" binding:"required,country_code=iso3166_1_alpha2"`
	NbLocalities int64       `json:"nbLocalities" edgedb:"nb_localities" example:"9"`
} // @name Country

func List(db *edgedb.Client) (countries []Country, err error) {
	query := `select
		location::Country {
			*, nb_localities := count(.localities)
		}
		order by (exists .localities) desc then .name asc;`
	err = db.Query(context.Background(), query, &countries)
	return
}
