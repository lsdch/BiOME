package controllers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
)

type LocationController struct {
	service *services.LocationService
}

func NewLocationController(service *services.LocationService) *LocationController {
	return &LocationController{
		service: service,
	}
}

func (c *LocationController) ListCountries(ctx context.Context, input *struct{}) (*BodyTransporter[[]models.Country], error) {
	countries, err := c.service.ListCountries(ctx)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.Country]{Body: countries}, nil
}

func (c *LocationController) ListCountriesSummary(ctx context.Context, input *struct{}) (*BodyTransporter[[]models.CountrySummary], error) {
	countries, err := c.service.ListCountriesSummary(ctx)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.CountrySummary]{Body: countries}, nil
}

func (c *LocationController) Register(r *router.Router) {
	group := func(r *router.Router) router.Group {
		return r.RouteGroup("/locations").
			WithTags([]string{"Location"})
	}
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
}
