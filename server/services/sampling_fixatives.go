package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/sirupsen/logrus"
)

func (s *SamplingService) ListSamplingFixatives(ctx context.Context, q db.Querier) ([]models.Fixative, error) {
	return s.store.ListSamplingFixatives(ctx, q)
}

func (s *SamplingService) CreateSamplingFixative(ctx context.Context, q db.Querier, input models.FixativeInput) (models.Fixative, error) {
	return s.store.CreateSamplingFixative(ctx, q, input)
}

func (s *SamplingService) UpdateSamplingFixative(ctx context.Context, q db.Querier, code string, input models.FixativeUpdateParams) (models.Fixative, error) {
	return s.store.UpdateSamplingFixative(ctx, q, code, input)
}

func (s *SamplingService) DeleteSamplingFixative(ctx context.Context, q db.Querier, code string) error {
	return s.store.DeleteSamplingFixative(ctx, q, code)
}

func (s *SamplingService) SetFixativesAtSampling(ctx context.Context, tx *db.Tx, samplingID uuid.UUID, codes []string) error {

	unknownCode, err := s.store.ListUnknownFixativeCodes(ctx, tx, codes)
	if err != nil {
		return err
	}

	if len(unknownCode) > 0 {
		errs := make([]error, 0, len(unknownCode))
		for _, code := range unknownCode {
			errs = append(errs, models.ErrUnknownFixativeCode(code))
		}
		return errors.Join(errs...)
	}

	return s.store.ReplaceFixativesAtSampling(ctx, tx, samplingID, codes)
}

func (s *SamplingService) BootstrapSamplingFixatives(ctx context.Context, q db.Querier, yamlBytes []byte) error {
	if fixatives, err := s.ListSamplingFixatives(ctx, q); err != nil {
		return fmt.Errorf("failed to list sampling fixatives: %w", err)
	} else if len(fixatives) > 0 {
		logrus.Infof("Found %d existing sampling fixatives, skipping bootstrap", len(fixatives))
		return nil
	}

	var fixatives []models.FixativeInput
	err := yaml.Unmarshal(yamlBytes, &fixatives)
	logrus.Infof("Bootstrapping %d sampling fixatives from YAML", len(fixatives))
	if err != nil {
		return fmt.Errorf("failed to unmarshal sampling fixatives: %w", err)
	}
	for _, fixative := range fixatives {
		_, err := s.CreateSamplingFixative(ctx, q, fixative)
		if err != nil {
			return fmt.Errorf("failed to create sampling fixative %s: %w", fixative.Code, err)
		}
	}
	return nil
}

// Initializes the sampling fixatives resolution state for a given import hash. This is typically called when starting a new import workflow.
func (s *SamplingService) InitFixativeResolution(ctx context.Context, q db.Querier, importID uuid.UUID) (state []models.SamplingFixativeResolution, err error) {
	return s.store.InitFixativeResolution(ctx, q, importID)
}

func (s *SamplingService) GetFixativesResolution(ctx context.Context, q db.Querier, importID uuid.UUID) (state []models.SamplingFixativeResolution, err error) {
	return s.store.GetFixativesResolution(ctx, q, importID)
}

func (s *SamplingService) ResolveFixative(ctx context.Context, q db.Querier, importID uuid.UUID, input models.SamplingFixativeResolutionInput) (models.SamplingFixativeResolution, error) {
	if err := input.Validate(); err != nil {
		return models.SamplingFixativeResolution{}, err
	}
	return s.store.ResolveFixative(ctx, q, importID, input)
}
