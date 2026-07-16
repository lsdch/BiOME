package stores

import (
	"context"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
)

func (s *SamplingStore) ListUnknownFixativeCodes(
	ctx context.Context,
	q db.Querier,
	codes []string,
) ([]string, error) {
	return q.Queries().ListUnknownFixativeCodes(ctx, codes)
}

func (s *SamplingStore) ReplaceFixativesAtSampling(
	ctx context.Context,
	q db.Querier,
	samplingID uuid.UUID,
	codes []string,
) error {
	return q.Queries().ReplaceFixativesAtSampling(ctx, samplingID, codes)
}

func (s *SamplingStore) GetSamplingFixativesAtEvent(
	ctx context.Context,
	q db.Querier,
	samplingID uuid.UUID,
) ([]models.Fixative, error) {

	rows, err := q.Queries().GetSamplingFixativesAtEvent(ctx, samplingID)
	if err != nil {
		return nil, err
	}

	res := make([]models.Fixative, len(rows))
	for i, r := range rows {
		res[i] = models.FixativeFromDB(r)
	}

	return res, nil
}
