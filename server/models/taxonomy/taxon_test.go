package taxonomy_test

import (
	"darco/proto/db"
	"darco/proto/models/taxonomy"
	"strings"
	"testing"

	"github.com/edgedb/edgedb-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaxonomyList(t *testing.T) {
	cases := []struct {
		filters taxonomy.ListFilters
		expect  func(taxa []taxonomy.TaxonDB)
	}{
		{
			taxonomy.ListFilters{}, func(taxa []taxonomy.TaxonDB) {
				assert.NotEmpty(t, taxa)
			},
		},
		{
			taxonomy.ListFilters{IsAnchor: edgedb.NewOptionalBool(false)},
			func(taxa []taxonomy.TaxonDB) {
				for _, taxon := range taxa {
					assert.False(t, taxon.Anchor)
				}
			},
		},
		{
			taxonomy.ListFilters{Pattern: "thisTAXONdoesntEXIST"},
			func(taxa []taxonomy.TaxonDB) {
				assert.Empty(t, taxa)
			},
		},
		{
			taxonomy.ListFilters{Pattern: "Asel"},
			func(taxa []taxonomy.TaxonDB) {
				assert.NotEmpty(t, taxa)
				for _, taxon := range taxa {
					assert.Contains(t, strings.ToLower(taxon.Name), "asel")
				}
			},
		},
		{
			taxonomy.ListFilters{Rank: taxonomy.Genus},
			func(taxa []taxonomy.TaxonDB) {
				assert.NotEmpty(t, taxa)
				for _, taxon := range taxa {
					assert.Equal(t, taxon.Rank, taxonomy.Genus)
				}
			},
		},
		{
			taxonomy.ListFilters{Status: taxonomy.Accepted},
			func(taxa []taxonomy.TaxonDB) {
				assert.NotEmpty(t, taxa)
				for _, taxon := range taxa {
					assert.Equal(t, taxon.Status, taxonomy.Accepted)
				}
			},
		},
	}

	for _, c := range cases {
		taxa, err := taxonomy.ListTaxa(db.Client(), &c.filters)
		require.NoError(t, err)
		c.expect(taxa)
	}
}

func TestTaxonomyFind(t *testing.T) {
	taxon, err := taxonomy.FindByCode(db.Client(), "Asellus")
	require.NoError(t, err)
	assert.Equal(t, taxon.Code, "Asellus")
}
