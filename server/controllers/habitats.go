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

func (c *HabitatsController) GetHabitatGroups(ctx context.Context, _ *struct{}) (*BodyTransporter[[]models.HabitatGroupWithElements], error) {
	groups, err := c.service.GetHabitatGroups(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.HabitatGroupWithElements]{Body: groups}, nil
}

func (c *HabitatsController) AddHabitatGroup(ctx context.Context, group *BodyTransporter[models.HabitatGroupInput]) (*struct{}, error) {
	return &struct{}{}, c.db.WithTx(ctx, func(tx *db.Tx) error {
		return c.service.CreateHabitatGroup(ctx, tx, group.Body)
	})
}

func (c *HabitatsController) DeleteHabitatGroup(ctx context.Context, input *UUIDInput) (*struct{}, error) {
	return &struct{}{}, c.db.WithTx(ctx, func(tx *db.Tx) error {
		return c.service.DeleteHabitatGroup(ctx, tx, input.ID)
	})
}

type UpdateHabitatGroupInput struct {
	UUIDInput
	Body models.HabitatGroupUpdate
}

func (c *HabitatsController) UpdateHabitatGroup(ctx context.Context, input *UpdateHabitatGroupInput) (*struct{}, error) {
	return &struct{}{}, c.db.WithTx(ctx, func(tx *db.Tx) error {
		return c.service.UpdateHabitatGroup(ctx, tx, input.ID, input.Body)
	})
}

func (c *HabitatsController) RegisterRoutes(r *router.Router) {
	habitatsAPI := r.RouteGroup("/habitats").WithTags([]string{"Habitats"})

	router.NewSpec(
		habitatsAPI,
		"GetHabitatGroups",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/",
			Summary: "Get all habitat groups with their elements",
		},
		c.GetHabitatGroups,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		habitatsAPI,
		"CreateHabitatGroup",
		huma.Operation{
			Method:  http.MethodPost,
			Path:    "/",
			Summary: "Create a new habitat group with its elements",
		},
		c.AddHabitatGroup,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleAdmin)).Register(r)

	router.NewSpec(
		habitatsAPI,
		"DeleteHabitatGroup",
		huma.Operation{
			Method:  http.MethodDelete,
			Path:    "/{id}",
			Summary: "Delete a habitat group by its ID",
		},
		c.DeleteHabitatGroup,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleAdmin)).Register(r)

	router.NewSpec(
		habitatsAPI,
		"UpdateHabitatGroup",
		huma.Operation{
			Method:  http.MethodPatch,
			Path:    "/{id}",
			Summary: "Update a habitat group by its ID",
		},
		c.UpdateHabitatGroup,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleAdmin)).Register(r)
}
