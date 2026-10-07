package stores

import (
	"strings"
	"testing"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/lsdch/biome/db/biomedb/biomedb/public/table"
	"github.com/stretchr/testify/require"
)

func TestOccurrenceHabitatAndAccessPointsFilters(t *testing.T) {
	const habitatID = "11111111-1111-1111-1111-111111111111"
	for _, tt := range []struct {
		name   string
		params ListOccurrencesParams
	}{
		{name: "absent"},
		{name: "empty", params: ListOccurrencesParams{Habitats: []string{}, AccessPoints: []string{}}},
		{name: "habitat", params: ListOccurrencesParams{Habitats: []string{habitatID}}},
		{name: "all habitats", params: ListOccurrencesParams{Habitats: []string{habitatID, "22222222-2222-2222-2222-222222222222"}}},
		{name: "duplicate habitat", params: ListOccurrencesParams{Habitats: []string{habitatID, habitatID}}},
		{name: "access points", params: ListOccurrencesParams{AccessPoints: []string{"Cave", "O'Brien"}}},
		{name: "combined", params: ListOccurrencesParams{Habitats: []string{habitatID}, AccessPoints: []string{"Cave", "O'Brien"}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stmt := SELECT(table.Occurrences.ID).FROM(table.Occurrences.
				INNER_JOIN(table.Samplings, table.Occurrences.SamplingID.EQ(table.Samplings.ID)))
			sql, args := tt.params.ApplyFilters(stmt).Sql()
			if len(tt.params.Habitats) > 0 {
				require.Contains(t, sql, "EXISTS")
				require.Contains(t, sql, "samplings_habitats.sampling_id = samplings.id")
				require.Equal(t, len(tt.params.Habitats), strings.Count(sql, "EXISTS"))
				require.Equal(t, len(tt.params.Habitats), strings.Count(sql, "samplings_habitats.habitat_id ="))
				require.NotContains(t, sql, " OR ")
				for _, id := range tt.params.Habitats {
					require.Contains(t, args, id)
				}
				require.Contains(t, sql, "::uuid")
				require.Contains(t, args, habitatID)
				require.NotContains(t, sql, habitatID)
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
