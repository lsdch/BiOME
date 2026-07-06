package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
)

type ArticleService struct {
}

func NewArticleService() *ArticleService {
	return &ArticleService{}
}

func (s *ArticleService) ListArticles(ctx context.Context, q db.Querier) ([]models.Article, error) {
	articles, err := q.Queries().ListArticles(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]models.Article, len(articles))
	for i, a := range articles {
		result[i] = models.ArticleFromDB(a)
	}
	return result, nil
}

func (s *ArticleService) GetArticleByID(ctx context.Context, q db.Querier, articleID uuid.UUID) (models.Article, error) {
	article, err := q.Queries().GetArticleByID(ctx, articleID)
	if err != nil {
		return models.Article{}, err
	}
	return models.ArticleFromDB(article), nil
}

func (s *ArticleService) GetArticleByDOI(ctx context.Context, q db.Querier, doi models.DOI) (models.Article, error) {
	doiStr := doi.String()
	article, err := q.Queries().GetArticleByDOI(ctx, &doiStr)
	if err != nil {
		return models.Article{}, err
	}
	return models.ArticleFromDB(article), nil
}

func (s *ArticleService) CreateArticle(ctx context.Context, q db.Querier, params models.CreateArticleParams) (models.Article, error) {
	article, err := q.Queries().CreateArticle(ctx, params.ToDBParams())
	if err != nil {
		return models.Article{}, err
	}
	return models.ArticleFromDB(article), nil
}

func (s *ArticleService) DeleteArticleByID(ctx context.Context, q db.Querier, articleID uuid.UUID) (models.Article, error) {
	article, err := q.Queries().DeleteArticleByID(ctx, articleID)
	if err != nil {
		return models.Article{}, err
	}
	return models.ArticleFromDB(article), nil
}
