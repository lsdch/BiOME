package csvmodels

import (
	"testing"

	"github.com/lsdch/biome/models"
	"github.com/stretchr/testify/require"
)

func TestSamplingCSVHeaderMatchesRecords(t *testing.T) {
	sampling := SamplingWithOccurrencesCSV{SamplingWithOccurrences: models.SamplingWithOccurrences{
		Occurrences: []models.BaseOccurrence{{
			Code:     "occ-1",
			Quantity: models.NewOptional(models.OccurrenceQuantity{Exact: 12}),
			Comments: models.NewOptional("observed"),
		}},
	}}
	header := sampling.Header()
	count := 0
	for record := range sampling.Records('|') {
		count++
		require.Len(t, record, len(header), "every CSV value needs a matching header")
		fields := make(map[string]string)
		for i, name := range header {
			fields[name] = record[i]
		}
		require.Equal(t, "occ-1", fields["occurrence_code"])
		require.Equal(t, "12", fields["specimen_quantity"])
		require.Equal(t, "observed", fields["occurrence_comments"])
	}
	require.Equal(t, 1, count)
}
