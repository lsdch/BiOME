package main

import (
	"context"
	"darco/proto/db"
	"embed"
	"flag"
	"fmt"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

//go:embed queries
var queries embed.FS

//go:embed data
var data embed.FS

func entityQueryPath(entity string) string {
	return fmt.Sprintf("queries/%s.edgeql", entity)
}
func entityDataPath(entity string) string {
	return fmt.Sprintf("data/%s.json", entity)
}

func Seed(tx *edgedb.Tx, entity string) error {
	queryPath := entityQueryPath(entity)
	dataPath := entityDataPath(entity)
	query, err := queries.ReadFile(queryPath)
	if err != nil {
		logrus.Errorf("Failed to load seed query @ %s: %v", queryPath, err)
		return err
	}

	data, err := data.ReadFile(dataPath)
	if err != nil {
		logrus.Errorf("Failed to load seed data @ %s: %v", dataPath, err)
		return err
	}

	err = tx.Execute(context.Background(), string(query), data)
	if err != nil {
		logrus.Errorf(
			"Query execution failed for query @ %s and data @ %s:\n%v",
			queryPath, dataPath, err,
		)
		return err
	}
	return nil
}

var entities = []string{
	"institutions",
	"persons",
	"users",
}

func main() {

	database := flag.String("db", "", "The name of the database to seed")
	flag.Parse()

	client := db.Connect(edgedb.Options{Database: *database})

	setupEmailConfig(client)

	if err := seedTaxonomyGBIF(client); err != nil {
		logrus.Errorf("Failed to load Asellidae taxonomy: %v", err)
		return
	}

	client.Tx(context.Background(),
		func(ctx context.Context, tx *edgedb.Tx) error {
			for _, entity := range entities {
				logrus.Infof("Seeding %s", entity)
				err := Seed(tx, entity)
				if err != nil {
					return err
				}
			}
			return nil
		})
}
