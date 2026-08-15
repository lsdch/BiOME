package models

import (
	"testing"

	"github.com/lsdch/biome/db/biomedb"
)

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
