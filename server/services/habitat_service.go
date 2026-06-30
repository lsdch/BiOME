package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
)

type HabitatService struct {
	db *db.DB
}

func NewHabitatService(db *db.DB) *HabitatService {
	return &HabitatService{db: db}
}

func (s *HabitatService) GetHabitatGroups(ctx context.Context) ([]models.HabitatGroupWithElements, error) {
	groups, err := s.db.Queries().ListHabitatGroups(ctx)
	if err != nil {
		return nil, err
	}

	habitats, err := s.db.Queries().ListHabitats(ctx)
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

func (s *HabitatService) AddHabitatToGroup(ctx context.Context, groupID uuid.UUID, habitat models.HabitatInput) error {
	_, err := s.db.Queries().InsertHabitatInGroup(ctx, biomedb.InsertHabitatInGroupParams{
		Label:          habitat.Label,
		Description:    habitat.Description.ToPtr(),
		HabitatGroupID: groupID,
	})
	return err
}

func (s *HabitatService) DeleteHabitat(ctx context.Context, habitatName string) error {
	return s.db.Queries().DeleteHabitatByName(ctx, habitatName)
}

func (s *HabitatService) CreateHabitatGroup(ctx context.Context, group models.HabitatGroupInput) error {
	created, err := s.db.Queries().InsertHabitatGroup(ctx, biomedb.InsertHabitatGroupParams{
		Label:           group.Label,
		ParentHabitatID: models.UUIDOpt(group.Depends),
		Description:     group.Description.ToPtr(),
	})
	if err != nil {
		return err
	}
	for _, element := range group.Elements {
		err = s.AddHabitatToGroup(ctx, created.ID, element)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *HabitatService) DeleteHabitatGroup(ctx context.Context, groupName string) error {
	return s.db.Queries().DeleteHabitatGroupByName(ctx, groupName)
}
