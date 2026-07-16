package services

import (
	"context"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/types"
)

func (s *OccurrencesService) AddCollection(ctx context.Context, q db.Querier, occurrenceID types.ULID, input models.CollectionInput) (*models.Collection, error) {
	created, err := q.Queries().AddOccurrenceCollection(ctx, input.ToDBParams(occurrenceID))
	if err != nil {
		return nil, err
	}

	collection := models.CollectionFromDB(created)
	return &collection, nil
}

func (s *OccurrencesService) loadOccurrenceCollections(ctx context.Context, q db.Querier, occurrenceID types.ULID) ([]models.Collection, error) {
	collections, err := q.Queries().GetOccurrenceCollections(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}
	result := make([]models.Collection, len(collections))
	for i, c := range collections {
		result[i] = models.CollectionFromDB(c)
	}
	return result, nil
}

func (s *OccurrencesService) ListCollectionNames(ctx context.Context, q db.Querier) ([]string, error) {
	return q.Queries().ListCollectionNames(ctx)
}
