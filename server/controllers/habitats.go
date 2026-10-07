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

type HabitatsController struct {
	db      *db.DB
	service *services.HabitatService
}

func NewHabitatsController(db *db.DB, service *services.HabitatService) *HabitatsController {
	return &HabitatsController{db: db, service: service}
}

func (c *HabitatsController) ListHabitats(ctx context.Context, _ *struct{}) (*BodyTransporter[[]models.Habitat], error) {
	habitats, err := c.service.ListHabitats(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.Habitat]{Body: habitats}, nil
}

func (c *HabitatsController) AddHabitat(ctx context.Context, input *BodyTransporter[models.HabitatInput]) (*BodyTransporter[models.Habitat], error) {
	habitat, err := c.service.CreateHabitat(ctx, c.db, input.Body)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[models.Habitat]{Body: habitat}, nil
}

func (c *HabitatsController) DeleteHabitat(ctx context.Context, input *UUIDInput) (*struct{}, error) {
	return &struct{}{}, c.db.WithTx(ctx, func(tx *db.Tx) error {
		return c.service.DeleteHabitat(ctx, tx, input.ID)
	})
}

type UpdateHabitatInput struct {
	UUIDInput
	Body models.HabitatUpdate
}

func (c *HabitatsController) UpdateHabitat(ctx context.Context, input *UpdateHabitatInput) (*struct{}, error) {
	return &struct{}{}, c.db.WithTx(ctx, func(tx *db.Tx) error {
		return c.service.UpdateHabitat(ctx, tx, input.ID, input.Body)
	})
}

func (c *HabitatsController) RegisterRoutes(r *router.Router) {
	habitatsAPI := r.RouteGroup("/habitats").WithTags([]string{"Habitats"})

	router.NewSpec(
		habitatsAPI,
		"ListHabitats",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/",
			Summary: "List all habitats",
		},
		c.ListHabitats,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		habitatsAPI,
		"CreateHabitat",
		huma.Operation{
			Method:  http.MethodPost,
			Path:    "/",
			Summary: "Create a new habitat",
		},
		c.AddHabitat,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleAdmin)).Register(r)

	router.NewSpec(
		habitatsAPI,
		"DeleteHabitat",
		huma.Operation{
			Method:  http.MethodDelete,
			Path:    "/{id}",
			Summary: "Delete a habitat by its ID",
		},
		c.DeleteHabitat,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleAdmin)).Register(r)

	router.NewSpec(
		habitatsAPI,
		"UpdateHabitat",
		huma.Operation{
			Method:  http.MethodPatch,
			Path:    "/{id}",
			Summary: "Update a habitat by its ID",
		},
		c.UpdateHabitat,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleAdmin)).Register(r)
}
