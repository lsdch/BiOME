package models

import (
	"fmt"

	"github.com/uber/h3-go/v4"
)

type Coordinates struct {
	Latitude  float64 `json:"latitude" csv:"latitude" query:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" csv:"longitude" query:"longitude" validate:"required,longitude"`
}

func (c Coordinates) ToCode() string {
	return fmt.Sprintf("%fN%fE", c.Latitude, c.Longitude)
}

func (c Coordinates) ToH3GeoCoord() h3.LatLng {
	return h3.LatLng{
		Lat: float64(c.Latitude),
		Lng: float64(c.Longitude),
	}
}

type CoordinatesWithPrecision struct {
	Coordinates
	Precision Optional[int32] `json:"precision,omitempty" validate:"gte=0"`
}

func CoordinatesWithPrecisionFromDB(lat float64, lon float64, precision *int32) CoordinatesWithPrecision {
	return CoordinatesWithPrecision{
		Coordinates: Coordinates{
			Latitude:  lat,
			Longitude: lon,
		},
		Precision: NewOptionalFromPtr(precision),
	}
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

type H3Cell struct {
	H3Index          h3.Cell `json:"h3_index"`
	SamplingsCount   int32   `json:"samplings_count"`
	OccurrencesCount int32   `json:"occurrences_count"`
}

func CellH3FromDB(h3Index h3.Cell, samplingsCount int32, occurrencesCount int32) H3Cell {
	return H3Cell{
		H3Index:          h3Index,
		SamplingsCount:   samplingsCount,
		OccurrencesCount: occurrencesCount,
	}
}

func (c H3Cell) WithDistance(distanceMeters int32) H3CellWithDistance {
	return H3CellWithDistance{
		H3Cell:         c,
		DistanceMeters: distanceMeters,
	}
}

func (c H3Cell) WithRichness(speciesRichness int64, genusRichness int64, familyRichness int64) H3CellWithRichness {
	return H3CellWithRichness{
		H3Cell:          c,
		SpeciesRichness: speciesRichness,
		GenusRichness:   genusRichness,
		FamilyRichness:  familyRichness,
	}
}

type H3CellWithDistance struct {
	H3Cell
	DistanceMeters int32 `json:"distance_meters"`
}

func (c H3CellWithDistance) WithRichness(speciesRichness int64, genusRichness int64, familyRichness int64) H3CellWithRichnessAndDistance {
	return H3CellWithRichnessAndDistance{
		H3CellWithRichness: c.H3Cell.WithRichness(speciesRichness, genusRichness, familyRichness),
		DistanceMeters:     c.DistanceMeters,
	}
}

type H3CellWithRichness struct {
	H3Cell
	SpeciesRichness int64 `json:"species_richness"`
	GenusRichness   int64 `json:"genus_richness"`
	FamilyRichness  int64 `json:"family_richness"`
	OccurringTaxa   int64 `json:"occurring_taxa"`
}

func (c H3CellWithRichness) WithDistance(distanceMeters int32) H3CellWithRichnessAndDistance {
	return H3CellWithRichnessAndDistance{
		H3CellWithRichness: c,
		DistanceMeters:     distanceMeters,
	}
}

type H3CellWithRichnessAndDistance struct {
	H3CellWithRichness
	DistanceMeters int32 `json:"distance_meters"`
}
