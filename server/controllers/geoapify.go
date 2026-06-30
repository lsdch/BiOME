package controllers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services/geoapify"
)

type GeoapifyController struct {
	service *geoapify.GeoapifyService
}

func NewGeoapifyController(service *geoapify.GeoapifyService) *GeoapifyController {
	return &GeoapifyController{service: service}
}

func (c *GeoapifyController) GetStatus(
	ctx context.Context,
	input *struct{},
) (*BodyTransporter[geoapify.GeoapifyStatus], error) {
	status, err := c.service.GetStatus(ctx)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[geoapify.GeoapifyStatus]{Body: status}, nil
}

func (c *GeoapifyController) ListUsage(
	ctx context.Context,
	input *struct{},
) (*BodyTransporter[[]geoapify.GeoapifyUsage], error) {
	usage, err := c.service.ListUsage(ctx)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]geoapify.GeoapifyUsage]{Body: usage}, nil
}

func (c *GeoapifyController) ReverseGeocode(
	ctx context.Context,
	input *struct {
		Latitude  float32 `query:"lat" required:"true"`
		Longitude float32 `query:"lon" required:"true"`
	},
) (*BodyTransporter[geoapify.GeoapifyResult], error) {
	result, err := c.service.ReverseGeocode(ctx, input.Latitude, input.Longitude)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[geoapify.GeoapifyResult]{Body: *result}, nil
}

func (c *GeoapifyController) BatchReverseGeocode(
	ctx context.Context,
	input *BodyTransporter[[]geoapify.LatLonCoords],
) (*BodyTransporter[[]geoapify.GeoapifyResult], error) {
	results, err := c.service.BatchReverseGeocode(ctx, input.Body)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]geoapify.GeoapifyResult]{Body: results}, nil
}

func (c *GeoapifyController) RegisterRoutes(r *router.Router) {
	geoapifyAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/geoapify").
			WithTags([]string{"Services"})
	}

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
