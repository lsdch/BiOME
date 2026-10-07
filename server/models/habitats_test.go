package models

import (
	"encoding/json"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/lsdch/biome/data"
	"github.com/stretchr/testify/require"
)

func TestHabitatUpdateJSON(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		set        bool
		value      *string
	}{
		{name: "omitted", body: `{}`},
		{name: "cleared", body: `{"facet":null}`, set: true},
		{name: "assigned", body: `{"facet":"Flow"}`, set: true, value: func() *string { s := "Flow"; return &s }()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var update HabitatUpdate
			require.NoError(t, json.Unmarshal([]byte(tc.body), &update))
			params := update.ToDBParams(uuid.New())
			require.Equal(t, tc.set, params.SetFacet)
			require.Equal(t, tc.value, params.Facet)
		})
	}
}

func TestBootstrapHabitatsYAML(t *testing.T) {
	bytes, err := data.DataFS.ReadFile("habitats.yaml")
	require.NoError(t, err)
	var habitats []HabitatInput
	require.NoError(t, yaml.Unmarshal(bytes, &habitats))
	require.NotEmpty(t, habitats)
	names := make(map[string]bool)
	for _, habitat := range habitats {
		require.NotEmpty(t, habitat.Name)
		require.False(t, names[habitat.Name], "duplicate habitat %s", habitat.Name)
		names[habitat.Name] = true
		require.True(t, habitat.Facet.IsSet, habitat.Name)
		require.NotEmpty(t, habitat.Facet.Value)
	}
	require.True(t, names["Shallow subterranean habitat (SSH)"])
}
