package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
)

type ArticleService struct {
	db *db.DB
}

func NewArticleService(db *db.DB) *ArticleService {
	return &ArticleService{
		db: db,
	}
}

func (s *ArticleService) ListArticles(ctx context.Context) ([]models.Article, error) {
	articles, err := s.db.Queries().ListArticles(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]models.Article, len(articles))
	for i, a := range articles {
		result[i] = models.ArticleFromDB(a)
	}
	return result, nil
}

func (s *ArticleService) GetArticleByID(ctx context.Context, articleID uuid.UUID) (models.Article, error) {
	article, err := s.db.Queries().GetArticleByID(ctx, articleID)
	if err != nil {
		return models.Article{}, err
	}
	return models.ArticleFromDB(article), nil
}

func (s *ArticleService) GetArticleByDOI(ctx context.Context, doi models.DOI) (models.Article, error) {
	doiStr := doi.String()
	article, err := s.db.Queries().GetArticleByDOI(ctx, &doiStr)
	if err != nil {
		return models.Article{}, err
	}
	return models.ArticleFromDB(article), nil
}

func (s *ArticleService) CreateArticle(ctx context.Context, params models.CreateArticleParams) (models.Article, error) {
	article, err := s.db.Queries().CreateArticle(ctx, params.ToDBParams())
	if err != nil {
		return models.Article{}, err
	}
	return models.ArticleFromDB(article), nil
}

func (s *ArticleService) DeleteArticleByID(ctx context.Context, articleID uuid.UUID) (models.Article, error) {
	article, err := s.db.Queries().DeleteArticleByID(ctx, articleID)
	if err != nil {
		return models.Article{}, err
	}
	return models.ArticleFromDB(article), nil
}
