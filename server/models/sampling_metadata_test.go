package models

import (
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSelectedVocabRequiresCandidate(t *testing.T) {
	for _, tt := range []struct {
		name  string
		id    Optional[uuid.UUID]
		valid bool
	}{
		{"absent", Optional[uuid.UUID]{}, false},
		{"nil UUID", NewOptional(uuid.Nil), false},
		{"candidate", NewOptional(uuid.New()), true},
	} {
		for _, kind := range []string{"method", "fixative"} {
			t.Run(kind+"/"+tt.name, func(t *testing.T) {
				var err error
				if kind == "method" {
					err = (SamplingMethodResolutionInput{Status: biomedb.VocabResolutionStatusSelected, ResolvedMethodId: tt.id}).Validate()
				} else {
					err = (SamplingFixativeResolutionInput{Status: biomedb.VocabResolutionStatusSelected, ResolvedFixativeID: tt.id}).Validate()
				}
				if tt.valid {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
				}
			})
		}
	}
}
