package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"seeds"
	"seeds/email"

	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/dataset"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/settings"

	_ "net/http/pprof"

	"github.com/sirupsen/logrus"
)

var entities = []string{
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

	// logrus.Println("Enabling pprof for profiling")
	// go func() {
	// 	logrus.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	database := flag.String("db", "", "The name of the database to seed")
	BATCH_SIZE := flag.Int("batch", 200, "The batch size for bulk inserts")
	N_CORES := flag.Int("cores", 10, "The number of cores to use for parallel inserts")
	flag.Parse()

	timeout, _ := geltypes.ParseDuration("15m")
	client := db.Connect(gelcfg.Options{Database: *database, Concurrency: 10}).WithConfig(map[string]interface{}{
		"session_idle_transaction_timeout": timeout,
	})

	// aselloidea, err := seeds.LoadSiteDataset(client, "data/Aselloidea/sites.json")
	// if err != nil {
	// 	logrus.Fatalf("Failed to load Asellidae sites: %v", err)
	// }
	path, err := os.Getwd()
	if err != nil {
		logrus.Println(err)
	}
	fmt.Println(path)

	err = client.WithConfig(map[string]interface{}{
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

	tracker := &occurrence.OccurrenceBatchProgressBar{}

	// err = client.WithConfig(map[string]interface{}{
	// 	"session_idle_transaction_timeout": timeout,
	// }).Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {

	var createdDatasets []dataset.Dataset
	rollbackDatasets := func() {
		for _, d := range createdDatasets {
			logrus.Infof("Rolling back dataset %s", d.Slug)
			if err := d.RollbackImport(client); err != nil {
				logrus.Errorf("Failed to rollback dataset %s: %v", d.Slug, err)
			}
		}
	}

	logrus.Info("⚙ Artificial datasets")
	datasets, err := seeds.LoadMultipleOccurrencesDatasets("data/datasets.json")
	if err != nil {
		logrus.Errorf("Failed to load datasets: %v", err)
		return
	}
	for i := range datasets {
		d, err := datasets[i].SetTracker(tracker).SaveParallel(client, *BATCH_SIZE, *N_CORES)
		if err != nil {
			rollbackDatasets()
			logrus.Errorf("❗Failed to seed occurrence dataset: %v", err)
			return
		}
		createdDatasets = append(createdDatasets, d)
	}

	logrus.Info("🧪 Empirical datasets")
	logrus.Infof("🌱 Seeding EGCop occurrences")
	copepoda, err := seeds.LoadOccurrencesDataset("data/Copepoda/Copepoda_occurrences.json")
	if err != nil {
		rollbackDatasets()
		logrus.Fatalf("Failed to load datasets: %v", err)
	}
	datasetCopepoda, err := copepoda.SetTracker(tracker).SaveParallel(client, *BATCH_SIZE, *N_CORES)
	if err != nil {
		rollbackDatasets()
		logrus.Fatalf("Failed to seed Copepoda occurrences: %v", err)
	}
	createdDatasets = append(createdDatasets, datasetCopepoda)

	logrus.Infof("🌱 Seeding WAD occurrences")
	aselloidea, err := seeds.LoadOccurrencesDataset("data/Aselloidea/Aselloidea_occurrences.json")
	if err != nil {
		rollbackDatasets()
		logrus.Fatalf("Failed to load datasets: %v", err)
	}

	datasetAselloidea, err := aselloidea.SetTracker(tracker).SaveParallel(client, *BATCH_SIZE, *N_CORES)
	if err != nil {
		rollbackDatasets()
		logrus.Fatalf("Failed to seed Aselloidea occurrences: %v", err)
	}
	createdDatasets = append(createdDatasets, datasetAselloidea)

	logrus.Infof("⚙ Postprocessing...")
	// logrus.Infof("• generate bio-material codes")
	// if err := tx.Execute(context.Background(),
	// 	`#edgeql
	// 		update occurrence::BioMaterial set {};
	// 	`); err != nil {
	// 	return err
	// }
	logrus.Infof("• generate sequence codes")
	if err := client.Execute(context.Background(),
		`#edgeql
				update seq::ExternalSequence set {};
			`); err != nil {
		logrus.Fatalf("Failed to generate sequence codes: %v", err)
	}
}
