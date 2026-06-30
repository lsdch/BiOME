package models

import (
	"encoding/json"

	"github.com/lsdch/biome/db/biomedb"
)

type Country biomedb.Country

func CountryFromDB(c biomedb.ListCountriesRow) Country {
	return Country{
		Name:         c.Name,
		Code:         c.Code,
		Continent:    c.Continent,
		Subcontinent: c.Subcontinent,
	}
}

type CountryInput struct {
	Name         string          `json:"name"`
	Code         string          `json:"code"`
	Continent    string          `json:"continent"`
	Subcontinent string          `json:"subcontinent"`
	GeometryJSON json.RawMessage `json:"geom"`
}

func (c CountryInput) ToParams() biomedb.InsertCountryParams {
	return biomedb.InsertCountryParams{
		Name:         c.Name,
		Code:         c.Code,
		Continent:    c.Continent,
		Subcontinent: c.Subcontinent,
		Geom:         c.GeometryJSON,
	}
}

type CountrySummary struct {
	Country
	SamplingCount   int64 `json:"sampling_count"`
	OccurrenceCount int64 `json:"occurrence_count"`
}

func CountrySummaryFromDB(c biomedb.ListCountriesSummaryRow) CountrySummary {
	return CountrySummary{
		Country: Country{
			Name:         c.Name,
			Code:         c.Code,
			Continent:    c.Continent,
			Subcontinent: c.Subcontinent,
		},
		SamplingCount:   c.SamplingCount,
		OccurrenceCount: c.OccurrenceCount,
	}
}
