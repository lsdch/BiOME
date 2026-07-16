package stores

import (
	"context"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
)

func (s *SamplingStore) ListUnknownSamplingMethodCodes(
	ctx context.Context,
	q db.Querier,
	codes []string,
) ([]string, error) {
	return q.Queries().ListUnknownSamplingMethodCodes(ctx, codes)
}

func (s *SamplingStore) ReplaceMethodsAtSampling(
	ctx context.Context,
	q db.Querier,
	samplingID uuid.UUID,
	codes []string,
) error {
	return q.Queries().ReplaceMethodsAtSampling(ctx, samplingID, codes)
}

func (s *SamplingStore) GetSamplingMethodsAtEvent(
	ctx context.Context,
	q db.Querier,
	samplingID uuid.UUID,
) ([]models.SamplingMethod, error) {

	rows, err := q.Queries().GetSamplingMethodsAtEvent(ctx, samplingID)
	if err != nil {
		return nil, err
	}

	res := make([]models.SamplingMethod, len(rows))
	for i, r := range rows {
		res[i] = models.SamplingMethodFromDB(r)
	}

	return res, nil
}
