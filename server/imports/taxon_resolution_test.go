package imports

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/services/gbif"
	"github.com/sirupsen/logrus"
)

func TestFetchCandidatesFromGBIF(t *testing.T) {
	ctx := context.Background()

	// À adapter avec ton constructeur réel
	gbifClient := gbif.NewClient(config.GBIFConfig{
		UserAgent:          "BiOME test",
		MaxConcurrent:      5,
		BackboneDatasetKey: "d7dddbf4-2cf0-4f39-9b2a-bb099caae36c",
	})

	r := NewTaxonResolutionService(
		gbifClient,
	)

	// Renseigne ici ton input
	toFetch := []models.TaxonResolution{
		{
			ID:              uuid.New(),
			InputName:       "Balkanostenasellus skopljensis skopljensis",
			InputAuthorship: models.NewOptional("(Karaman, 1937) "),
			ScientificName:  "Balkanostenasellus skopljensis skopljensis (Karaman, 1937) ",
			InputRank:       models.NewOptional("SUBSPECIES"),
			// InputRank: ...
		},
	}

	candidates, err := r.FetchCandidatesFromGBIF(
		ctx,
		1, // Animalia kingdom key
		toFetch,
	)

	if err != nil {
		t.Fatalf("FetchCandidatesFromGBIF failed: %v", err)
	}

	for resolutionID, matches := range candidates {
		logrus.Infof("Resolution %s: %d candidates", resolutionID, len(matches))

		for _, candidate := range matches {
			t.Logf(
				"candidate: scientificName=%s rank=%s status=%s key=%d",
				candidate.ScientificName,
				candidate.Rank,
				candidate.Status,
				candidate.Key,
			)
		}
	}

	if len(candidates[toFetch[0].ID]) == 0 {
		t.Errorf("no GBIF candidates returned")
	}
}
