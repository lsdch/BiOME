package models

import (
	"github.com/uber/h3-go/v4"
)

type Coordinates struct {
	Latitude  float32 `json:"latitude" csv:"latitude"`
	Longitude float32 `json:"longitude" csv:"longitude"`
}

func (c Coordinates) ToH3GeoCoord() h3.LatLng {
	return h3.LatLng{
		Lat: float64(c.Latitude),
		Lng: float64(c.Longitude),
	}
}

type CoordinatesWithPrecision struct {
	Coordinates
	Precision Optional[int32] `json:"precision,omitempty"`
}

func CoordinatesWithPrecisionFromDB(lat float32, lon float32, precision *int32) CoordinatesWithPrecision {
	return CoordinatesWithPrecision{
		Coordinates: Coordinates{
			Latitude:  lat,
			Longitude: lon,
		},
		Precision: NewOptionalFromPtr(precision),
	}
}

type CoordinatesWithPrecisionInput struct {
	Coordinates
	PrecisionM *int32 `csv:"coordinates_precision_m,omitempty"`
}

type BoundingBox struct {
	Corners [4]Coordinates `json:",inline" minItems:"4" maxItems:"4"`
}

func (b BoundingBox) ToH3Polygon() h3.GeoPolygon {
	loop := make([]h3.LatLng, len(b.Corners))
	for i, corner := range b.Corners {
		loop[i] = corner.ToH3GeoCoord()
	}
	geoPolygon := h3.GeoPolygon{
		GeoLoop: loop,
	}
	return geoPolygon
}
