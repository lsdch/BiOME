package services

import (
	"context"
	"fmt"
	"io"
	"net/http"

	_ "embed"

	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/data"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/sirupsen/logrus"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"gopkg.in/yaml.v3"
)

type BootstrapService struct {
	db        *db.DB
	config    config.BootstrapConfig
	settings  *SettingsService
	samplings *SamplingService
	habitats  *HabitatService
}

func NewBootstrapService(db *db.DB,
	config config.BootstrapConfig,
	settings *SettingsService,
	samplings *SamplingService,
	habitats *HabitatService,
) *BootstrapService {
	return &BootstrapService{
		db:        db,
		config:    config,
		settings:  settings,
		samplings: samplings,
		habitats:  habitats,
	}
}

func (s *BootstrapService) Bootstrap(ctx context.Context) error {

	if err := s.settings.Bootstrap(ctx); err != nil {
		return fmt.Errorf("bootstrap settings: %w", err)
	}

	ok, err := s.settings.TestSMTPConnection(ctx)
	if err != nil {
		return fmt.Errorf("test SMTP connection: %w", err)
	}
	if !ok {
		return fmt.Errorf("SMTP connection failed")
	} else {
		logrus.Infof("SMTP connection successful")
	}

	if err := s.BootstrapCountries(ctx); err != nil {
		return fmt.Errorf("bootstrap countries: %w", err)
	}

	return nil
}

func (s *BootstrapService) BootstrapCountries(ctx context.Context) error {
	existingCountries, err := s.db.Queries().ListCountries(ctx)
	if err != nil {
		return err
	}
	if len(existingCountries) > 0 {
		logrus.Infof("Countries already exist in the database, skipping bootstrap")
		return nil
	}

	logrus.Infof("Bootstrapping countries from %s", s.config.CountriesJSON_URL)
	countriesJSON, err := downloadCountriesJSON(ctx, s.config.CountriesJSON_URL)
	if err != nil {
		return fmt.Errorf("failed to download countries.json: %w", err)
	}

	// Parse the JSON and insert countries into the database
	// Assuming you have a function to parse the JSON and return a slice of country structs
	countries, err := s.parseCountriesJSON(countriesJSON)
	if err != nil {
		return fmt.Errorf("failed to parse countries.json: %w", err)
	}

	return s.db.WithTx(ctx, func(q *biomedb.Queries) error {
		for _, country := range countries {
			err := q.InsertCountry(ctx, country.ToParams())
			if err != nil {
				return fmt.Errorf("failed to insert country: %w", err)
			}
		}
		return nil
	})
}

func (s *BootstrapService) parseCountriesJSON(data []byte) ([]models.CountryInput, error) {
	var fc = new(geojson.FeatureCollection)
	err := fc.UnmarshalJSON(data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal countries.json: %w", err)
	}

	var countries = make([]models.CountryInput, len(fc.Features))

	for i, feature := range fc.Features {

		name, ok := feature.Properties[s.config.CountryNameKey].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid '%s' property for feature %d", s.config.CountryNameKey, i)
		}
		code, ok := feature.Properties[s.config.CountryCodeKey].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid '%s' property for feature %d", s.config.CountryCodeKey, i)
		}
		region, ok := feature.Properties[s.config.CountryContinentKey].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid '%s' property for feature %d", s.config.CountryContinentKey, i)
		}
		subRegion, ok := feature.Properties[s.config.CountrySubcontinentKey].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid '%s' property for feature %d", s.config.CountrySubcontinentKey, i)
		}
		geom, ok := feature.Geometry.(*geom.MultiPolygon)
		if !ok || geom == nil {
			return nil, fmt.Errorf("missing geometry for feature %d", i)
		}
		geomBytes, err := geojson.Marshal(geom)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal geometry for feature %d: %w", i, err)
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

func (s *BootstrapService) BootstrapSamplingMetadata(ctx context.Context) error {
	methodsBytes, err := data.DataFS.ReadFile("sampling_methods.yaml")
	if err != nil {
		return fmt.Errorf("failed to read sampling methods: %w", err)
	}
	var methods []biomedb.CreateSamplingMethodParams
	err = yaml.Unmarshal(methodsBytes, &methods)
	if err != nil {
		return fmt.Errorf("failed to unmarshal sampling methods: %w", err)
	}
	for _, method := range methods {
		_, err := s.samplings.CreateSamplingMethod(ctx, method)
		if err != nil {
			return fmt.Errorf("failed to create sampling method %s: %w", method.Code, err)
		}
	}

	fixativesBytes, err := data.DataFS.ReadFile("fixatives.yaml")
	if err != nil {
		return fmt.Errorf("failed to read fixatives: %w", err)
	}
	var fixatives []biomedb.CreateFixativeParams
	err = yaml.Unmarshal(fixativesBytes, &fixatives)
	if err != nil {
		return fmt.Errorf("failed to unmarshal fixatives: %w", err)
	}
	for _, fixative := range fixatives {
		_, err := s.samplings.CreateSamplingFixative(ctx, fixative)
		if err != nil {
			return fmt.Errorf("failed to create fixative %s: %w", fixative.Code, err)
		}
	}

	habitatsBytes, err := data.DataFS.ReadFile("habitats.yaml")
	if err != nil {
		return fmt.Errorf("failed to read habitats: %w", err)
	}
	var habitatGroups []models.HabitatGroupInput
	err = yaml.Unmarshal(habitatsBytes, &habitatGroups)
	if err != nil {
		return fmt.Errorf("failed to unmarshal habitats: %w", err)
	}
	for _, group := range habitatGroups {
		err := s.habitats.CreateHabitatGroup(ctx, group)
		if err != nil {
			return fmt.Errorf("failed to create habitat group %s: %w", group.Label, err)
		}
	}

	return nil
}
