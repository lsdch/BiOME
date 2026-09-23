package models

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSamplingCodePreservesDatePrecision(t *testing.T) {
	sampling := Sampling{Site: Site{Code: NewOptional("site-1")}}
	require.Equal(t, "site-1|NA", sampling.Code())
	sampling.PerformedOn = NewOptional(DateWithPrecision{Date: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), Precision: EventDatePrecision("month")})
	require.Equal(t, "site-1|2024-02", sampling.Code())
	sampling.Site.Code.Clear()
	sampling.Coordinates.Coordinates = Coordinates{Latitude: -12.5, Longitude: 45.25}
	require.Equal(t, "-12.500000N45.250000E|2024-02", sampling.Code())
}

func TestProximityParamsPreserveExplicitZero(t *testing.T) {
	input := ListSamplingsAtProximityInput{Latitude: 0, Longitude: 2, RadiusMeters: 500, ExcludeIds: []uuid.UUID{uuid.New()}}
	params := input.ToParams()
	require.False(t, params.EventDate.Valid)
	require.Equal(t, int32(30), params.DateIntervalDays)
	input.DateIntervalDays = NewOptional(int32(0))
	input.EventDate = NewOptional(CompositeDate{Year: 2024, Month: 2, Day: 29})
	params = input.ToParams()
	require.Zero(t, params.DateIntervalDays)
	require.Equal(t, "2024-02-29", params.EventDate.Time.Format("2006-01-02"))
	require.Equal(t, input.ExcludeIds, params.ExcludeSamplingIds)
	require.Equal(t, input.ExcludeIds, input.ToParamsH3().ExcludeSamplingIds)
}

func TestHabitatGroupUpdateDistinguishesOmittedNullAndFalse(t *testing.T) {
	input := HabitatGroupUpdate{}
	require.False(t, input.HasUpdateInfos())
	omitted := input.ToDBParams(uuid.New())
	require.False(t, omitted.SetParentHabitatID)
	require.Nil(t, omitted.ExclusiveElements)
	input.Depends.SetNull()
	input.Exclusive = NewOptional(false)
	params := input.ToDBParams(omitted.GroupID)
	require.True(t, input.HasUpdateInfos())
	require.True(t, params.SetParentHabitatID)
	require.False(t, params.ParentHabitatID.Valid)
	require.NotNil(t, params.ExclusiveElements)
	require.False(t, *params.ExclusiveElements)
}
