package models

import (
	"testing"

	"github.com/lsdch/biome/db/biomedb"
	"github.com/stretchr/testify/require"
)

func TestTaxonGBIFToStagingPreservesClassification(t *testing.T) {
	taxon := TaxonGBIF{HigherClassificationMap: map[int32]string{
		1: "Animalia", 54: "Arthropoda", 216: "Insecta", 1457: "Coleoptera",
	}}
	// Check associations without relying on Go map iteration order.
	for range 32 {
		params := taxon.ToStaging()
		require.Len(t, params.HigherTaxonKeys, len(taxon.HigherClassificationMap))
		require.Len(t, params.HigherTaxonNames, len(params.HigherTaxonKeys))
		classification := make(map[int32]string)
		for i, key := range params.HigherTaxonKeys {
			classification[key] = params.HigherTaxonNames[i]
		}
		require.Equal(t, taxon.HigherClassificationMap, classification)
	}
}

func TestTaxonGBIFStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   TaxonStatus
	}{
		{"accepted", "ACCEPTED", biomedb.TaxonStatusAccepted},
		{"doubtful", "DOUBTFUL", biomedb.TaxonStatusDoubtful},
		{"other status", "other status", biomedb.TaxonStatusUnclassified},
		{"synonym", "SYNONYM", biomedb.TaxonStatusSynonym},
		{"synonym lowercase", "synonym", biomedb.TaxonStatusSynonym},
		{"synonym mixed case", "SyNoNyM", biomedb.TaxonStatusSynonym},
		{"synonym with extra text", "SYNONYM (some extra text)", biomedb.TaxonStatusSynonym},
		{"synonym with extra text lowercase", "synonym (some extra text)", biomedb.TaxonStatusSynonym},
		{"synonym with extra text mixed case", "SyNoNyM (some extra text)", biomedb.TaxonStatusSynonym},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taxon := TaxonGBIF{Status: tt.status}
			if got := taxon.GetStatus(); got != tt.want {
				t.Errorf("TaxonGBIF.GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
