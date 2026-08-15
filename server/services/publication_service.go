package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/types"
)

type PublicationService struct {
}

func NewPublicationService() *PublicationService {
	return &PublicationService{}
}

func (s *PublicationService) ListPublications(ctx context.Context, q db.Querier) ([]models.Publication, error) {
	articles, err := q.Queries().ListPublications(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]models.Publication, len(articles))
	for i, a := range articles {
		result[i] = models.PublicationFromDB(a)
	}
	return result, nil
}

func (s *PublicationService) GetPublicationByID(ctx context.Context, q db.Querier, articleID uuid.UUID) (models.Publication, error) {
	article, err := q.Queries().GetPublicationByID(ctx, articleID)
	if err != nil {
		return models.Publication{}, err
	}
	return models.PublicationFromDB(article), nil
}

func (s *PublicationService) GetPublicationByDOI(ctx context.Context, q db.Querier, doi types.DOI) (models.Publication, error) {
	article, err := q.Queries().GetPublicationByDOI(ctx, &doi)
	if err != nil {
		return models.Publication{}, err
	}
	return models.PublicationFromDB(article), nil
}

func (s *PublicationService) CreatePublication(ctx context.Context, q db.Querier, params models.CreatePublicationParams) (models.Publication, error) {
	article, err := q.Queries().CreatePublication(ctx, params.ToDBParams())
	if err != nil {
		return models.Publication{}, err
	}
	return models.PublicationFromDB(article), nil
}

func (s *PublicationService) DeletePublicationByID(ctx context.Context, q db.Querier, articleID uuid.UUID) (models.Publication, error) {
	article, err := q.Queries().DeletePublicationByID(ctx, articleID)
	if err != nil {
		return models.Publication{}, err
	}
	return models.PublicationFromDB(article), nil
}
