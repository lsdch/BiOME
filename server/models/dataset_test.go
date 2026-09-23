package models

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDatasetCreationIdentityAndVisibility(t *testing.T) {
	owner := uuid.New()
	for _, public := range []bool{false, true} {
		input := DatasetInput{Label: "Faune des rivières", Public: public}
		params := input.ToParams(owner)
		require.Equal(t, owner, params.OwnerID)
		require.Equal(t, public, params.IsPublic)
		require.Equal(t, "faune-des-rivieres", params.Slug)
		require.Equal(t, input.Label, params.Label)
		require.NotEqual(t, params.ULID, input.ToParams(owner).ULID)
	}
}
