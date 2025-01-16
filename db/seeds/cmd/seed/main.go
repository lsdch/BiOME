package main

import (
	"context"
	"darco/proto/db"
	"darco/proto/models"
	"darco/proto/models/occurrence"
	"darco/proto/models/people"
	"darco/proto/services/setup"
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

var superAdmin = people.SuperAdminInput{
	UserInput: people.UserInput{
		Login:         "lsdch",
		EmailField:    people.EmailField{Email: "louis.duchemin@univ-lyon1.fr"},
		PasswordInput: people.PasswordInput{Password: "superadmin", ConfirmPwd: "superadmin"},
	},
	PersonIdentity: people.PersonIdentity{
		FirstName: "Louis",
		LastName:  "Duchemin",
	},
	Alias: models.OptionalInput[string]{
		IsSet: true,
		Value: "lsdch",
	},
	Institution: people.InstitutionInput{
		InstitutionInfos: people.InstitutionInfos{
			Name: "Laboratoire d'Écologie des Hydrosystèmes Naturels et Anthropisés",
			Code: "LEHNA",
			Kind: "Lab",
		},
	},
}

func main() {

	database := flag.String("db", "", "The name of the database to seed")
	flag.Parse()

	wadSites, err := seeds.LoadSites("data/Aselloidea/sites.json", 500)
	if err != nil {
		logrus.Fatalf("Failed to load WAD sites: %v", err)
	}

	client := db.Connect(edgedb.Options{Database: *database})

	err = client.Tx(context.Background(), func(ctx context.Context, tx *edgedb.Tx) error {

		logrus.Infof("⚙ Initializing settings with superadmin account")
		if err := setup.InitTx(tx, superAdmin); err != nil {
			return err
		}

		if err := email.SetupEmailConfig(client, email.EmailSetupArgs{}); err != nil {
			return err
		}

		logrus.Infof("🌱 Seeding habitats")
		if err := occurrence.InitialHabitatsSetup(tx); err != nil {
			logrus.Errorf("Failed to seed habitats: %v", err)
			return err
		}

		if err := seeds.SeedTaxonomyGBIF(tx); err != nil {
			logrus.Errorf("Failed to load Asellidae taxonomy: %v", err)
			return err
		}

		logrus.Infof("🌱 Seeding...")
		for _, entity := range entities {
			logrus.Infof("• %s", entity)
			err := seeds.Seed(tx, entity)
			if err != nil {
				return err
			}
		}

		logrus.Infof("⚙ Postprocessing...")
		logrus.Infof("• generate sequence codes")
		if err := tx.Execute(context.Background(),
			`#edgeql
				update seq::ExternalSequence set {};
			`); err != nil {
			return err
		}

		logrus.Info("🧪 Empirical datasets")
		logrus.Infof("🌱 Seeding WAD sampling sites")
		if err := seeds.SeedSites(tx, *wadSites); err != nil {
			logrus.Errorf("Failed to seed WAD sampling sites")
			return err
		}

		return nil
	})

	if err != nil {
		logrus.Errorf("Seeding failed: %v", err)
	}
}
