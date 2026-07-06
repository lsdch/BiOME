package controllers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
)

type DatasetController struct {
	db      *db.DB
	service *services.DatasetsService
}

func NewDatasetController(db *db.DB, service *services.DatasetsService) *DatasetController {
	return &DatasetController{
		db:      db,
		service: service,
	}
}

func (c *DatasetController) LoadDatasetsForOccurrence(ctx context.Context, input *ULIDPath) (*BodyTransporter[[]models.Dataset], error) {
	datasets, err := c.service.LoadDatasetsForOccurrence(ctx, c.db, input.ULID)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.Dataset]{Body: datasets}, nil
}

func (c *DatasetController) ListDatasets(ctx context.Context, input *struct{}) (*BodyTransporter[[]models.Dataset], error) {
	datasets, err := c.service.ListDatasets(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.Dataset]{Body: datasets}, nil
}

func (c *DatasetController) GetDatasetByID(ctx context.Context, input *ULIDPath) (*BodyTransporter[*models.Dataset], error) {
	dataset, err := c.service.GetDatasetByID(ctx, c.db, input.ULID)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[*models.Dataset]{Body: dataset}, nil
}

func (c *DatasetController) LoadOccurrencesForDataset(ctx context.Context, input *ULIDPath) (*BodyTransporter[[]models.Occurrence], error) {
	occurrences, err := c.service.LoadOccurrencesForDataset(ctx, c.db, input.ULID)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.Occurrence]{Body: occurrences}, nil
}

func (c *DatasetController) RegisterRoutes(r *router.Router) {
	group := r.RouteGroup("/datasets").WithTags([]string{"Datasets"})
	router.NewSpec(group, "ListDatasets",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/",
			Summary: "List datasets",
		},
		c.ListDatasets,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(r.RouteGroup("/occurrences").WithTags([]string{"Occurrences"}), "LoadDatasetsForOccurrence",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/{ulid}/datasets",
			Summary: "Load datasets for occurrence",
		},
		c.LoadDatasetsForOccurrence,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(group, "GetDatasetByID",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/{ulid}",
			Summary: "Get dataset by ID",
		},
		c.GetDatasetByID,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(group, "LoadOccurrencesForDataset",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/{ulid}/occurrences",
			Summary: "Load occurrences for dataset",
		},
		c.LoadOccurrencesForDataset,
	).WithAccessPolicy(auth.Public()).Register(r)
}
