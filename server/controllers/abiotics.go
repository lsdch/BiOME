package controllers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
)

type AbioticsController struct {
	db      *db.DB
	service *services.AbioticService
}

func NewAbioticsController(db *db.DB, service *services.AbioticService) *AbioticsController {
	return &AbioticsController{
		db:      db,
		service: service,
	}
}

func (c AbioticsController) ListAbioticParameters(ctx context.Context, input *struct{}) (*BodyTransporter[[]models.AbioticParam], error) {
	params, err := c.service.ListAbioticParameters(ctx, c.db)
	if err != nil {
		return nil, err
	}

	return &BodyTransporter[[]models.AbioticParam]{Body: params}, nil
}

func (c *AbioticsController) RegisterRoutes(r *router.Router) {
	abioticsAPI := r.RouteGroup("/abiotics").WithTags([]string{"Abiotics"})
	router.NewSpec(abioticsAPI,
		"ListAbioticParameters",
		huma.Operation{
			Method:  "GET",
			Path:    "/parameters",
			Summary: "List abiotic parameters",
		},
		c.ListAbioticParameters,
	).WithAccessPolicy(auth.Public()).Register(r)
}
