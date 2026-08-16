package controllers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
)

func (c *SamplingController) ListSamplingFixatives(
	ctx context.Context,
	input *struct{},
) (*BodyTransporter[[]models.Fixative], error) {
	fixatives, err := c.service.ListSamplingFixatives(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.Fixative]{
		Body: fixatives,
	}, nil
}

func (c *SamplingController) CreateSamplingFixative(
	ctx context.Context,
	input *BodyTransporter[models.FixativeInput],
) (*BodyTransporter[models.Fixative], error) {
	fixative, err := c.service.CreateSamplingFixative(ctx, c.db, input.Body)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[models.Fixative]{
		Body: fixative,
	}, nil
}

func (c *SamplingController) UpdateSamplingFixative(
	ctx context.Context,
	input *struct {
		Body models.FixativeUpdateParams
		CodePath
	},
) (*BodyTransporter[models.Fixative], error) {
	fixative, err := c.service.UpdateSamplingFixative(ctx, c.db, input.Code, input.Body)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[models.Fixative]{
		Body: fixative,
	}, nil
}

func (c *SamplingController) DeleteSamplingFixative(
	ctx context.Context,
	input *CodePath,
) (*struct{}, error) {
	err := c.service.DeleteSamplingFixative(ctx, c.db, input.Code)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *SamplingController) RegisterSamplingFixativesRoutes(r *router.Router) {
	group := r.RouteGroup("/fixatives").WithTags([]string{"Sampling"})

	router.NewSpec(group, "ListFixatives", huma.Operation{
		Method:  "GET",
		Path:    "/",
		Summary: "List sampling fixatives",
	}, c.ListSamplingFixatives).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(group, "CreateFixative", huma.Operation{
		Method:  "POST",
		Path:    "/",
		Summary: "Create a new sampling fixative",
	}, c.CreateSamplingFixative).WithAccessPolicy(auth.Role(biomedb.UserRoleMaintainer)).Register(r)

	router.NewSpec(group, "UpdateFixative", huma.Operation{
		Method:  "PATCH",
		Path:    "/{code}",
		Summary: "Update a sampling fixative by code",
	}, c.UpdateSamplingFixative).WithAccessPolicy(auth.Role(biomedb.UserRoleMaintainer)).Register(r)

	router.NewSpec(group, "DeleteFixative", huma.Operation{
		Method:  "DELETE",
		Path:    "/{code}",
		Summary: "Delete a sampling fixative by code",
	}, c.DeleteSamplingFixative).WithAccessPolicy(auth.Role(biomedb.UserRoleMaintainer)).Register(r)

}
