package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCollectionArrayInputUnmarshalCSV(t *testing.T) {
	var got CollectionArrayInput
	require.NoError(t, got.UnmarshalCSV([]byte("Museum[A1;B2]|University")))
	require.Equal(t, CollectionArrayInput{
		{Name: "Museum", Vouchers: []string{"A1", "B2"}},
		{Name: "University", Vouchers: []string{}},
	}, got)

	for _, input := range []string{"Museum[", "Museum[A1", "Museum[A1]]", "University|Museum["} {
		t.Run(input, func(t *testing.T) {
			var got CollectionArrayInput
			require.NotPanics(t, func() {
				require.Error(t, got.UnmarshalCSV([]byte(input)))
			})
		})
	}
}
