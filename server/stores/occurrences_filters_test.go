package stores

import (
	"strings"
	"testing"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/lsdch/biome/db/biomedb/biomedb/public/table"
	"github.com/stretchr/testify/require"
)

func TestOccurrenceHabitatAndAccessPointsFilters(t *testing.T) {
	const habitatName = "Forest"
	for _, tt := range []struct {
		name   string
		params ListOccurrencesParams
	}{
		{name: "absent"},
		{name: "empty", params: ListOccurrencesParams{Habitats: []string{}, AccessPoints: []string{}}},
		{name: "habitat", params: ListOccurrencesParams{Habitats: []string{habitatName}}},
		{name: "all habitats", params: ListOccurrencesParams{Habitats: []string{habitatName, "O'Brien"}}},
		{name: "duplicate habitat", params: ListOccurrencesParams{Habitats: []string{habitatName, habitatName}}},
		{name: "access points", params: ListOccurrencesParams{AccessPoints: []string{"Cave", "O'Brien"}}},
		{name: "combined", params: ListOccurrencesParams{Habitats: []string{habitatName}, AccessPoints: []string{"Cave", "O'Brien"}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stmt := SELECT(table.Occurrences.ID).FROM(table.Occurrences.
				INNER_JOIN(table.Samplings, table.Occurrences.SamplingID.EQ(table.Samplings.ID)))
			sql, args := tt.params.ApplyFilters(stmt).Sql()
			if len(tt.params.Habitats) > 0 {
				require.Contains(t, sql, "EXISTS")
				require.Contains(t, sql, "samplings_habitats.sampling_id = samplings.id")
				require.Equal(t, len(tt.params.Habitats), strings.Count(sql, "EXISTS"))
				require.Contains(t, sql, "samplings_habitats.habitat_id = habitats.id")
				require.Equal(t, len(tt.params.Habitats), strings.Count(sql, "habitats.label ="))
				require.NotContains(t, sql, " OR ")
				for _, name := range tt.params.Habitats {
					require.Contains(t, args, name)
					require.NotContains(t, sql, name)
				}
				require.NotContains(t, sql, "::uuid")
				require.NotContains(t, args, int64(1))
				require.Contains(t, sql, "SELECT 1")
			} else {
				require.NotContains(t, sql, "samplings_habitats")
			}
			if len(tt.params.AccessPoints) > 0 {
				require.Contains(t, sql, "samplings.access_points && CAST(")
				require.Contains(t, sql, "AS citext[])")
				require.Contains(t, args, tt.params.AccessPoints)
				require.NotContains(t, sql, "O'Brien")
			} else {
				require.NotContains(t, sql, "access_points")
			}
		})
	}
}
