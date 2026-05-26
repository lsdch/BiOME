package main

import (
	"context"
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
	"persons",
	"users",
	"articles",
	// "programs",
	"sampling_methods",
	"fixatives",
	"abiotic",
	"genes",
	"data_sources",
}

func main() {
	// logrus.SetLevel(logrus.DebugLevel)
	// logrus.SetLevel(logrus.InfoLevel)

	// fh, err := os.OpenFile("seed.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	// if err != nil {
	// 	logrus.Fatalf("Failed to open log file: %v", err)
	// }
	// defer fh.Close()
	// logrus.SetOutput(fh)
	//
	// logrus.Println("Enabling pprof for profiling")
	// go func() {
	// 	logrus.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	configPath := flag.String("cfg", "db/seeds/config", "Path to the configuration directory")

	database := flag.String("db", "", "The name of the database to seed")
	BATCH_SIZE := flag.Int("batch", 100, "The batch size for bulk inserts")
	N_CORES := flag.Int("cores", 10, "The number of cores to use for parallel inserts")
	LOG_OUTPUT := flag.String("log", "", "Path to log output file (default: stdout)")
	DEBUG := flag.Bool("debug", false, "Enable debug logging")

	flag.Parse()

	if *DEBUG {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.Debug("Debug logging enabled")

	if LOG_OUTPUT != nil && *LOG_OUTPUT != "" {
		fh, err := os.OpenFile("seed.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			logrus.Fatalf("Failed to open log file: %v", err)
		}
		defer fh.Close()
		logrus.SetOutput(fh)
	}

	timeout, _ := geltypes.ParseDuration("15m")
	client := db.Connect(gelcfg.Options{Database: *database, Concurrency: 10}).WithConfig(map[string]interface{}{
		"session_idle_transaction_timeout": timeout,
	})

	logrus.Infof("Loading configuration\n")

	cwd, _ := os.Getwd()
	logrus.Infof("Using configuration directory: %s/%s", cwd, *configPath)
	v := viper.New()
	v.AddConfigPath(filepath.Join(cwd, *configPath))

	logrus.Infof("Loading JSON datasets...")
	datasets, err := seeds.LoadMultipleOccurrencesDatasets("data/datasets.json")
	if err != nil {
		logrus.Fatalf("Failed to load datasets: %v", err)
	}
	copepoda, err := seeds.LoadOccurrencesDataset("data/datasets/Copepoda_occurrences.json")
	if err != nil {
		logrus.Fatalf("Failed to load datasets: %v", err)
	}
	ostracoda, err := seeds.LoadOccurrencesDataset("data/datasets/Ostracoda_occurrences.json")
	if err != nil {
		logrus.Fatalf("Failed to load datasets: %v", err)
	}

	aselloidea, err := seeds.LoadOccurrencesDataset("data/datasets/Aselloidea_occurrences.json")
	if err != nil {
		logrus.Fatalf("Failed to load datasets: %v", err)
	}

	austria, err := seeds.LoadOccurrencesDataset("data/datasets/Asellidae_Austria.json")
	if err != nil {
		logrus.Fatalf("Failed to load datasets: %v", err)
	}

	spiders, err := seeds.LoadOccurrencesDataset("data/datasets/Spiders.json")
	if err != nil {
		logrus.Fatalf("Failed to load datasets: %v", err)
	}

	subcommand := flag.Arg(0)
	logrus.Infof("Running seed subcommand: %s", subcommand)
	switch subcommand {

	case "":
		logrus.Fatalf("Please provide a subcommand")

	case "init":
		logrus.Infof("🌱 Initializing instance configuration")
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
		if err := seeds.SeedCountriesGeoJSON(client, "../../data/remote/countries.json"); err != nil {
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

	case "datasets":
		target_dataset := flag.Arg(1)
		err = client.WithConfig(map[string]interface{}{
			"session_idle_transaction_timeout": timeout,
		}).Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {
			tracker := &occurrence.OccurrenceBatchProgressBar{}

			// err = client.WithConfig(map[string]interface{}{
			// 	"session_idle_transaction_timeout": timeout,
			// }).Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {

			var createdDatasets []dataset.Dataset
			handleDatasetImportError := func() {

				for _, d := range createdDatasets {
					logrus.Infof("Dataset %s was imported.", d.Label)
					// logrus.Infof("Rolling back dataset %s", d.Slug)
					// if err := d.RollbackImport(client); err != nil {
					// 	logrus.Errorf("Failed to rollback dataset %s: %v", d.Slug, err)
					// }
				}
			}

			if target_dataset == "" || target_dataset == "artificial" {
				logrus.Info("⚙ Artificial datasets")
				if err != nil {
					return fmt.Errorf("Failed to load datasets: %v", err)
				}
				for i := range datasets {
					datasets[i].OccurrenceBatchMetadataInputs.Taxa = nil
					d, err := datasets[i].SetTracker(tracker).SaveParallel(client, *BATCH_SIZE, *N_CORES)
					if err != nil {
						handleDatasetImportError()
						return fmt.Errorf("❗Failed to seed occurrence dataset: %v", err)
					}
					createdDatasets = append(createdDatasets, d)
				}
			}

			logrus.Info("🧪 Empirical datasets")

			if target_dataset == "" || target_dataset == "ostracoda" {
				logrus.Infof("🌱 Seeding Ostracoda occurrences")
				datasetOstracoda, err := ostracoda.SetTracker(tracker).SaveParallel(client, *BATCH_SIZE, *N_CORES)
				if err != nil {
					handleDatasetImportError()
					logrus.Fatalf("Failed to seed Ostracoda occurrences: %v", err)
				}
				createdDatasets = append(createdDatasets, datasetOstracoda)
			}

			if target_dataset == "" || target_dataset == "copepoda" {
				logrus.Infof("🌱 Seeding Copepoda occurrences")
				datasetCopepoda, err := copepoda.SetTracker(tracker).SaveParallel(client, *BATCH_SIZE, *N_CORES)
				if err != nil {
					handleDatasetImportError()
					logrus.Fatalf("Failed to seed Copepoda occurrences: %v", err)
				}
				createdDatasets = append(createdDatasets, datasetCopepoda)
			}
			if target_dataset == "" || target_dataset == "aselloidea" {
				logrus.Infof("🌱 Seeding Aselloidea occurrences")
				datasetAselloidea, err := aselloidea.SetTracker(tracker).SaveParallel(client, *BATCH_SIZE, *N_CORES)
				if err != nil {
					handleDatasetImportError()
					logrus.Fatalf("Failed to seed Aselloidea occurrences: %v", err)
				}
				createdDatasets = append(createdDatasets, datasetAselloidea)
			}

			if target_dataset == "" || target_dataset == "austria" {
				logrus.Infof("🌱 Seeding Asellidae occurrences from Austria")
				datasetAustria, err := austria.SetTracker(tracker).SaveParallel(client, *BATCH_SIZE, *N_CORES)
				if err != nil {
					handleDatasetImportError()
					logrus.Fatalf("Failed to seed Asellidae occurrences from Austria: %v", err)
				}
				createdDatasets = append(createdDatasets, datasetAustria)
			}

			if target_dataset == "" || target_dataset == "spiders" {
				logrus.Infof("🌱 Seeding Spider occurrences")
				datasetSpiders, err := spiders.SetTracker(tracker).SaveParallel(client, *BATCH_SIZE, *N_CORES)
				if err != nil {
					handleDatasetImportError()
					logrus.Fatalf("Failed to seed Spider occurrences: %v", err)
				}
				createdDatasets = append(createdDatasets, datasetSpiders)
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
			if err := client.Execute(context.Background(),
				`#edgeql
				update seq::ExternalSequence set {};
			`); err != nil {
				logrus.Fatalf("Failed to generate sequence codes: %v", err)
			}
			return nil
		})
		if err != nil {
			logrus.Errorf("Seeding datasets failed: %v", err)
		}
		return
	}
}
