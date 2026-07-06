package controllers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
)

type ArticlesController struct {
	db      *db.DB
	service *services.ArticleService
}

func NewArticlesController(db *db.DB, service *services.ArticleService) *ArticlesController {
	return &ArticlesController{
		db:      db,
		service: service,
	}
}

func (c *ArticlesController) ListArticles(ctx context.Context, _ *struct{}) (*BodyTransporter[[]models.Article], error) {
	articles, err := c.service.ListArticles(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.Article]{
		Body: articles,
	}, nil
}

func (c *ArticlesController) CreateArticle(ctx context.Context, input *BodyTransporter[models.CreateArticleParams]) (*BodyTransporter[models.Article], error) {
	article, err := c.service.CreateArticle(ctx, c.db, input.Body)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[models.Article]{
		Body: article,
	}, nil
}

func (c *ArticlesController) DeleteArticle(ctx context.Context, input *UUIDInput) (*struct{}, error) {
	_, err := c.service.DeleteArticleByID(ctx, c.db, input.ID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *ArticlesController) RegisterRoutes(r *router.Router) {
	articlesAPI := r.RouteGroup("/articles").WithTags([]string{"Bibliography"})

	router.NewSpec(articlesAPI, "ListArticles", huma.Operation{
		Method:  http.MethodGet,
		Path:    "/",
		Summary: "List bibliography",
	}, c.ListArticles).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(articlesAPI, "CreateArticle", huma.Operation{
		Method:  http.MethodPost,
		Path:    "/",
		Summary: "Create a new article",
	}, c.CreateArticle).WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).Register(r)

	router.NewSpec(articlesAPI, "DeleteArticle", huma.Operation{
		Method:  http.MethodDelete,
		Path:    "/{id}",
		Summary: "Delete an article by ID",
	}, c.DeleteArticle).WithAccessPolicy(auth.Role(biomedb.UserRoleMaintainer)).Register(r)
}
