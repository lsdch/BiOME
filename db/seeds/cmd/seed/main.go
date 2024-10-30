package main

import (
	"context"
	"darco/proto/db"
	"darco/proto/models/location"
	"flag"
	"seeds"
	"seeds/email"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

var entities = []string{
	"institutions",
	"persons",
	"users",
	"programs",
	"sampling_methods",
	"fixatives",
	"abiotic",
}

func main() {

	database := flag.String("db", "", "The name of the database to seed")
	flag.Parse()

	client := db.Connect(edgedb.Options{Database: *database})

	email.SetupEmailConfig(client, email.EmailSetupArgs{})

	logrus.Infof("Seeding habitats")
	if err := location.InitialHabitatsSetup(client); err != nil {
		logrus.Fatalf("Failed to seed habitats: %v", err)
	}

	if err := seeds.SeedTaxonomyGBIF(client); err != nil {
		logrus.Errorf("Failed to load Asellidae taxonomy: %v", err)
		return
	}

	err := client.Tx(context.Background(),
		func(ctx context.Context, tx *edgedb.Tx) error {
			for _, entity := range entities {
				logrus.Infof("Seeding %s", entity)
				err := seeds.Seed(tx, entity)
				if err != nil {
					return err
				}
			}
			return nil
		})
	if err != nil {
		logrus.Errorf("Seeding failed: %v", err)
	}
}
