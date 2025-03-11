package taxonomy_test

import (
	"strings"
	"testing"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/taxonomy"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaxonomyList(t *testing.T) {
	cases := []struct {
		filters taxonomy.ListFilters
		expect  func(taxa []taxonomy.TaxonWithParentRef)
	}{
		{
			taxonomy.ListFilters{}, func(taxa []taxonomy.TaxonWithParentRef) {
				assert.NotEmpty(t, taxa)
			},
		},
		{
			taxonomy.ListFilters{IsAnchor: geltypes.NewOptionalBool(false)},
			func(taxa []taxonomy.TaxonWithParentRef) {
				for _, taxon := range taxa {
					assert.False(t, taxon.Anchor)
				}
			},
		},
		{
			taxonomy.ListFilters{Pattern: "Asel"},
			func(taxa []taxonomy.TaxonWithParentRef) {
				assert.NotEmpty(t, taxa)
				assert.Contains(t, strings.ToLower(taxa[0].Name), "asel")
			},
		},
		{
			taxonomy.ListFilters{Ranks: []taxonomy.TaxonRank{taxonomy.Genus}},
			func(taxa []taxonomy.TaxonWithParentRef) {
				assert.NotEmpty(t, taxa)
				for _, taxon := range taxa {
					assert.Equal(t, taxon.Rank, taxonomy.Genus)
				}
			},
		},
		{
			taxonomy.ListFilters{Status: taxonomy.Accepted},
			func(taxa []taxonomy.TaxonWithParentRef) {
				assert.NotEmpty(t, taxa)
				for _, taxon := range taxa {
					assert.Equal(t, taxon.Status, taxonomy.Accepted)
				}
			},
		},
	}

	for _, c := range cases {
		taxa, err := taxonomy.ListTaxa(db.Client(), c.filters)
		require.NoError(t, err)
		c.expect(taxa)
	}
}

func TestTaxonomyFind(t *testing.T) {
	taxon, err := taxonomy.FindByCode(db.Client(), "Asellus")
	require.NoError(t, err)
	assert.Equal(t, taxon.Code, "Asellus")
}
