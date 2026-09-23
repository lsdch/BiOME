package csvmodels

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/stretchr/testify/require"
)

func TestValidateRowsPreservesErrors(t *testing.T) {
	v := validator.New()
	valid := OccurrenceImportRow{TaxonName: "Asellus aquaticus"}
	valid.Coordinates.Latitude, valid.Coordinates.Longitude = 48, 2
	require.NoError(t, ValidateRows([]OccurrenceImportRow{valid}, v))

	invalidName := valid
	invalidName.TaxonName = ""
	invalidSubspecies := valid
	rank := biomedb.TaxonRankSubspecies
	invalidSubspecies.TaxonRank = &rank
	invalidSubspecies.SetRowNumber(3)

	err := ValidateRows([]OccurrenceImportRow{invalidName, invalidSubspecies}, v)
	var batch *BatchValidationErrors
	require.ErrorAs(t, err, &batch)
	require.Len(t, batch.Errors, 2)
	var fieldErrors validator.ValidationErrors
	require.ErrorAs(t, batch.Errors[0], &fieldErrors)
	var parseError *CSVParseError
	require.ErrorAs(t, batch.Errors[1], &parseError)
	require.Equal(t, int32(3), parseError.RowNumber)
}
