package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
)

func (s *SamplingService) InitHabitatResolution(ctx context.Context, tx *db.Tx, importID uuid.UUID) error {
	return s.store.InitHabitatResolution(ctx, tx, importID)
}

func (s *SamplingService) GetHabitatsResolution(ctx context.Context, q db.Querier, importID uuid.UUID) ([]models.SamplingHabitatResolution, error) {
	return s.store.GetHabitatsResolution(ctx, q, importID)
}

func (s *SamplingService) ResolveHabitats(ctx context.Context, q db.Querier, importID uuid.UUID, input models.SamplingHabitatResolutionInput) (models.SamplingHabitatResolution, error) {
	if err := input.Validate(); err != nil {
		return models.SamplingHabitatResolution{}, err
	}
	return s.store.ResolveHabitats(ctx, q, importID, input)
}
