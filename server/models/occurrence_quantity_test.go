package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuantityInputUnmarshalCSV(t *testing.T) {
	for _, tt := range []struct {
		input string
		want  QuantityInput
	}{
		{"", QuantityInput{}},
		{" 12 ", QuantityInput{Exact: NewOptional(int32(12))}},
		{"2 - 5", QuantityInput{Lower: NewOptional(int32(2)), Upper: NewOptional(int32(5))}},
	} {
		t.Run(tt.input, func(t *testing.T) {
			var got QuantityInput
			require.NoError(t, got.UnmarshalCSV([]byte(tt.input)))
			require.Equal(t, tt.want, got)
		})
	}
	for _, input := range []string{"1.5", "12abc", "2-5abc", "2.5-5", "2147483648", "1-2-3"} {
		t.Run(input, func(t *testing.T) {
			var got QuantityInput
			require.Error(t, got.UnmarshalCSV([]byte(input)))
		})
	}
}
