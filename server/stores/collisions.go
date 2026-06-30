package stores

import (
	"context"

	"github.com/lsdch/biome/db/biomedb"
)

type CollisionStore struct {
	q *biomedb.Queries
}

func NewCollisionStore(q *biomedb.Queries) *CollisionStore {
	return &CollisionStore{q: q}
}

type CollisionDetectionParams struct {
	RadiusMeters     int32
	DateIntervalDays int32
}

func (s *CollisionStore) DetectBatchSamplingCollisions(ctx context.Context, importHash string, params CollisionDetectionParams) (existing []biomedb.DetectBatchSamplingCollisionsRow, err error) {
	return s.q.DetectBatchSamplingCollisions(ctx, biomedb.DetectBatchSamplingCollisionsParams{
		ImportHash:       importHash,
		RadiusMeters:     params.RadiusMeters,
		DateIntervalDays: params.DateIntervalDays,
	})
}

func (s *CollisionStore) DetectBatchOccurrencesCollisions(ctx context.Context, importHash string, params CollisionDetectionParams) (existing []biomedb.DetectBatchOccurrenceCollisionsRow, err error) {
	return s.q.DetectBatchOccurrenceCollisions(ctx, biomedb.DetectBatchOccurrenceCollisionsParams{
		ImportHash:       importHash,
		RadiusMeters:     params.RadiusMeters,
		DateIntervalDays: params.DateIntervalDays,
	})
}
