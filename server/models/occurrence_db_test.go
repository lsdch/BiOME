package models

import (
	"github.com/lsdch/biome/db/biomedb"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOccurrenceFromDBPreservesVerbatimAndUnknowns(t *testing.T) {
	verbatim := "  Asellus cf. aquaticus [original spelling]  "
	lower, upper := int32(2), int32(5)
	occurrence := BaseOccurrenceFromDB(biomedb.Occurrence{
		VerbatimIdentification: &verbatim, IdentificationConfer: true,
		QuantityLower: &lower, QuantityUpper: &upper,
	}, biomedb.Taxon{Name: "Asellus aquaticus"})
	require.Equal(t, verbatim, occurrence.Identification.Verbatim.Value)
	require.True(t, occurrence.Identification.Confer)
	require.True(t, occurrence.Identification.IdentifiedOn.Missing())
	require.True(t, occurrence.TypeStatus.Missing())
	require.Equal(t, "2-5", occurrence.Quantity.Value.String())
}
