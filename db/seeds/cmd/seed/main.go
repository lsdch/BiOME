package main

import (
	"context"
	"flag"
	"fmt"
	"seeds"
	"seeds/email"

	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/settings"

	"github.com/sirupsen/logrus"
)

var entities = []string{
	// "countries",
	"organisations",
	"persons",
	"users",
	"articles",
	"programs",
	"sampling_methods",
	"fixatives",
	"abiotic",
	"genes",
	"data_sources",
	// "datasets",
}

var superAdminInput = people.SuperAdminInput{
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
	Organisation: people.OrganisationInput{
		OrganisationInfos: people.OrganisationInfos{
			Name: "Laboratoire d'Écologie des Hydrosystèmes Naturels et Anthropisés",
			Code: "LEHNA",
			Kind: "Lab",
		},
	},
}

func main() {

	database := flag.String("db", "", "The name of the database to seed")
	flag.Parse()

	client := db.Connect(gelcfg.Options{Database: *database})

	// aselloidea, err := seeds.LoadSiteDataset(client, "data/Aselloidea/sites.json")
	// if err != nil {
	// 	logrus.Fatalf("Failed to load Asellidae sites: %v", err)
	// }

	timeout, _ := geltypes.ParseDuration("15m")
	err := client.WithConfig(map[string]interface{}{
		"session_idle_transaction_timeout": timeout,
	}).Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {

		logrus.Infof("🌱 Seeding countries")
		if err := seeds.SeedCountriesGeoJSON(tx, "../../data/remote/countries.json"); err != nil {
			return fmt.Errorf("Failed to seed countries: %v", err)
		}

		logrus.Infof("⚙ Initializing settings with superadmin account")
		superAdmin, err := superAdminInput.Save(tx)
		if err != nil {
			return fmt.Errorf("Failed to initialize super admin account: %v", err)
		}

		if err := (settings.SettingsInput{
			SuperAdminID: superAdmin.ID,
			Instance: settings.InstanceSettingsInput{
				InstanceSettingsInner: settings.InstanceSettingsInner{
					Name: "[BiOME prototype]",
				},
				Description: models.NewOptionalNull("Prototype BiOME instance"),
			},
		}).SaveTx(tx); err != nil {
			return fmt.Errorf("Failed to initialize settings: %v", err)
		}

		if err := email.SetupEmailConfig(client, email.EmailSetupArgs{}); err != nil {
			return err
		}

		logrus.Infof("🌱 Seeding habitats")
		if err := occurrence.InitialHabitatsSetup(tx); err != nil {
			logrus.Errorf("Failed to seed habitats: %v", err)
			return err
		}

		if err := seeds.SeedTaxonomyGBIF(tx,
			// Aselloidea families
			"Asellidae", "Stenasellidae",
			// Copepods families
			"Cyclopidae", "Parastenocarididae", "Canthocamptidae", "Ameiridae",
			"Chappuisiidae", "Diaptomidae", "Ectinosomatidae", "Gelyellidae",
			"Halicyclopidae", "Miraciidae", "Phyllognathopodidae",
		); err != nil {
			logrus.Errorf("Failed to seed taxonomy: %v", err)
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
		return nil
	})

	if err != nil {
		logrus.Errorf("Seeding failed: %v", err)
	}

	err = client.WithConfig(map[string]interface{}{
		"session_idle_transaction_timeout": timeout,
	}).Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {

		logrus.Info("⚙ Artificial datasets")
		datasets, err := seeds.LoadMultipleOccurrencesDatasets("data/datasets.json")
		if err != nil {
			logrus.Errorf("Failed to load datasets: %v", err)
			return err
		}
		if err := seeds.SeedOccurrencesDatasets(tx, datasets); err != nil {
			logrus.Errorf("Failed to seed datasets: %v", err)
			return err
		}

		logrus.Info("🧪 Empirical datasets")
		// logrus.Infof("🌱 Seeding WAD sampling sites")
		// if err := seeds.SeedSites(tx, *aselloidea); err != nil {
		// 	logrus.Errorf("Failed to seed Aselloidea sampling sites")
		// 	return err
		// }
		logrus.Infof("🌱 Seeding EGCop occurrences")
		copepoda, err := seeds.LoadOccurrencesDataset("data/Copepoda/Copepoda_occurrences.json")
		if err != nil {
			logrus.Errorf("Failed to load datasets: %v", err)
			return err
		}
		if err := seeds.SeedOccurrencesDatasets(tx, []occurrence.OccurrenceDatasetInput{*copepoda}); err != nil {
			return err
		}

		logrus.Infof("🌱 Seeding WAD occurrences")
		aselloidea, err := seeds.LoadOccurrencesDataset("data/Aselloidea/Aselloidea_occurrences.json")
		if err != nil {
			logrus.Errorf("Failed to load datasets: %v", err)
			return err
		}
		if err := seeds.SeedOccurrencesDatasets(tx, []occurrence.OccurrenceDatasetInput{*aselloidea}); err != nil {
			return err
		}

		logrus.Infof("⚙ Postprocessing...")
		// logrus.Infof("• generate bio-material codes")
		// if err := tx.Execute(context.Background(),
		// 	`#edgeql
		// 		update occurrence::BioMaterial set {};
		// 	`); err != nil {
		// 	return err
		// }
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
