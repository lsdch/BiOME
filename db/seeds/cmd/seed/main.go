package main

import (
	"context"
	"darco/proto/db"
	"darco/proto/models/occurrence"
	"flag"
	"seeds"
	"seeds/email"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

var entities = []string{
	"countries",
	"institutions",
	"persons",
	"users",
	"articles",
	"programs",
	"sampling_methods",
	"fixatives",
	"abiotic",
	"genes",
	"seq_databases",
	"datasets",
}

func main() {

	database := flag.String("db", "", "The name of the database to seed")
	flag.Parse()

	client := db.Connect(edgedb.Options{Database: *database})

	err := client.Tx(context.Background(), func(ctx context.Context, tx *edgedb.Tx) error {
		if err := email.SetupEmailConfig(client, email.EmailSetupArgs{}); err != nil {
			return err
		}

		logrus.Infof("Seeding habitats")
		if err := occurrence.InitialHabitatsSetup(tx); err != nil {
			logrus.Fatalf("Failed to seed habitats: %v", err)
		}

		if err := seeds.SeedTaxonomyGBIF(tx); err != nil {
			logrus.Errorf("Failed to load Asellidae taxonomy: %v", err)
			return err
		}

		logrus.Infof("Seeding...")
		for _, entity := range entities {
			logrus.Infof("• %s", entity)
			err := seeds.Seed(tx, entity)
			if err != nil {
				return err
			}
		}

		logrus.Infof("Postprocessing...")
		logrus.Infof("• generate sequence codes")
		if err := tx.Execute(context.Background(),
			`#edgeql
				update seq::ExternalSequence set {};
			`); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logrus.Errorf("Seeding failed: %v", err)
	}
}
