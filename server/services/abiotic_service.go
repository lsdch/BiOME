package services

import (
	"context"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
)

type AbioticService struct {
}

func NewAbioticService(db *db.DB) *AbioticService {
	return &AbioticService{}
}

func (s *AbioticService) ListAbioticParameters(ctx context.Context, q db.Querier) ([]models.AbioticParam, error) {
	paramsDB, err := q.Queries().ListAbioticParams(ctx)
	if err != nil {
		return nil, err
	}

	params := make([]models.AbioticParam, len(paramsDB))
	for i, p := range paramsDB {
		params[i] = models.AbioticParamFromDB(p)
	}
	return params, nil
}

func (s *AbioticService) CreateAbioticParam(ctx context.Context, q db.Querier, input models.AbioticParamInput) (models.AbioticParam, error) {
	paramDB, err := q.Queries().CreateAbioticParam(ctx, input.ToDB())
	if err != nil {
		return models.AbioticParam{}, err
	}

	return models.AbioticParamFromDB(paramDB), nil
}
