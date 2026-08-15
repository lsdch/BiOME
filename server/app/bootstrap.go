package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "embed"

	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/data"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/stores"
	"github.com/sirupsen/logrus"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

type AppBootstrap struct {
	db              *db.DB
	config          config.BootstrapConfig
	services        *AppServices
	taxonResolution *stores.TaxonResolutionStore
}

func NewAppBootstrap(db *db.DB,
	config config.BootstrapConfig,
	services *AppServices,
) *AppBootstrap {
	taxonResolution := stores.NewTaxonResolutionStore()
	return &AppBootstrap{
		db:              db,
		config:          config,
		services:        services,
		taxonResolution: taxonResolution,
	}
}

func (s *AppBootstrap) Bootstrap(ctx context.Context) error {

	if err := s.services.SettingsService.Bootstrap(ctx, s.db); err != nil {
		return fmt.Errorf("bootstrap settings: %w", err)
	}

	ok, err := s.services.SettingsService.TestSMTPConnection(ctx)
	if err != nil {
		return fmt.Errorf("test SMTP connection: %w", err)
	}
	if !ok {
		return fmt.Errorf("SMTP connection failed")
	} else {
		logrus.Infof("SMTP connection successful")
	}

	if err := s.services.AccountsService.BootstrapUsers(ctx, s.db); err != nil {
		return fmt.Errorf("bootstrap users: %w", err)
	}

	if err := s.BootstrapCountries(ctx); err != nil {
		return fmt.Errorf("bootstrap countries: %w", err)
	}

	if err := s.BootstrapGBIFKingdoms(ctx); err != nil {
		return fmt.Errorf("bootstrap GBIF kingdoms: %w", err)
	}

	if err := s.BootstrapSamplingMetadata(ctx); err != nil {
		return fmt.Errorf("bootstrap sampling metadata: %w", err)
	}

	logrus.Infof("Bootstrap completed successfully")

	return nil
}

func (s *AppBootstrap) BootstrapGBIFKingdoms(ctx context.Context) error {

	kingdoms, err := s.services.TaxonomyService.FetchGBIFKingdoms(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch GBIF kingdoms: %w", err)
	}
	logrus.Infof("Bootstrapping %d GBIF kingdoms into the database", len(kingdoms))
	return s.taxonResolution.InsertGBIFBatch(ctx, s.db, kingdoms)
}

func loadCountriesJSON(ctx context.Context, url string, cachePath string) ([]byte, error) {
	// Try cache first
	info, err := os.Stat(cachePath)
	// invalidate cache after 1 year
	if err == nil && time.Since(info.ModTime()) < 365*30*24*time.Hour {
		logrus.Infof("Loading countries from cache at %s", cachePath)
		return os.ReadFile(cachePath)
	}

	logrus.Infof("Downloading countries from %s", url)

	data, err := downloadCountriesJSON(ctx, url)
	if err != nil {
		return nil, err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return nil, fmt.Errorf("create cache directory: %w", err)
	}

	// Save cache
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		logrus.Warnf("failed to write countries cache: %v", err)
	} else {
		logrus.Infof("Saved countries.json cache to %s", cachePath)
	}

	return data, nil
}

func (s *AppBootstrap) BootstrapCountries(ctx context.Context) error {
	existingCountries, err := s.db.Queries().ListCountries(ctx)
	if err != nil {
		return err
	}
	if len(existingCountries) > 0 {
		logrus.Infof("%d countries already exist in the database, skipping bootstrap", len(existingCountries))
		return nil
	}

	countriesJSON, err := loadCountriesJSON(ctx, s.config.Countries.CountriesJSON_URL, s.config.Countries.CountryJSON_CachePath)
	if err != nil {
		return fmt.Errorf("failed to load countries.json: %w", err)
	}

	// Parse the JSON and insert countries into the database
	// Assuming you have a function to parse the JSON and return a slice of country structs
	countries, err := s.parseCountriesJSON(countriesJSON)
	if err != nil {
		return fmt.Errorf("failed to parse countries.json: %w", err)
	}

	logrus.Infof("Bootstrapping %d countries into the database", len(countries))

	return s.db.WithTx(ctx, func(tx *db.Tx) error {
		for _, country := range countries {
			err := tx.Queries().InsertCountry(ctx, country.ToParams())
			if err != nil {
				return fmt.Errorf("failed to insert country: %w", err)
			}
		}
		return nil
	})
}

func (s *AppBootstrap) parseCountriesJSON(data []byte) ([]models.CountryInput, error) {
	var fc = new(geojson.FeatureCollection)
	err := fc.UnmarshalJSON(data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal countries.json: %w", err)
	}

	var countries = make([]models.CountryInput, len(fc.Features))

	for i, feature := range fc.Features {

		name, ok := feature.Properties[s.config.Countries.CountryNameKey].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid '%s' property for feature %d", s.config.Countries.CountryNameKey, i)
		}
		code, ok := feature.Properties[s.config.Countries.CountryCodeKey].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid '%s' property for feature %d", s.config.Countries.CountryCodeKey, i)
		}
		region, ok := feature.Properties[s.config.Countries.CountryContinentKey].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid '%s' property for feature %d", s.config.Countries.CountryContinentKey, i)
		}
		subRegion, ok := feature.Properties[s.config.Countries.CountrySubcontinentKey].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid '%s' property for feature %d", s.config.Countries.CountrySubcontinentKey, i)
		}

		var (
			geometry    geom.T
			hasGeometry = true
		)
		switch g := feature.Geometry.(type) {
		case *geom.MultiPolygon:
			geometry = g
		case *geom.Polygon:
			geometry = g
		default:
			logrus.Warnf("unsupported geometry type %T for feature %d", g, i)
			hasGeometry = false
		}
		var geomBytes []byte
		if hasGeometry {
			geomBytes, err = geojson.Marshal(geometry)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal geometry for feature %d: %w", i, err)
			}
		}
		countries[i] = models.CountryInput{
			Name:         name,
			Code:         code,
			Continent:    region,
			Subcontinent: subRegion,
			GeometryJSON: geomBytes,
		}
	}

	return countries, nil
}

func downloadCountriesJSON(ctx context.Context, countriesURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, countriesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download countries.json: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return data, nil
}

func (s *AppBootstrap) BootstrapSamplingMetadata(ctx context.Context) error {
	methodsBytes, err := data.DataFS.ReadFile("sampling_methods.yaml")
	if err != nil {
		return fmt.Errorf("failed to read sampling methods: %w", err)
	}
	if err := s.services.SamplingsService.BootstrapSamplingMethods(ctx, s.db, methodsBytes); err != nil {
		return fmt.Errorf("failed to bootstrap sampling methods: %w", err)
	}

	fixativesBytes, err := data.DataFS.ReadFile("fixatives.yaml")
	if err != nil {
		return fmt.Errorf("failed to read fixatives: %w", err)
	}
	if err := s.services.SamplingsService.BootstrapSamplingFixatives(ctx, s.db, fixativesBytes); err != nil {
		return fmt.Errorf("failed to bootstrap sampling fixatives: %w", err)
	}

	habitatsBytes, err := data.DataFS.ReadFile("habitats.yaml")
	if err != nil {
		return fmt.Errorf("failed to read habitats: %w", err)
	}
	return s.db.WithTx(ctx, func(tx *db.Tx) error {
		if err := s.services.HabitatService.BootstrapHabitats(ctx, tx, habitatsBytes); err != nil {
			return fmt.Errorf("failed to bootstrap habitats: %w", err)
		}
		return nil
	})

}
