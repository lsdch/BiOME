package csvmodels

import (
	"github.com/google/uuid"
	"github.com/lsdch/biome/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPublicationResolutionPreservesExplicitDOIAndVerbatim(t *testing.T) {
	doi := types.DOI("10.1234/example")
	verbatim := "  Original citation [uncertain]  "
	for _, source := range []*string{nil, &verbatim} {
		input := PublicationResolutionInput{DOI: &doi, Verbatim: source, RowNumbers: []int32{2, 7}}
		params := input.ToParams(uuid.New())
		require.Equal(t, &doi, params.DOI)
		require.Equal(t, source, params.Verbatim)
		require.Equal(t, input.RowNumbers, params.RowNumbers)
		require.Nil(t, params.Year)
	}
}
