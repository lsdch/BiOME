package controllers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
	"github.com/lsdch/biome/types"
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

func (c *DatasetController) ListDatasets(ctx context.Context, input *struct{}) (*BodyTransporter[[]models.DatasetWithSummary], error) {
	datasets, err := c.service.ListDatasets(ctx, c.db)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.DatasetWithSummary]{Body: datasets}, nil
}

func (c *DatasetController) GetDatasetByID(ctx context.Context, input *ULIDPath) (*BodyTransporter[*models.DatasetWithMaintainers], error) {
	dataset, err := c.service.GetDatasetByID(ctx, c.db, input.ULID)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[*models.DatasetWithMaintainers]{Body: dataset}, nil
}

func (c *DatasetController) LoadOccurrencesForDataset(ctx context.Context, input *ULIDPath) (*BodyTransporter[[]models.Occurrence], error) {
	occurrences, err := c.service.LoadOccurrencesForDataset(ctx, c.db, input.ULID)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.Occurrence]{Body: occurrences}, nil
}

func (c *DatasetController) ListImportBatchesInDataset(ctx context.Context, input *ULIDPath) (*BodyTransporter[[]models.ImportBatchListItem], error) {
	importBatches, err := c.service.ListImportBatchesInDataset(ctx, c.db, input.ULID)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[[]models.ImportBatchListItem]{Body: importBatches}, nil
}

func (c *DatasetController) CreateDataset(ctx context.Context, input *BodyTransporter[models.DatasetInput]) (*BodyTransporter[*models.DatasetWithMaintainers], error) {
	var dataset *models.DatasetWithMaintainers
	session, _ := auth.SessionFromContext(ctx)
	err := c.db.WithTx(ctx, func(tx *db.Tx) error {
		d, err := c.service.CreateDataset(ctx, tx, session.UserID, input.Body)
		if err != nil {
			return err
		}
		dataset = d
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[*models.DatasetWithMaintainers]{Body: dataset}, nil
}

func (c *DatasetController) AddCuratorToDataset(ctx context.Context, input *struct {
	ULIDPath
	UserID uuid.UUID `path:"user_id" format:"uuid"`
}) (*BodyTransporter[*models.DatasetWithMaintainers], error) {

	err := c.service.AddCuratorToDataset(ctx, c.db, input.ULID, input.UserID)
	if err != nil {
		return nil, err
	}
	d, err := c.service.GetDatasetByID(ctx, c.db, input.ULID)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[*models.DatasetWithMaintainers]{Body: d}, nil
}

func (c *DatasetController) RemoveCuratorFromDataset(ctx context.Context, input *struct {
	ULIDPath
	UserID uuid.UUID `path:"user_id" format:"uuid"`
}) (*BodyTransporter[*models.DatasetWithMaintainers], error) {
	err := c.service.RemoveCuratorFromDataset(ctx, c.db, input.ULID, input.UserID)
	if err != nil {
		return nil, err
	}
	d, err := c.service.GetDatasetByID(ctx, c.db, input.ULID)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[*models.DatasetWithMaintainers]{Body: d}, nil
}

func (c *DatasetController) AddOccurrence(ctx context.Context, input *struct {
	ULIDPath
	OccurrenceID types.ULID `path:"occurrence_id"`
}) (*struct{}, error) {
	err := c.service.AddOccurrenceToDataset(ctx, c.db, input.ULID, input.OccurrenceID)
	return nil, err
}

func (c *DatasetController) RemoveOccurrence(ctx context.Context, input *struct {
	ULIDPath
	OccurrenceID types.ULID `path:"occurrence_id"`
}) (*struct{}, error) {
	err := c.service.RemoveOccurrenceFromDataset(ctx, c.db, input.ULID, input.OccurrenceID)
	return nil, err
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

	router.NewSpec(group, "CreateDataset",
		huma.Operation{
			Method:  http.MethodPost,
			Path:    "/",
			Summary: "Create a new dataset",
		},
		c.CreateDataset,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).Register(r)

	router.NewSpec(group, "AddCuratorToDataset",
		huma.Operation{
			Method:  http.MethodPost,
			Path:    "/{ulid}/curators/{user_id}",
			Summary: "Add a curator to a dataset",
		},
		c.AddCuratorToDataset,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).Register(r)

	router.NewSpec(group, "RemoveCuratorFromDataset",
		huma.Operation{
			Method:  http.MethodDelete,
			Path:    "/{ulid}/curators/{user_id}",
			Summary: "Remove a curator from a dataset",
		},
		c.RemoveCuratorFromDataset,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).Register(r)

	router.NewSpec(group, "ListImportBatchesInDataset",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/{ulid}/batches",
			Summary: "List import batches in a dataset",
		},
		c.ListImportBatchesInDataset,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(group, "AddOccurrenceToDataset",
		huma.Operation{
			Method:  http.MethodPost,
			Path:    "/{ulid}/occurrences/{occurrence_id}",
			Summary: "Add an occurrence to a dataset",
		},
		c.AddOccurrence,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).Register(r)

	router.NewSpec(group, "RemoveOccurrenceFromDataset",
		huma.Operation{
			Method:  http.MethodDelete,
			Path:    "/{ulid}/occurrences/{occurrence_id}",
			Summary: "Remove an occurrence from a dataset",
		},
		c.RemoveOccurrence,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleContributor)).Register(r)
}
