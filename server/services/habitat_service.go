package services

import (
	"context"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/sirupsen/logrus"
)

type HabitatService struct{}

func NewHabitatService() *HabitatService {
	return &HabitatService{}
}

func (s *HabitatService) GetHabitatGroups(ctx context.Context, q db.Querier) ([]models.HabitatGroupWithElements, error) {
	groups, err := q.Queries().ListHabitatGroups(ctx)
	if err != nil {
		return nil, err
	}

	habitats, err := q.Queries().ListHabitats(ctx)
	if err != nil {
		return nil, err
	}

	byGroup := make(map[uuid.UUID][]models.Habitat)

	for _, habitat := range habitats {
		byGroup[habitat.HabitatGroupID] = append(
			byGroup[habitat.HabitatGroupID],
			models.HabitatFromDB(habitat),
		)
	}

	result := make([]models.HabitatGroupWithElements, 0, len(groups))

	for _, group := range groups {
		result = append(result, models.HabitatGroupWithElements{
			HabitatGroup: models.HabitatGroupFromDB(group),
			Elements:     byGroup[group.ID],
		})
	}

	return result, nil
}

func (s *HabitatService) AddHabitatToGroup(ctx context.Context, q *biomedb.Queries, groupID uuid.UUID, habitat models.HabitatInput) error {
	_, err := q.InsertHabitatInGroup(ctx, biomedb.InsertHabitatInGroupParams{
		Label:          habitat.Label,
		Description:    habitat.Description.ToPtr(),
		HabitatGroupID: groupID,
	})
	return err
}

func (s *HabitatService) DeleteHabitat(ctx context.Context, q *biomedb.Queries, habitatID uuid.UUID) error {
	return q.DeleteHabitatByID(ctx, habitatID)
}

func (s *HabitatService) CreateHabitatGroup(ctx context.Context, tx *db.Tx, group models.HabitatGroupInput) error {
	created, err := tx.Queries().InsertHabitatGroup(ctx, group.ToDBParams())
	if err != nil {
		return err
	}
	for _, element := range group.Elements {
		if err := s.AddHabitatToGroup(ctx, tx.Queries(), created.ID, element); err != nil {
			return err
		}
	}
	return nil
}

func (s *HabitatService) DeleteHabitatGroup(ctx context.Context, q db.Querier, groupID uuid.UUID) error {
	return q.Queries().DeleteHabitatGroup(ctx, groupID)
}

func (s *HabitatService) UpdateHabitatGroup(ctx context.Context, tx *db.Tx, groupID uuid.UUID, update models.HabitatGroupUpdate) error {
	if update.HasUpdateInfos() {
		if err := tx.Queries().UpdateHabitatGroupInfo(ctx, update.ToDBParams(groupID)); err != nil {
			return err
		}
	}

	for _, newHabitat := range update.CreateElements {
		if err := s.AddHabitatToGroup(ctx, tx.Queries(), groupID, newHabitat); err != nil {
			return err
		}
	}

	for habitatID, habitatUpdate := range update.UpdateElements {
		if err := tx.Queries().UpdateHabitat(ctx, habitatUpdate.ToDBParams(habitatID)); err != nil {
			return err
		}
	}

	for _, habitatID := range update.DeleteElements {
		if err := s.DeleteHabitat(ctx, tx.Queries(), habitatID); err != nil {
			return err
		}
	}

	return nil
}

func (s *HabitatService) BootstrapHabitats(ctx context.Context, tx *db.Tx, yamlBytes []byte) error {
	if groups, err := s.GetHabitatGroups(ctx, tx); err != nil {
		return fmt.Errorf("failed to list habitat groups: %w", err)
	} else if len(groups) > 0 {
		logrus.Infof("Found existing habitat groups, skipping bootstrap")
		return nil
	}

	var groups []models.HabitatGroupInput
	err := yaml.Unmarshal(yamlBytes, &groups)
	if err != nil {
		return fmt.Errorf("failed to unmarshal habitat groups: %w", err)
	}

	logrus.Infof("Bootstrapping %d habitat groups from YAML", len(groups))

	for _, group := range groups {
		err = s.CreateHabitatGroup(ctx, tx, group)
		if err != nil {
			return fmt.Errorf("failed to create habitat group %s: %w", group.Label, err)
		}
	}

	return nil
}
