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

type LocationController struct {
	db      *db.DB
	service *services.LocationService
}

func NewLocationController(db *db.DB, service *services.LocationService) *LocationController {
	return &LocationController{
		db:      db,
		service: service,
	}
}

func (c *LocationController) ListCountries(ctx context.Context, input *struct{}) (*BodyTransporter[[]models.Country], error) {
	countries, err := c.service.ListCountries(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.Country]{Body: countries}, nil
}

func (c *LocationController) ListCountriesSummary(ctx context.Context, input *struct{}) (*BodyTransporter[[]models.CountrySummary], error) {
	countries, err := c.service.ListCountriesSummary(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.CountrySummary]{Body: countries}, nil
}

func (c *LocationController) CoordinatesToCountry(ctx context.Context, input *models.Coordinates) (*BodyTransporter[*models.Country], error) {
	country, err := c.service.CoordinatesToCountry(ctx, c.db, input.Latitude, input.Longitude)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[*models.Country]{Body: country}, nil
}

func (c *LocationController) RegisterRoutes(r *router.Router) {
	group := r.RouteGroup("/locations").WithTags([]string{"Location"})

	router.NewSpec(group, "ListCountries", huma.Operation{
		Path:    "/countries",
		Method:  "GET",
		Summary: "List countries",
	}, c.ListCountries).
		WithAccessPolicy(auth.Public()).
		Register(r)

	router.NewSpec(group, "ListCountriesSummary", huma.Operation{
		Path:    "/countries/summary",
		Method:  "GET",
		Summary: "List countries summary",
	}, c.ListCountriesSummary).
		WithAccessPolicy(auth.Public()).
		Register(r)

	router.NewSpec(group, "CoordinatesToCountry", huma.Operation{
		Path:    "/countries/from-coordinates",
		Method:  "GET",
		Summary: "Get country from coordinates",
	}, c.CoordinatesToCountry).
		WithAccessPolicy(auth.Public()).
		Register(r)
}
