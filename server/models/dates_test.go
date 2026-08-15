package models

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/stretchr/testify/assert"
)

func TestDateWithPrecision_ToCode(t *testing.T) {
	tests := []struct {
		name      string
		date      time.Time
		precision EventDatePrecision
		expected  string
	}{
		{
			name:      "year precision",
			date:      time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			precision: biomedb.EventDatePrecisionYear,
			expected:  "2023",
		},
		{
			name:      "month precision",
			date:      time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			precision: biomedb.EventDatePrecisionMonth,
			expected:  "2023-05",
		},
		{
			name:      "day precision (default)",
			date:      time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			precision: biomedb.EventDatePrecisionDay,
			expected:  "2023-05-15",
		},
		{
			name:      "unknown precision defaults to day",
			date:      time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			precision: EventDatePrecision("invalid"),
			expected:  "2023-05-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DateWithPrecision{
				Date:      tt.date,
				Precision: tt.precision,
			}
			assert.Equal(t, tt.expected, d.String())
		})
	}
}

func TestMaybeDateWithPrecisionFromDB(t *testing.T) {
	t.Run("should return unset optional when date is invalid", func(t *testing.T) {
		date := pgtype.Date{Valid: false}
		precision := biomedb.EventDatePrecisionDay
		result := MaybeDateWithPrecisionFromDB(date, &precision)
		assert.False(t, result.IsSet)
	})

	t.Run("should return set optional when both are valid", func(t *testing.T) {
		now := time.Now()
		date := pgtype.Date{Valid: true, Time: now}
		precision := biomedb.EventDatePrecisionMonth
		result := MaybeDateWithPrecisionFromDB(date, &precision)
		assert.True(t, result.IsSet)
		assert.Equal(t, now, result.Value.Date)
		assert.Equal(t, biomedb.EventDatePrecisionMonth, result.Value.Precision)
	})
}

func TestCompositeDate(t *testing.T) {
	t.Run("ToTime", func(t *testing.T) {
		d := CompositeDate{
			Day:   15,
			Month: 5,
			Year:  2023,
		}
		expected := time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)
		assert.Equal(t, expected, d.ToTime())
	})

	t.Run("ToPgDate", func(t *testing.T) {
		d := CompositeDate{
			Day:   15,
			Month: 5,
			Year:  2023,
		}
		pgDate := d.ToPgDate()
		assert.True(t, pgDate.Valid)
		assert.Equal(t, time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), pgDate.Time)
	})

}

func TestEventDateInput(t *testing.T) {
	now := time.Now()

	t.Run("PrecisionPtr", func(t *testing.T) {
		t.Run("should return nil when unset", func(t *testing.T) {
			input := &EventDateInput{IsSet: false}
			assert.Nil(t, input.PrecisionPtr())
		})

		t.Run("should return pointer when set", func(t *testing.T) {
			input := &EventDateInput{
				Precision: biomedb.EventDatePrecisionDay,
				IsSet:     true,
			}
			ptr := input.PrecisionPtr()
			assert.NotNil(t, ptr)
			assert.Equal(t, biomedb.EventDatePrecisionDay, *ptr)
		})
	})

	t.Run("DatePg", func(t *testing.T) {
		t.Run("should return invalid pgtype.Date when unset", func(t *testing.T) {
			input := &EventDateInput{IsSet: false}
			pgDate := input.DatePg()
			assert.False(t, pgDate.Valid)
		})

		t.Run("should return valid pgtype.Date when set", func(t *testing.T) {
			input := &EventDateInput{
				Date:  now,
				IsSet: true,
			}
			pgDate := input.DatePg()
			assert.True(t, pgDate.Valid)
			assert.Equal(t, now, pgDate.Time)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Run("should return empty string when unset", func(t *testing.T) {
			input := &EventDateInput{IsSet: false}
			assert.Equal(t, "", input.String())
		})

		t.Run("should format by precision", func(t *testing.T) {
			tests := []struct {
				name      string
				precision EventDatePrecision
				expected  string
			}{
				{"day", biomedb.EventDatePrecisionDay, "2023-05-15"},
				{"month", biomedb.EventDatePrecisionMonth, "2023-05"},
				{"year", biomedb.EventDatePrecisionYear, "2023"},
				{"unknown", EventDatePrecision("invalid"), "2023-05-15"},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					input := &EventDateInput{
						Date:      time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
						Precision: tt.precision,
						IsSet:     true,
					}
					assert.Equal(t, tt.expected, input.String())
				})
			}
		})
	})

	t.Run("UnmarshalCSV", func(t *testing.T) {
		t.Run("should handle empty input", func(t *testing.T) {
			input := &EventDateInput{}
			err := input.UnmarshalCSV([]byte(""))
			assert.NoError(t, err)
			assert.False(t, input.IsSet)
		})

		t.Run("should parse day precision", func(t *testing.T) {
			input := &EventDateInput{}
			err := input.UnmarshalCSV([]byte("2023-05-15"))
			assert.NoError(t, err)
			assert.True(t, input.IsSet)
			assert.Equal(t, biomedb.EventDatePrecisionDay, input.Precision)
			assert.Equal(t, time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), input.Date)
		})

		t.Run("should parse month precision", func(t *testing.T) {
			input := &EventDateInput{}
			err := input.UnmarshalCSV([]byte("2023-05"))
			assert.NoError(t, err)
			assert.True(t, input.IsSet)
			assert.Equal(t, biomedb.EventDatePrecisionMonth, input.Precision)
			assert.Equal(t, time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC), input.Date)
		})

		t.Run("should parse year precision", func(t *testing.T) {
			input := &EventDateInput{}
			err := input.UnmarshalCSV([]byte("2023"))
			assert.NoError(t, err)
			assert.True(t, input.IsSet)
			assert.Equal(t, biomedb.EventDatePrecisionYear, input.Precision)
			assert.Equal(t, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), input.Date)
		})

		t.Run("should return error on invalid format", func(t *testing.T) {
			input := &EventDateInput{}
			err := input.UnmarshalCSV([]byte("invalid"))
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid date format")
		})

		t.Run("should trim whitespace", func(t *testing.T) {
			input := &EventDateInput{}
			err := input.UnmarshalCSV([]byte("  2023-05-15  "))
			assert.NoError(t, err)
			assert.True(t, input.IsSet)
		})
	})
}
