package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"seeds"
	"seeds/config"

	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/dataset"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/spf13/viper"

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

func main() {

	// logrus.Println("Enabling pprof for profiling")
	// go func() {
	// 	logrus.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	configPath := flag.String("cfg", "db/seeds/config", "Path to the configuration directory")

	database := flag.String("db", "", "The name of the database to seed")
	BATCH_SIZE := flag.Int("batch", 200, "The batch size for bulk inserts")
	N_CORES := flag.Int("cores", 10, "The number of cores to use for parallel inserts")
	flag.Parse()

	timeout, _ := geltypes.ParseDuration("15m")
	client := db.Connect(gelcfg.Options{Database: *database, Concurrency: 10}).WithConfig(map[string]interface{}{
		"session_idle_transaction_timeout": timeout,
	})

	cwd, _ := os.Getwd()
	logrus.Infof("Using configuration directory: %s/%s", cwd, *configPath)
	v := viper.New()
	v.AddConfigPath(filepath.Join(cwd, *configPath))

	logrus.Infof("Loading JSON datasets...")
	datasets, err := seeds.LoadMultipleOccurrencesDatasets("data/datasets.json")
	if err != nil {
		logrus.Fatalf("Failed to load datasets: %v", err)
	}
	copepoda, err := seeds.LoadOccurrencesDataset("data/Copepoda/Copepoda_occurrences.json")
	if err != nil {
		logrus.Fatalf("Failed to load datasets: %v", err)
	}
	aselloidea, err := seeds.LoadOccurrencesDataset("data/Aselloidea/Aselloidea_occurrences.json")
	if err != nil {
		logrus.Fatalf("Failed to load datasets: %v", err)
	}

	subcommand := flag.Arg(0)
	switch subcommand {

	case "":
		logrus.Fatalf("Please provide a subcommand")

	case "init":
		err := client.Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {
			cfg, err := config.LoadConfig[config.InstanceConfig](v, "config")
			if err != nil {
				return fmt.Errorf("Failed to load instance configuration: %v", err)
			}
			err = cfg.SaveTx(tx)
			if err != nil {
				return fmt.Errorf("Failed to initialize instance configuration: %v", err)
			}
			return v.WriteConfig()
		})
		if err != nil {
			logrus.Fatal(err)
		}
		return

	case "countries":
		logrus.Infof("🌱 Seeding countries")
		if err := seeds.SeedCountriesGeoJSON(client, "./data/remote/countries.json"); err != nil {
			logrus.Fatalf("Failed to seed countries: %v", err)
		}
		return

	case "habitats":
		habitats, err := occurrence.ListHabitats(client)
		if err != nil {
			logrus.Fatalf("Failed to list habitats: %v", err)
		}
		if len(habitats) > 0 {
			logrus.Infof("✅ Habitats already seeded, skipping.")
			return
		}
		logrus.Infof("🌱 Seeding habitats")
		if err := occurrence.InitialHabitatsSetup(client); err != nil {
			logrus.Errorf("Failed to seed habitats: %v", err)
			return
		}
		return

	case "taxonomy":
		if err := seeds.SeedTaxonomyGBIF(client,
			// Aselloidea families
			"Asellidae", "Stenasellidae",
			// Copepods families
			"Cyclopidae", "Parastenocarididae", "Canthocamptidae", "Ameiridae",
			"Chappuisiidae", "Diaptomidae", "Ectinosomatidae", "Gelyellidae",
			"Halicyclopidae", "Miraciidae", "Phyllognathopodidae",
		); err != nil {
			logrus.Errorf("Failed to seed taxonomy: %v", err)
			return
		}

	case "metadata":
		logrus.Infof("🌱 Seeding...")
		for _, entity := range entities {
			logrus.Infof("• %s", entity)
			err := seeds.Seed(client, entity)
			if err != nil {
				logrus.Errorf("Failed to seed %s: %v", entity, err)
			}
		}
		return

	case "missing-taxa":
		logrus.Infof("Checking for missing taxa in datasets...")
		err = CheckMissingTaxa(client, append(datasets, copepoda, aselloidea)...)
		if err != nil {
			logrus.Fatalf("Missing taxa detected: %v", err)
		}
		return

	case "datasets":
		err = client.WithConfig(map[string]interface{}{
			"session_idle_transaction_timeout": timeout,
		}).Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {
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
			if err != nil {
				return fmt.Errorf("Failed to load datasets: %v", err)
			}
			for i := range datasets {
				datasets[i].OccurrenceBatchMetadataInputs.Taxa = nil
				d, err := datasets[i].SetTracker(tracker).SaveParallel(client, *BATCH_SIZE, *N_CORES)
				if err != nil {
					rollbackDatasets()
					return fmt.Errorf("❗Failed to seed occurrence dataset: %v", err)
				}
				createdDatasets = append(createdDatasets, d)
			}

			logrus.Info("🧪 Empirical datasets")
			logrus.Infof("🌱 Seeding EGCop occurrences")
			copepoda.OccurrenceBatchMetadataInputs.Taxa = nil
			datasetCopepoda, err := copepoda.SetTracker(tracker).SaveParallel(client, *BATCH_SIZE, *N_CORES)
			if err != nil {
				rollbackDatasets()
				logrus.Fatalf("Failed to seed Copepoda occurrences: %v", err)
			}
			createdDatasets = append(createdDatasets, datasetCopepoda)

			logrus.Infof("🌱 Seeding WAD occurrences")

			aselloidea.OccurrenceBatchMetadataInputs.Taxa = nil
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
			return nil
		})
		return
	}

	// err = client.WithConfig(map[string]interface{}{
	// 	"session_idle_transaction_timeout": timeout,
	// }).Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {

	// 	// logrus.Infof("⚙ Initializing settings with superadmin account")
	// 	// superAdmin, err := superAdminInput.Save(tx)
	// 	// if err != nil {
	// 	// 	return fmt.Errorf("Failed to initialize super admin account: %v", err)
	// 	// }

	// 	// if err := (settings.SettingsInput{
	// 	// 	SuperAdminID: superAdmin.ID,
	// 	// 	Instance: settings.InstanceSettingsInput{
	// 	// 		InstanceSettingsInner: settings.InstanceSettingsInner{
	// 	// 			Name: "[BiOME prototype]",
	// 	// 		},
	// 	// 		Description: models.NewOptionalNull("Prototype BiOME instance"),
	// 	// 	},
	// 	// }).SaveTx(tx); err != nil {
	// 	// 	return fmt.Errorf("Failed to initialize settings: %v", err)
	// 	// }

	// 	// if err := email.SetupEmailConfig(client, email.EmailSetupArgs{}); err != nil {
	// 	// 	return err
	// 	// }

	// 	// logrus.Infof("🌱 Seeding habitats")
	// 	// if err := occurrence.InitialHabitatsSetup(tx); err != nil {
	// 	// 	logrus.Errorf("Failed to seed habitats: %v", err)
	// 	// 	return err
	// 	// }

	// 	// if err := seeds.SeedTaxonomyGBIF(tx,
	// 	// 	// Aselloidea families
	// 	// 	"Asellidae", "Stenasellidae",
	// 	// 	// Copepods families
	// 	// 	"Cyclopidae", "Parastenocarididae", "Canthocamptidae", "Ameiridae",
	// 	// 	"Chappuisiidae", "Diaptomidae", "Ectinosomatidae", "Gelyellidae",
	// 	// 	"Halicyclopidae", "Miraciidae", "Phyllognathopodidae",
	// 	// ); err != nil {
	// 	// 	logrus.Errorf("Failed to seed taxonomy: %v", err)
	// 	// 	return err
	// 	// }

	// 	// logrus.Infof("🌱 Seeding...")
	// 	// for _, entity := range entities {
	// 	// 	logrus.Infof("• %s", entity)
	// 	// 	err := seeds.Seed(tx, entity)
	// 	// 	if err != nil {
	// 	// 		return err
	// 	// 	}
	// 	// }
	// 	// return nil
	// // })

	// if err != nil {
	// 	logrus.Errorf("Seeding failed: %v", err)
	// }

	// logrus.Infof("Checking for missing taxa in datasets...")
	// err = CheckMissingTaxa(client, append(datasets, copepoda, aselloidea)...)
	// if err != nil {
	// 	logrus.Fatalf("Missing taxa detected: %v", err)
	// }

	// tracker := &occurrence.OccurrenceBatchProgressBar{}

	// // err = client.WithConfig(map[string]interface{}{
	// // 	"session_idle_transaction_timeout": timeout,
	// // }).Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {

	// var createdDatasets []dataset.Dataset
	// rollbackDatasets := func() {
	// 	for _, d := range createdDatasets {
	// 		logrus.Infof("Rolling back dataset %s", d.Slug)
	// 		if err := d.RollbackImport(client); err != nil {
	// 			logrus.Errorf("Failed to rollback dataset %s: %v", d.Slug, err)
	// 		}
	// 	}
	// }

	// logrus.Info("⚙ Artificial datasets")
	// if err != nil {
	// 	logrus.Errorf("Failed to load datasets: %v", err)
	// 	return
	// }
	// for i := range datasets {
	// 	datasets[i].OccurrenceBatchMetadataInputs.Taxa = nil
	// 	d, err := datasets[i].SetTracker(tracker).SaveParallel(client, *BATCH_SIZE, *N_CORES)
	// 	if err != nil {
	// 		rollbackDatasets()
	// 		logrus.Errorf("❗Failed to seed occurrence dataset: %v", err)
	// 		return
	// 	}
	// 	createdDatasets = append(createdDatasets, d)
	// }

	// logrus.Info("🧪 Empirical datasets")
	// logrus.Infof("🌱 Seeding EGCop occurrences")
	// copepoda.OccurrenceBatchMetadataInputs.Taxa = nil
	// datasetCopepoda, err := copepoda.SetTracker(tracker).SaveParallel(client, *BATCH_SIZE, *N_CORES)
	// if err != nil {
	// 	rollbackDatasets()
	// 	logrus.Fatalf("Failed to seed Copepoda occurrences: %v", err)
	// }
	// createdDatasets = append(createdDatasets, datasetCopepoda)

	// logrus.Infof("🌱 Seeding WAD occurrences")

	// aselloidea.OccurrenceBatchMetadataInputs.Taxa = nil
	// datasetAselloidea, err := aselloidea.SetTracker(tracker).SaveParallel(client, *BATCH_SIZE, *N_CORES)
	// if err != nil {
	// 	rollbackDatasets()
	// 	logrus.Fatalf("Failed to seed Aselloidea occurrences: %v", err)
	// }
	// createdDatasets = append(createdDatasets, datasetAselloidea)

	// logrus.Infof("⚙ Postprocessing...")
	// // logrus.Infof("• generate bio-material codes")
	// // if err := tx.Execute(context.Background(),
	// // 	`#edgeql
	// // 		update occurrence::BioMaterial set {};
	// // 	`); err != nil {
	// // 	return err
	// // }
	// logrus.Infof("• generate sequence codes")
	// if err := client.Execute(context.Background(),
	// 	`#edgeql
	// 			update seq::ExternalSequence set {};
	// 		`); err != nil {
	// 	logrus.Fatalf("Failed to generate sequence codes: %v", err)
	// }
}

