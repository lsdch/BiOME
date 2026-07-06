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

type AccountsController struct {
	db      *db.DB
	service *services.AccountService
}

func NewAccountsController(db *db.DB, service *services.AccountService) *AccountsController {
	return &AccountsController{db: db, service: service}
}

func (c *AccountsController) ListUsers(ctx context.Context, _ *struct{}) (*BodyTransporter[[]models.User], error) {
	users, err := c.service.ListUsers(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.User]{Body: users}, nil
}

func (c *AccountsController) RegisterRoutes(r *router.Router) {
	accountsAPI := r.RouteGroup("/accounts").WithTags([]string{"Accounts"})

	router.NewSpec(
		accountsAPI,
		"ListUsers",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/",
			Summary: "List all users",
		},
		c.ListUsers,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleAdmin)).Register(r)
}
