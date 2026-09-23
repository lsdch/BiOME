package models

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestErrorPathComposition(t *testing.T) {
	cause := errors.New("invalid value")
	field := WrapErrorPath(cause, "name")
	nested := WrapErrorPath(WrapErrorIndex(field, 2), "taxa")
	require.Equal(t, "taxa[2].name", nested.Path)
	require.ErrorIs(t, nested, cause)
	require.Equal(t, nested.Path, nested.AsErrorDetail().Location)
	t.Run("wrapped error", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := WrapErrorPath(fmt.Errorf("validation: %w", field), "taxon")
			require.Equal(t, "taxon.name", got.Path)
			require.ErrorIs(t, got, cause)
		})
	})
}
