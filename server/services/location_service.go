package services

import (
	"context"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
)

type LocationService struct {
	db *db.DB
}

func NewLocationService(db *db.DB) *LocationService {
	return &LocationService{
		db: db,
	}
}

func (s *LocationService) ListCountries(ctx context.Context) ([]models.Country, error) {
	countries, err := s.db.Queries().ListCountries(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]models.Country, len(countries))
	for i, c := range countries {
		result[i] = models.CountryFromDB(c)
	}
	return result, nil
}

func (s *LocationService) ListCountriesSummary(ctx context.Context) ([]models.CountrySummary, error) {
	countries, err := s.db.Queries().ListCountriesSummary(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]models.CountrySummary, len(countries))
	for i, c := range countries {
		result[i] = models.CountrySummaryFromDB(c)
	}
	return result, nil
}
