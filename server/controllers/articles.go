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
	service *services.PublicationService
}

func NewArticlesController(db *db.DB, service *services.PublicationService) *ArticlesController {
	return &ArticlesController{
		db:      db,
		service: service,
	}
}

func (c *ArticlesController) ListPublications(ctx context.Context, _ *struct{}) (*BodyTransporter[[]models.Publication], error) {
	articles, err := c.service.ListPublications(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.Publication]{
		Body: articles,
	}, nil
}

func (c *ArticlesController) CreatePublication(ctx context.Context, input *BodyTransporter[models.CreatePublicationParams]) (*BodyTransporter[models.Publication], error) {
	article, err := c.service.CreatePublication(ctx, c.db, input.Body)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[models.Publication]{
		Body: article,
	}, nil
}

func (c *ArticlesController) DeletePublication(ctx context.Context, input *UUIDInput) (*struct{}, error) {
	_, err := c.service.DeletePublicationByID(ctx, c.db, input.ID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *ArticlesController) RegisterRoutes(r *router.Router) {
	publicationsAPI := r.RouteGroup("/publications").WithTags([]string{"Bibliography"})

	router.NewSpec(publicationsAPI, "ListPublications", huma.Operation{
		Method:  http.MethodGet,
		Path:    "/",
		Summary: "List bibliography",
	}, c.ListPublications).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(publicationsAPI, "CreatePublication", huma.Operation{
		Method:  http.MethodPost,
		Path:    "/",
		Summary: "Create a new publication",
	}, c.CreatePublication).WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).Register(r)

	router.NewSpec(publicationsAPI, "DeletePublication", huma.Operation{
		Method:  http.MethodDelete,
		Path:    "/{id}",
		Summary: "Delete a publication by ID",
	}, c.DeletePublication).WithAccessPolicy(auth.Role(biomedb.UserRoleMaintainer)).Register(r)
}
