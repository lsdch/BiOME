package models

import (
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTaxonRelationsPreserveSynonymTarget(t *testing.T) {
	parent, accepted := uuid.New(), uuid.New()
	taxon := TaxonFromDB(&biomedb.Taxon{ID: uuid.New(), Name: "Original name", Status: biomedb.TaxonStatusSynonym, ParentID: UUIDToPg(parent), AcceptedTaxonID: UUIDToPg(accepted)})
	require.Equal(t, parent, taxon.ParentID.Value)
	require.Equal(t, accepted, taxon.AcceptedID.Value)
	require.Equal(t, "Original name", taxon.Name)
	require.True(t, taxon.GBIF_ID.Missing())
	require.Nil(t, TaxonFromDB(nil))
}

func TestTaxonRankTransitions(t *testing.T) {
	for _, tt := range []struct{ rank, parent, child TaxonRank }{
		{biomedb.TaxonRankSpecies, biomedb.TaxonRankGenus, biomedb.TaxonRankSubspecies},
		{biomedb.TaxonRankGenus, biomedb.TaxonRankFamily, biomedb.TaxonRankSpecies},
		{biomedb.TaxonRankKingdom, biomedb.TaxonRankKingdom, biomedb.TaxonRankPhylum},
	} {
		t.Run(string(tt.rank), func(t *testing.T) {
			require.Equal(t, tt.parent, ParentRank(tt.rank))
			require.Equal(t, tt.child, ChildRank(tt.rank))
		})
	}
}
