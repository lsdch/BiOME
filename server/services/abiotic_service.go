package services

import (
	"context"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
)

type AbioticService struct {
	db *db.DB
}

func NewAbioticService(db *db.DB) *AbioticService {
	return &AbioticService{db: db}
}

func (s *AbioticService) ListAbioticParameters(ctx context.Context) ([]models.AbioticParam, error) {
	paramsDB, err := s.db.Queries().ListAbioticParams(ctx)
	if err != nil {
		return nil, err
	}

	params := make([]models.AbioticParam, len(paramsDB))
	for i, p := range paramsDB {
		params[i] = models.AbioticParamFromDB(p)
	}
	return params, nil
}

func (s *AbioticService) CreateAbioticParam(ctx context.Context, input models.AbioticParamInput) (models.AbioticParam, error) {
	paramDB, err := s.db.Queries().CreateAbioticParam(ctx, input.ToDB())
	if err != nil {
		return models.AbioticParam{}, err
	}

	return models.AbioticParamFromDB(paramDB), nil
}
