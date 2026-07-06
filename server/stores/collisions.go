package stores

import (
	"context"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
)

type CollisionStore struct {
}

func NewCollisionStore() *CollisionStore {
	return &CollisionStore{}
}

type CollisionDetectionParams struct {
	RadiusMeters     int32
	DateIntervalDays int32
}

func (s *CollisionStore) DetectBatchSamplingCollisions(ctx context.Context, q db.Querier, importHash string, params CollisionDetectionParams) (existing []biomedb.DetectBatchSamplingCollisionsRow, err error) {
	return q.Queries().DetectBatchSamplingCollisions(ctx, biomedb.DetectBatchSamplingCollisionsParams{
		ImportHash:       importHash,
		RadiusMeters:     params.RadiusMeters,
		DateIntervalDays: params.DateIntervalDays,
	})
}

func (s *CollisionStore) DetectBatchOccurrencesCollisions(ctx context.Context, q db.Querier, importHash string, params CollisionDetectionParams) (existing []biomedb.DetectBatchOccurrenceCollisionsRow, err error) {
	return q.Queries().DetectBatchOccurrenceCollisions(ctx, biomedb.DetectBatchOccurrenceCollisionsParams{
		ImportHash:       importHash,
		RadiusMeters:     params.RadiusMeters,
		DateIntervalDays: params.DateIntervalDays,
	})
}
