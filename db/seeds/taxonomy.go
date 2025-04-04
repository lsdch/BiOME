package seeds

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	gbif "github.com/lsdch/biome/models/taxonomy/GBIF"

	"github.com/geldata/gel-go/geltypes"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

func SeedTaxonomyGBIF(db geltypes.Tx, groups ...string) error {
	_ = db.Execute(context.Background(), "delete taxonomy::Taxon")
	for _, g := range groups {
		err := seedTaxonomyGroup(db, g)
		if err != nil {
			return fmt.Errorf("failed to seed taxonomy group %s: %w", g, err)
		}
	}
	return nil
}

func FetchSpeciesSuggestions(name string) (*gbif.TaxonInnerGBIF, error) {

	// url := "https://api.gbif.org/v1/species/suggest?datasetKey=d7dddbf4-2cf0-4f39-9b2a-bb099caae36c&status=ACCEPTED&nameType=SCIENTIFIC&q=Stenasellidae"

	query := url.Values{}
	query.Add("status", "ACCEPTED")
	query.Add("nameType", "SCIENTIFIC")
	query.Add("q", name)

	u := url.URL{
		Host:     "api.gbif.org",
		Path:     "v1/species/suggest",
		Scheme:   "https",
		RawQuery: query.Encode(),
	}

	// Create a new request with a timeout context.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute the request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	// Check for a non-200 status code.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	// Read and print the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var taxa []gbif.TaxonInnerGBIF
	if err := json.Unmarshal(body, &taxa); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if len(taxa) == 0 {
		return nil, fmt.Errorf("no taxa found for name: %s", name)
	}

	return &taxa[0], nil
}

func seedTaxonomyGroup(db geltypes.Tx, group string) error {
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription(fmt.Sprintf("Importing %s taxonomy from GBIF", group)),
		progressbar.OptionSetMaxDetailRow(2),
		progressbar.OptionShowCount(),
	)
	taxon, err := FetchSpeciesSuggestions(group)
	if err != nil {
		return fmt.Errorf("failed to fetch taxon by name: %w", err)
	}
	_ = bar.AddDetail(fmt.Sprintf("Taxon: %s (%d)", taxon.Name, taxon.Key))
	var total int
	err = gbif.ImportTaxonTx(db,
		gbif.ImportRequestGBIF{Key: taxon.Key, Children: true},
		func(p *gbif.ImportProcess) {
			total = p.Imported
			_ = bar.Set(p.Imported)
			if p.Error != nil {
				bar.Describe(fmt.Sprintf("%+v", p))
				logrus.Fatalf("Failed to import taxonomy: %v", p.Error)
			}
		})
	_ = bar.AddDetail(fmt.Sprintf("Taxonomy setup: %d taxa imported", total))
	_ = bar.Close()

	return err
}
