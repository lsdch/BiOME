package services

import (
	"context"
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/sirupsen/logrus"
)

type HabitatService struct{}

func NewHabitatService() *HabitatService { return &HabitatService{} }

func (s *HabitatService) GetHabitats(ctx context.Context, q db.Querier) ([]models.Habitat, error) {
	rows, err := q.Queries().ListHabitats(ctx)
	if err != nil {
		return nil, err
	}
	habitats := make([]models.Habitat, len(rows))
	for i, h := range rows {
		habitats[i] = models.HabitatFromDB(h)
	}
	return habitats, nil
}

func (s *HabitatService) CreateHabitat(ctx context.Context, q db.Querier, input models.HabitatInput) (models.Habitat, error) {
	h, err := q.Queries().InsertHabitat(ctx, input.ToDBParams())
	if err != nil {
		return models.Habitat{}, err
	}
	return models.HabitatFromDB(h), nil
}

func (s *HabitatService) DeleteHabitat(ctx context.Context, q db.Querier, id uuid.UUID) error {
	return q.Queries().DeleteHabitatByID(ctx, id)
}

func (s *HabitatService) UpdateHabitat(ctx context.Context, q db.Querier, id uuid.UUID, update models.HabitatUpdate) error {
	return q.Queries().UpdateHabitat(ctx, update.ToDBParams(id))
}

func (s *HabitatService) BootstrapHabitats(ctx context.Context, tx *db.Tx, yamlBytes []byte) error {
	habitats, err := s.GetHabitats(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to list habitats: %w", err)
	}
	if len(habitats) > 0 {
		logrus.Info("Found existing habitats, skipping bootstrap")
		return nil
	}
	var inputs []models.HabitatInput
	if err := yaml.Unmarshal(yamlBytes, &inputs); err != nil {
		return fmt.Errorf("failed to unmarshal habitats: %w", err)
	}
	for _, input := range inputs {
		if _, err := s.CreateHabitat(ctx, tx, input); err != nil {
			return fmt.Errorf("failed to create habitat %s: %w", input.Name, err)
		}
	}
	return nil
}