func CheckMissingTaxa(client geltypes.Executor, datasets ...*occurrence.OccurrenceDatasetInput) (err error) {
	var missingTaxa = make(map[string][]string)
	for _, dataset := range datasets {
		for _, taxon := range dataset.OccurrenceBatchMetadataInputs.Taxa {
			if _, err := taxon.Save(client); err != nil {
				return fmt.Errorf("failed to save taxon %s: %w", taxon.Name, err)
			}
		}
		dataset.OccurrenceBatchMetadataInputs.Taxa = nil
		mt, err := dataset.ListMissingTaxa(client)
		if err != nil {
			logrus.Fatalf("Failed to find missing taxa in dataset %s: %v", dataset.Label, err)
		}
		if len(mt) > 0 {
			missingTaxa[dataset.Label] = mt
		}
	}
	if len(missingTaxa) > 0 {
		logrus.Errorf("❗ Missing taxa detected in datasets")
		fh, _ := os.Create("missing_taxa.txt")
		defer fh.Close()
		for ds, taxa := range missingTaxa {
			for _, t := range taxa {
				fmt.Fprintf(fh, "%s\t%s\n", ds, t)
			}
		}
		return errors.New("Please add the missing taxa listed in missing_taxa.txt to the taxonomy before proceeding")
	}
	return nil
}
