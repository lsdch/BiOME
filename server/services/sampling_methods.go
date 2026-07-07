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

func (s *SamplingService) ListSamplingMethods(ctx context.Context, q db.Querier) ([]models.SamplingMethod, error) {
	return s.store.ListSamplingMethods(ctx, q)
}

func (s *SamplingService) CreateSamplingMethod(ctx context.Context, q db.Querier, input models.SamplingMethodInput) (models.SamplingMethod, error) {
	return s.store.CreateSamplingMethod(ctx, q, input)
}

func (s *SamplingService) UpdateSamplingMethod(ctx context.Context, q db.Querier, code string, input models.SamplingMethodUpdateParams) (models.SamplingMethod, error) {
	return s.store.UpdateSamplingMethod(ctx, q, code, input)
}

func (s *SamplingService) DeleteSamplingMethod(ctx context.Context, q db.Querier, code string) error {
	return s.store.DeleteSamplingMethod(ctx, q, code)
}

func (s *SamplingService) SetMethodsAtSampling(ctx context.Context, tx *db.Tx, samplingID uuid.UUID, codes []string) error {

	unknownCode, err := s.store.ListUnknownSamplingMethodCodes(ctx, tx, codes)
	if err != nil {
		return err
	}

	if len(unknownCode) > 0 {
		errs := make([]error, 0, len(unknownCode))
		for _, code := range unknownCode {
			errs = append(errs, models.ErrUnknownSamplingMethodCode(code))
		}
		return errors.Join(errs...)
	}

	return s.store.ReplaceMethodsAtSampling(ctx, tx, samplingID, codes)
}

// Initializes the sampling methods resolution state for a given import hash. This is typically called when starting a new import workflow.
func (s *SamplingService) InitMethodResolution(ctx context.Context, q db.Querier, importID uuid.UUID) (state []models.SamplingMethodResolution, err error) {
	return s.store.InitMethodResolution(ctx, q, importID)
}

func (s *SamplingService) GetMethodsResolution(ctx context.Context, q db.Querier, importID uuid.UUID) (state []models.SamplingMethodResolution, err error) {
	return s.store.GetMethodsResolution(ctx, q, importID)
}

func (s *SamplingService) ResolveMethod(ctx context.Context, q db.Querier, importID uuid.UUID, input models.SamplingMethodResolutionInput) (models.SamplingMethodResolution, error) {
	if err := input.Validate(); err != nil {
		return models.SamplingMethodResolution{}, err
	}
	return s.store.ResolveMethod(ctx, q, importID, input)
}

func (s *SamplingService) BootstrapSamplingMethods(ctx context.Context, q db.Querier, yamlBytes []byte) error {
	if methods, err := s.ListSamplingMethods(ctx, q); err != nil {
		return fmt.Errorf("failed to list sampling methods: %w", err)
	} else if len(methods) > 0 {
		logrus.Infof("Found %d existing sampling methods, skipping bootstrap", len(methods))
		return nil
	}

	var methods []models.SamplingMethodInput
	err := yaml.Unmarshal(yamlBytes, &methods)
	if err != nil {
		return fmt.Errorf("failed to unmarshal sampling methods: %w", err)
	}

	logrus.Infof("Bootstrapping %d sampling methods from YAML", len(methods))

	for _, method := range methods {
		_, err := s.CreateSamplingMethod(ctx, q, method)
		if err != nil {
			return fmt.Errorf("failed to create sampling method %s: %w", method.Code, err)
		}
	}
	return nil
}
