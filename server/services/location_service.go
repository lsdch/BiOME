package services

import (
	"context"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
)

type LocationService struct {
}

func NewLocationService() *LocationService {
	return &LocationService{}
}

func (s *LocationService) ListCountries(ctx context.Context, q db.Querier) ([]models.Country, error) {
	countries, err := q.Queries().ListCountries(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]models.Country, len(countries))
	for i, c := range countries {
		result[i] = models.CountryFromDB(c)
	}
	return result, nil
}

func (s *LocationService) ListCountriesSummary(ctx context.Context, q db.Querier) ([]models.CountrySummary, error) {
	countries, err := q.Queries().ListCountriesSummary(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]models.CountrySummary, len(countries))
	for i, c := range countries {
		result[i] = models.CountrySummaryFromDB(c)
	}
	return result, nil
}

func (s *LocationService) CoordinatesToCountry(ctx context.Context, q db.Querier, latitude, longitude float64) (*models.Country, error) {
	countryDB, err := q.Queries().CoordinatesToCountry(ctx, longitude, latitude)
	if err != nil {
		return nil, err
	}
	country := models.CountryFromDB(countryDB)
	return &country, nil
}
