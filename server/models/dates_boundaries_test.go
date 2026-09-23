package models

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDatePrecisionCalendarBounds(t *testing.T) {
	for _, tt := range []struct{ input, upper string }{
		{"2024", "2024-12-31"}, {"2024-02", "2024-02-29"}, {"2023-02", "2023-02-28"}, {"2024-02-15", "2024-02-15"},
	} {
		t.Run(tt.input, func(t *testing.T) {
			date, err := ParseDateWithPrecision(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.input, date.String())
			upper := date.UpperBound()
			require.Equal(t, tt.upper, upper.Date.Format("2006-01-02"))
			require.Equal(t, date.Precision, upper.Precision)
		})
	}
	for _, input := range []string{"2023-02-29", "2024-13", "2024-04-31"} {
		_, err := ParseDateWithPrecision(input)
		require.Error(t, err, input)
	}
	date, err := ParseDateWithPrecision(" ")
	require.NoError(t, err)
	require.Nil(t, date)
}
