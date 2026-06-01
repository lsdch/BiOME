package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/db/pgutils"
	"github.com/lsdch/biome/models"
)

type SamplingService struct {
	q *biomedb.Queries
}

func NewSamplingService(q *biomedb.Queries) *SamplingService {
	return &SamplingService{q: q}
}

func (s *SamplingService) ListSamplingMethods(ctx context.Context) ([]biomedb.SamplingMethod, error) {
	return s.q.ListSamplingMethods(ctx)
}

func (s *SamplingService) CreateSamplingMethod(ctx context.Context, input biomedb.CreateSamplingMethodParams) (biomedb.SamplingMethod, error) {
	return s.q.CreateSamplingMethod(ctx, input)
}

func (s *SamplingService) UpdateSamplingMethod(ctx context.Context, code string, input SamplingMethodUpdateParams) (biomedb.SamplingMethod, error) {
	return s.q.UpdateSamplingMethod(ctx, input.ToParams(code))
}

func (s *SamplingService) DeleteSamplingMethod(ctx context.Context, code string) error {
	return s.q.DeleteSamplingMethod(ctx, code)
}

type SamplingMethodUpdateParams struct {
	Name        models.OptionalInput[string] `json:"name"`
	Code        models.OptionalInput[string] `json:"code"`
	Description models.OptionalNull[string]  `json:"description"`
}

func (s *SamplingMethodUpdateParams) ToParams(oldCode string) biomedb.UpdateSamplingMethodParams {
	return biomedb.UpdateSamplingMethodParams{
		Code:           pgutils.TextOpt(s.Code),
		Name:           pgutils.TextOpt(s.Name),
		SetDescription: s.Description.IsSet,
		Description:    pgutils.TextOpt(s.Description),
		OldCode:        oldCode,
	}
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

type SamplingMethodResolutionInput struct {
	InputText        string                          `json:"input_text"`
	ResolvedMethodId models.OptionalInput[uuid.UUID] `json:"resolved_method_id,omitempty"`
	Status           biomedb.MethodResolutionStatus  `json:"status"`
}

func (i SamplingMethodResolutionInput) Validate() error {
	if i.Status == biomedb.MethodResolutionStatusSelected && !i.ResolvedMethodId.IsSet {
		return models.WrapErrorPath(fmt.Errorf("resolution status is 'selected' but no method was provided"), "resolved_method_id")
	}
	return nil
}

func (i SamplingMethodResolutionInput) ToParams(importHash string) biomedb.ResolveMethodParams {
	return biomedb.ResolveMethodParams{
		ImportHash:       importHash,
		InputText:        i.InputText,
		ResolvedMethodID: pgutils.UUIDOpt(i.ResolvedMethodId),
		Status:           i.Status,
	}
}

func (s *SamplingService) ResolveMethod(ctx context.Context, importHash string, input SamplingMethodResolutionInput) (biomedb.SamplingMethodsResolution, error) {
	if err := input.Validate(); err != nil {
		return biomedb.SamplingMethodsResolution{}, err
	}
	return s.q.ResolveMethod(ctx, input.ToParams(importHash))
}
