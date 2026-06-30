package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/oklog/ulid/v2"
)

type SamplingService struct {
	q *biomedb.Queries
}

func NewSamplingService(db db.Querier) *SamplingService {
	return &SamplingService{q: db.Queries()}
}

func (s *SamplingService) GetSamplingBatchWithOccurrences(ctx context.Context, samplingIDs []uuid.UUID, occurrenceIDs []ulid.ULID) (map[uuid.UUID]*models.SamplingWithOccurrences, error) {
	sBatch, err := s.q.GetSamplingBatch(ctx, samplingIDs)
	if err != nil {
		return nil, err
	}
	oBatch, err := s.q.GetOccurrencesAtSamplingsBatch(ctx, samplingIDs, occurrenceIDs)
	if err != nil {
		return nil, err
	}
	res := make(map[uuid.UUID]*models.SamplingWithOccurrences, len(sBatch))
	for _, s := range sBatch {
		res[s.Sampling.ID] = &models.SamplingWithOccurrences{
			Sampling: models.NewSamplingFromDB(s.Sampling, s.Country),
			// preallocate with capacity 1, as we expect most samplings to have a single occurrence
			Occurrences: make([]models.BaseOccurrence, 0, 1),
		}
	}
	for _, o := range oBatch {
		if s, ok := res[o.Occurrence.SamplingID]; ok {
			s.Occurrences = append(s.Occurrences, models.BaseOccurrenceFromDB(o.Occurrence, o.Taxon))
		}
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

func (s *SamplingService) LoadSamplingMethods(ctx context.Context, samplingID uuid.UUID) ([]models.SamplingMethod, error) {
	methods, err := s.q.GetSamplingMethodsAtEvent(ctx, samplingID)
	if err != nil {
		return nil, err
	}
	result := make([]models.SamplingMethod, len(methods))
	for i, m := range methods {
		result[i] = models.SamplingMethodFromDB(m)
	}
	return result, nil
}

func (s *SamplingService) LoadSamplingFixatives(ctx context.Context, samplingID uuid.UUID) ([]models.Fixative, error) {
	fixatives, err := s.q.GetSamplingFixativesAtEvent(ctx, samplingID)
	if err != nil {
		return nil, err
	}
	result := make([]models.Fixative, len(fixatives))
	for i, f := range fixatives {
		result[i] = models.FixativeFromDB(f)
	}
	return result, nil
}

func (s *SamplingService) LoadSamplingHabitats(ctx context.Context, samplingID uuid.UUID) ([]models.HabitatWithGroupName, error) {
	habitats, err := s.q.GetHabitatsAtEvent(ctx, samplingID)
	if err != nil {
		return nil, err
	}
	result := make([]models.HabitatWithGroupName, len(habitats))
	for i, h := range habitats {
		result[i] = models.HabitatWithGroupNameFromDB(h.Habitat, h.HabitatGroup.Label)
	}
	return result, nil
}

func (s *SamplingService) LoadSamplingTargetTaxa(ctx context.Context, samplingID uuid.UUID) ([]models.Taxon, error) {
	taxa, err := s.q.GetSamplingTargetTaxa(ctx, samplingID)
	if err != nil {
		return nil, err
	}
	result := make([]models.Taxon, len(taxa))
	for i, t := range taxa {
		result[i] = *models.TaxonFromDB(&t)
	}
	return result, nil
}

func (s *SamplingService) LoadSamplingMetadata(ctx context.Context, samplingID uuid.UUID) (models.SamplingMetadata, error) {
	methods, err := s.LoadSamplingMethods(ctx, samplingID)
	if err != nil {
		return models.SamplingMetadata{}, err
	}
	fixatives, err := s.LoadSamplingFixatives(ctx, samplingID)
	if err != nil {
		return models.SamplingMetadata{}, err
	}
	habitats, err := s.LoadSamplingHabitats(ctx, samplingID)
	if err != nil {
		return models.SamplingMetadata{}, err
	}
	taxa, err := s.LoadSamplingTargetTaxa(ctx, samplingID)
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

func (s *SamplingService) ListSamplingFixatives(ctx context.Context) ([]biomedb.Fixative, error) {
	return s.q.ListFixatives(ctx)
}

func (s *SamplingService) CreateSamplingFixative(ctx context.Context, input biomedb.CreateFixativeParams) (biomedb.Fixative, error) {
	return s.q.CreateFixative(ctx, input)
}

func (s *SamplingService) UpdateSamplingFixative(ctx context.Context, code string, input models.SamplingFixativeUpdateParams) (biomedb.Fixative, error) {
	return s.q.UpdateFixative(ctx, input.ToParams(code))
}

func (s *SamplingService) DeleteSamplingFixative(ctx context.Context, code string) error {
	tag, err := s.q.DeleteFixative(ctx, code)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("fixative with code %s not found", code)
	}
	return nil
}

func (s *SamplingService) ListSamplingMethods(ctx context.Context) ([]biomedb.SamplingMethod, error) {
	return s.q.ListSamplingMethods(ctx)
}

func (s *SamplingService) CreateSamplingMethod(ctx context.Context, input biomedb.CreateSamplingMethodParams) (biomedb.SamplingMethod, error) {
	return s.q.CreateSamplingMethod(ctx, input)
}

func (s *SamplingService) UpdateSamplingMethod(ctx context.Context, code string, input models.SamplingMethodUpdateParams) (biomedb.SamplingMethod, error) {
	return s.q.UpdateSamplingMethod(ctx, input.ToParams(code))
}

func (s *SamplingService) DeleteSamplingMethod(ctx context.Context, code string) error {
	return s.q.DeleteSamplingMethod(ctx, code)
}

// Initializes the sampling methods resolution state for a given import hash. This is typically called when starting a new import workflow.
func (s *SamplingService) InitMethodResolution(ctx context.Context, importHash string) (state []biomedb.SamplingMethodsResolution, err error) {
	resolution, err := s.q.InitMethodsResolution(ctx, importHash)
	if err != nil {
		return nil, err
	}
	return resolution, nil
}

func (s *SamplingService) GetMethodsResolution(ctx context.Context, importHash string) (state []biomedb.SamplingMethodsResolution, err error) {
	return s.q.GetMethodsResolution(ctx, importHash)
}

func (s *SamplingService) ResolveMethod(ctx context.Context, importHash string, input models.SamplingMethodResolutionInput) (biomedb.SamplingMethodsResolution, error) {
	if err := input.Validate(); err != nil {
		return biomedb.SamplingMethodsResolution{}, err
	}
	return s.q.ResolveMethod(ctx, input.ToParams(importHash))
}
