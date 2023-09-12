package country

import (
	"context"
	"darco/proto/models"
	"embed"

	"github.com/edgedb/edgedb-go"
)

//go:embed setup_countries.edgeql
var setupCountriesCmd string

//go:embed countries.json
var seed embed.FS

func Setup() error {
	json, err := seed.ReadFile("countries.json")
	if err != nil {
		return err
	}
	err = models.DB.Execute(context.Background(), setupCountriesCmd, json)
	return err
}

type Country struct {
	ID           edgedb.UUID `json:"id" edgedb:"id" example:"<UUID>"`
	Name         string      `json:"name" edgedb:"name" example:"Germany"`
	Code         string      `json:"code" edgedb:"code" example:"DE"`
	NbLocalities int64       `json:"nbLocalities" edgedb:"nb_localities" example:"9"`
}

func List() (countries []Country, err error) {
	query := `select
		location::Country {
			id, name, code, nb_localities := count(.localities)
		}
		order by (exists .localities) desc then .name asc;`
	err = models.DB.Query(context.Background(), query, &countries)
	return
}
