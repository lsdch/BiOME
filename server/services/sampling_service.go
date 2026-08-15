package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/stores"
	"github.com/lsdch/biome/types"
)

type SamplingService struct {
	store *stores.SamplingStore
}

func NewSamplingService() *SamplingService {
	return &SamplingService{
		store: stores.NewSamplingStore(),
	}
}

func (s *SamplingService) ListSamplingsAtProximity(ctx context.Context, q db.Querier, input models.ListSamplingsAtProximityInput) ([]models.SamplingWithDistance, error) {
	return s.store.ListSamplingsAtProximity(ctx, q, input)
}

func (s *SamplingService) ListSamplingsH3AtProximity(ctx context.Context, q db.Querier, input models.ListSamplingsAtProximityInput) ([]models.H3CellWithRichnessAndDistance, error) {
	return s.store.ListSamplingsH3AtProximity(ctx, q, input)
}

func (s *SamplingService) CreateSampling(ctx context.Context, tx *db.Tx, input models.SamplingInput) (*models.SamplingWithDetails, error) {
	sampling, err := s.store.CreateSampling(ctx, tx, input)
	if err != nil {
		return nil, err
	}

	var samplingDetails = models.SamplingMetadata{}

	if len(input.Methods) > 0 {
		err = s.SetMethodsAtSampling(ctx, tx, sampling.ID, input.Methods)
		if err != nil {
			return nil, err
		}
		samplingDetails.SamplingMethods, err = s.LoadSamplingMethods(ctx, tx, sampling.ID)
		if err != nil {
			return nil, err
		}
	}

	if len(input.Fixatives) > 0 {
		err = s.SetFixativesAtSampling(ctx, tx, sampling.ID, input.Fixatives)
		if err != nil {
			return nil, err
		}
		samplingDetails.Fixatives, err = s.LoadSamplingFixatives(ctx, tx, sampling.ID)
		if err != nil {
			return nil, err
		}
	}

	if len(input.TargetTaxa) > 0 {
		err = s.ReplaceSamplingTargetTaxa(ctx, tx, sampling.ID, input.TargetTaxa)
		if err != nil {
			return nil, err
		}
		samplingDetails.TargetTaxa, err = s.LoadSamplingTargetTaxa(ctx, tx, sampling.ID)
		if err != nil {
			return nil, err
		}
	}

	res := sampling.WithDetails(samplingDetails)
	return &res, nil
}

func (s *SamplingService) GetSampling(ctx context.Context, q db.Querier, samplingID uuid.UUID) (*models.Sampling, error) {
	return s.store.GetSampling(ctx, q, samplingID)
}

func (s *SamplingService) GetSamplingBatchWithOccurrences(ctx context.Context, q db.Querier, samplingIDs []uuid.UUID, occurrenceIDs []types.ULID) (map[uuid.UUID]*models.SamplingWithOccurrences, error) {
	sBatch, err := s.store.GetSamplingBatch(ctx, q, samplingIDs)
	if err != nil {
		return nil, err
	}
	oBatch, err := s.store.GetOccurrencesAtSamplingsBatch(ctx, q, samplingIDs, occurrenceIDs)
	if err != nil {
		return nil, err
	}
	res := make(map[uuid.UUID]*models.SamplingWithOccurrences, len(sBatch))
	for _, sampling := range sBatch {
		res[sampling.ID] = &models.SamplingWithOccurrences{
			Sampling:    sampling,
			Occurrences: oBatch[sampling.ID],
		}
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

func (s *SamplingService) ReplaceSamplingTargetTaxa(ctx context.Context, q db.Querier, samplingID uuid.UUID, taxonIDs []uuid.UUID) error {
	return q.Queries().ReplaceSamplingTargetTaxa(ctx, samplingID, taxonIDs)
}

func (s *SamplingService) LoadSamplingMethods(ctx context.Context, q db.Querier, samplingID uuid.UUID) ([]models.SamplingMethod, error) {
	return s.store.GetSamplingMethodsAtEvent(ctx, q, samplingID)

}

func (s *SamplingService) LoadSamplingFixatives(ctx context.Context, q db.Querier, samplingID uuid.UUID) ([]models.Fixative, error) {
	return s.store.GetSamplingFixativesAtEvent(ctx, q, samplingID)
}

func (s *SamplingService) LoadSamplingHabitats(ctx context.Context, q db.Querier, samplingID uuid.UUID) ([]models.HabitatWithGroupName, error) {
	return s.store.GetHabitatsAtEvent(ctx, q, samplingID)
}

func (s *SamplingService) LoadSamplingTargetTaxa(ctx context.Context, q db.Querier, samplingID uuid.UUID) ([]models.Taxon, error) {
	return s.store.GetSamplingTargetTaxa(ctx, q, samplingID)
}

func (s *SamplingService) LoadSamplingMetadata(ctx context.Context, q db.Querier, samplingID uuid.UUID) (models.SamplingMetadata, error) {
	methods, err := s.LoadSamplingMethods(ctx, q, samplingID)
	if err != nil {
		return models.SamplingMetadata{}, err
	}
	fixatives, err := s.LoadSamplingFixatives(ctx, q, samplingID)
	if err != nil {
		return models.SamplingMetadata{}, err
	}
	habitats, err := s.LoadSamplingHabitats(ctx, q, samplingID)
	if err != nil {
		return models.SamplingMetadata{}, err
	}
	taxa, err := s.LoadSamplingTargetTaxa(ctx, q, samplingID)
	if err != nil {
		return models.SamplingMetadata{}, err
	}
	return models.SamplingMetadata{
		SamplingMethods: methods,
		Fixatives:       fixatives,
		Habitats:        habitats,
		TargetTaxa:      taxa,
	}, nil
}

func (s *SamplingService) ListSamplingAccessPoints(ctx context.Context, q db.Querier) ([]string, error) {
	return q.Queries().ListSamplingAccessPoints(ctx)
}

/*
Materializes samplings with their associated metadata for a given batch.
*/
func (s *SamplingService) MaterializeSamplings(ctx context.Context, tx *db.Tx, importID uuid.UUID) error {
	if err := s.store.MaterializeSamplings(ctx, tx, importID); err != nil {
		return err
	}
	if err := s.store.MaterializeSamplingMethods(ctx, tx, importID); err != nil {
		return err
	}
	if err := s.store.MaterializeSamplingFixatives(ctx, tx, importID); err != nil {
		return err
	}
	if err := s.store.MaterializeSamplingTargets(ctx, tx, importID); err != nil {
		return err
	}
	return nil
}
