package controllers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
)

func (c *SamplingController) ListSamplingMethods(
	ctx context.Context,
	input *struct{},
) (*BodyTransporter[[]models.SamplingMethod], error) {
	methods, err := c.service.ListSamplingMethods(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.SamplingMethod]{
		Body: methods,
	}, nil
}

func (c *SamplingController) CreateSamplingMethod(ctx context.Context, input *BodyTransporter[models.SamplingMethodInput]) (*BodyTransporter[models.SamplingMethod], error) {
	method, err := c.service.CreateSamplingMethod(ctx, c.db, input.Body)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[models.SamplingMethod]{
		Body: method,
	}, nil
}

type SamplingMethodUpdate struct {
	Body models.SamplingMethodUpdateParams
	CodePath
}

func (c *SamplingController) UpdateSamplingMethod(ctx context.Context, input *SamplingMethodUpdate) (*BodyTransporter[models.SamplingMethod], error) {
	method, err := c.service.UpdateSamplingMethod(ctx, c.db, input.Code, input.Body)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[models.SamplingMethod]{
		Body: method,
	}, nil
}

func (c *SamplingController) DeleteSamplingMethod(ctx context.Context, input *CodePath) (*struct{}, error) {
	err := c.service.DeleteSamplingMethod(ctx, c.db, input.Code)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *SamplingController) RegisterSamplingMethodsRoutes(r *router.Router) {
	methodsGroup := r.RouteGroup("/sampling-methods").WithTags([]string{"Sampling"})

	router.NewSpec(methodsGroup, "ListSamplingMethods", huma.Operation{
		Method:  "GET",
		Path:    "/",
		Summary: "List sampling methods",
	}, c.ListSamplingMethods).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(methodsGroup, "CreateSamplingMethod", huma.Operation{
		Method:  "POST",
		Path:    "/",
		Summary: "Create a new sampling method",
	}, c.CreateSamplingMethod).WithAccessPolicy(auth.Role(biomedb.UserRoleMaintainer)).Register(r)

	router.NewSpec(methodsGroup, "UpdateSamplingMethod", huma.Operation{
		Method:  "PATCH",
		Path:    "/{code}",
		Summary: "Update a sampling method by code",
	}, c.UpdateSamplingMethod).WithAccessPolicy(auth.Role(biomedb.UserRoleMaintainer)).Register(r)

	router.NewSpec(methodsGroup, "DeleteSamplingMethod", huma.Operation{
		Method:  "DELETE",
		Path:    "/{code}",
		Summary: "Delete a sampling method by code",
	}, c.DeleteSamplingMethod).WithAccessPolicy(auth.Role(biomedb.UserRoleMaintainer)).Register(r)
}
