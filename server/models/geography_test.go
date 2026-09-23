package models

import (
	"github.com/stretchr/testify/require"
	"github.com/uber/h3-go/v4"
	"testing"
)

func TestBoundingBoxPreservesCoordinateOrder(t *testing.T) {
	box := BoundingBox{Corners: [4]Coordinates{{Latitude: 48, Longitude: 2}, {Latitude: 49, Longitude: 2}, {Latitude: 49, Longitude: 3}, {Latitude: 48, Longitude: 3}}}
	polygon := box.ToH3Polygon()
	require.Equal(t, h3.GeoLoop{{Lat: 48, Lng: 2}, {Lat: 49, Lng: 2}, {Lat: 49, Lng: 3}, {Lat: 48, Lng: 3}}, polygon.GeoLoop)
	require.Empty(t, polygon.Holes)
}
