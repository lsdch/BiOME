package country

import (
	"context"
	"darco/proto/models"
	"embed"
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
