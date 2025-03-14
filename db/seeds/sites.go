package seeds

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/models/occurrence"

	"github.com/sirupsen/logrus"
)

func LoadSiteDatasetJSON(file string) (dataset occurrence.SiteDatasetInput) {
	data, err := os.ReadFile(file)
	if err != nil {
		logrus.Fatalf("Failed to read sites JSON: %v", err)
	}

	if err := json.Unmarshal(data, &dataset); err != nil {
		logrus.Fatalf("Failed to parse sites JSON: %v", err)
	}
	return
}

func LoadSiteDataset(db geltypes.Executor, file string, maxAmount int) (*occurrence.SiteDatasetInput, error) {
	cfg, _ := config.LoadConfig("../../server", "config")
	logrus.Infof("Loaded config: %+v", cfg)

	dataset := LoadSiteDatasetJSON(file)
	if maxAmount > 0 {
		dataset.NewSites = dataset.NewSites[0:maxAmount]
	}

	// logrus.Infof("Making API call to Geoapify")
	// err := dataset.NewSites.FillPlaces(db, cfg.GeoApifyApiKey)
	// if err != nil {
	// 	return nil, err
	// }
	return &dataset, nil

}

func SeedSites(tx geltypes.Tx, dataset occurrence.SiteDatasetInput) error {
	dataset.InferCountry = true
	validated, errs := dataset.Validate(tx)
	if errs != nil {
		return errors.Join(errs...)
	}

	_, err := validated.SaveTx(tx)
	return err
}
