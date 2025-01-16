package seeds

import (
	"darco/proto/config"
	"darco/proto/models/occurrence"
	"encoding/json"
	"os"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

func loadSitesJSON(file string) occurrence.SiteImportDataset {
	data, err := os.ReadFile(file)
	if err != nil {
		logrus.Fatalf("Failed to read sites JSON: %v", err)
	}

	var sites occurrence.SiteImportDataset
	if err := json.Unmarshal(data, &sites); err != nil {
		logrus.Fatalf("Failed to parse sites JSON: %v", err)
	}
	return sites
}

func LoadSites(file string, maxAmount int) (*occurrence.SiteImportDataset, error) {
	cfg, _ := config.LoadConfig("../../server", "config")
	logrus.Infof("Loaded config: %+v", cfg)

	sitesInput := loadSitesJSON(file)
	sitesInput = sitesInput[0:maxAmount]

	logrus.Infof("Making API call to Geoapify")
	err := sitesInput.FillPlaces(cfg.GeoApifyApiKey)
	if err != nil {
		return nil, err
	}
	logrus.Infof("SITES\n%+v\n\n", sitesInput)
	return &sitesInput, nil
}

func SeedSites(db edgedb.Executor, sites occurrence.SiteImportDataset) error {
	_, err := sites.Save(db)
	return err
}
