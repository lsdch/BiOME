package controllers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services/geoapify"
)

type GeoapifyController struct {
	db      *db.DB
	service *geoapify.GeoapifyService
}

func NewGeoapifyController(db *db.DB, service *geoapify.GeoapifyService) *GeoapifyController {
	return &GeoapifyController{db: db, service: service}
}

func (c *GeoapifyController) GetStatus(
	ctx context.Context,
	input *struct{},
) (*BodyTransporter[models.GeoapifyStatus], error) {
	status, err := c.service.GetStatus(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[models.GeoapifyStatus]{Body: status}, nil
}

func (c *GeoapifyController) ListUsage(
	ctx context.Context,
	input *struct{},
) (*BodyTransporter[[]models.GeoapifyUsage], error) {
	usage, err := c.service.ListUsage(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.GeoapifyUsage]{Body: usage}, nil
}

func (c *GeoapifyController) ReverseGeocode(
	ctx context.Context,
	input *models.Coordinates,
) (*BodyTransporter[models.GeoapifyResult], error) {
	result, err := c.service.ReverseGeocode(ctx, c.db, input.Latitude, input.Longitude)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[models.GeoapifyResult]{Body: *result}, nil
}

func (c *GeoapifyController) BatchReverseGeocode(
	ctx context.Context,
	input *BodyTransporter[[]models.GeoapifyCoords],
) (*BodyTransporter[[]models.GeoapifyResult], error) {
	results, err := c.service.BatchReverseGeocode(ctx, c.db, input.Body)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.GeoapifyResult]{Body: results}, nil
}

func (c *GeoapifyController) RegisterRoutes(r *router.Router) {
	geoapifyAPI := r.RouteGroup("/geoapify").WithTags([]string{"Services"})

	router.NewSpec(
		geoapifyAPI,
		"ListGeoapifyUsage",
		huma.Operation{
			Path:    "/usage/history",
			Method:  http.MethodGet,
			Summary: "List Geoapify usage history",
		},
		c.ListUsage,
	).
		WithAccessPolicy(auth.Authenticated()).
		Register(r)

	router.NewSpec(
		geoapifyAPI,
		"GetGeoapifyStatus",
		huma.Operation{
			Path:    "/status",
			Method:  http.MethodGet,
			Summary: "Get Geoapify API status",
		},
		c.GetStatus,
	).
		WithAccessPolicy(auth.Authenticated()).
		Register(r)

	router.NewSpec(
		geoapifyAPI,
		"ReverseGeocode",
		huma.Operation{
			Path:    "/reverse-geocode",
			Method:  http.MethodGet,
			Summary: "Reverse geocode coordinates using Geoapify API",
		},
		c.ReverseGeocode,
	).
		WithAccessPolicy(auth.Authenticated()).
		Register(r)

	router.NewSpec(
		geoapifyAPI,
		"BatchReverseGeocode",
		huma.Operation{
			Path:    "/reverse-geocode/batch",
			Method:  http.MethodPost,
			Summary: "Batch reverse geocode coordinates using Geoapify API",
		},
		c.BatchReverseGeocode,
	).
		WithAccessPolicy(auth.Authenticated()).
		Register(r)

}
