package models

import (
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/lsdch/biome/db/biomedb/biomedb/public/table"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOccurrenceSortingSQL(t *testing.T) {
	for _, tt := range []struct {
		key   OccurrenceSortKey
		order SortOrder
		want  string
	}{
		{OccurrenceSortKeyCode, SortAsc, "occurrences.code ASC NULLS LAST"},
		{OccurrenceSortKeyEventDate, SortDesc, "samplings.event_date DESC NULLS LAST"},
		{OccurrenceSortKeyTaxonName, SortAsc, "taxa.name ASC NULLS LAST"},
	} {
		t.Run(string(tt.key), func(t *testing.T) {
			sort := SortBy[OccurrenceSortKey]{Key: NewOptional(tt.key), Order: tt.order}
			sql, _ := SELECT(table.Occurrences.ID).ORDER_BY(sort.ToOrderByClause()).Sql()
			require.Contains(t, sql, tt.want)
		})
	}
}
