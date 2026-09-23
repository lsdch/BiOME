package csvmodels

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/lsdch/biome/models"
	"github.com/stretchr/testify/require"
)

func TestSamplingHashUndated(t *testing.T) {
	a := OccurrenceImportRow{}
	b := a
	// These row numbers used to collapse to the same Unicode replacement rune.
	a.SetRowNumber(55296)
	b.SetRowNumber(55297)
	require.NotEqual(t, a.SamplingHash(false), b.SamplingHash(false))
	require.Equal(t, a.SamplingHash(true), b.SamplingHash(true))

	require.NoError(t, a.EventDate.UnmarshalCSV([]byte("2024-02")))
	b.EventDate = a.EventDate
	require.Equal(t, a.SamplingHash(false), b.SamplingHash(false))
}

func TestOccurrenceImportRowCoordinateValidation(t *testing.T) {
	v := validator.New()
	for _, tt := range []struct {
		name     string
		lat, lon float64
		valid    bool
	}{
		{"equator", 0, 2, true},
		{"Greenwich", 48, 0, true},
		{"latitude out of range", 91, 2, false},
		{"longitude out of range", 48, -181, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			row := OccurrenceImportRow{TaxonName: "Asellus aquaticus"}
			row.Coordinates.Coordinates = models.Coordinates{Latitude: tt.lat, Longitude: tt.lon}
			if tt.valid {
				require.NoError(t, row.Validate(v))
			} else {
				require.Error(t, row.Validate(v))
			}
		})
	}
}
