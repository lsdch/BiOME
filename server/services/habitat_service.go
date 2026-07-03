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

func (s *HabitatService) CreateHabitatGroup(ctx context.Context, group models.HabitatGroupInput) error {
	return s.db.WithTx(ctx, func(q *biomedb.Queries) error {
		created, err := q.InsertHabitatGroup(ctx, group.ToDBParams())
		if err != nil {
			return err
		}
		for _, element := range group.Elements {
			err = s.AddHabitatToGroup(ctx, q, created.ID, element)
			if err != nil {
				return err
			}
		}
		return nil

	})
}

func (s *HabitatService) DeleteHabitatGroup(ctx context.Context, groupID uuid.UUID) error {
	return s.db.Queries().DeleteHabitatGroup(ctx, groupID)
}

func (s *HabitatService) UpdateHabitatGroup(ctx context.Context, groupID uuid.UUID, update models.HabitatGroupUpdate) error {

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

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

	err = tx.Commit(ctx)
	return err
}
