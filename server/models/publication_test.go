package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPublicationUpdatePreservesPatchIntent(t *testing.T) {
	var input UpdatePublicationParams
	require.NoError(t, json.Unmarshal([]byte(`{"title":null,"authors":[],"comments":"verbatim note"}`), &input))
	params := input.ToDBParams(uuid.New())
	require.True(t, params.SetTitle)
	require.Nil(t, params.Title)
	require.False(t, params.SetDOI)
	require.Nil(t, params.Year)
	require.NotNil(t, params.Authors)
	require.Empty(t, params.Authors)
	require.Equal(t, "verbatim note", *params.Comments)
}

func TestPublicationCreationPreservesVerbatimAndUnknownYear(t *testing.T) {
	input := CreatePublicationParams{Verbatim: "  Smith, circa 1900 [uncertain]  "}
	params := input.ToDBParams()
	require.Equal(t, input.Verbatim, params.Verbatim)
	require.Nil(t, params.Year)
	require.Nil(t, params.DOI)
}
