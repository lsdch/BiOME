package controllers

import (
	"context"
	"io"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
	"github.com/lsdch/biome/services/storage"
	"github.com/sirupsen/logrus"
)

type ImportBatchController struct {
	DB          *db.DB
	service     *services.ImportBatchService
	fileStorage storage.RawFileStorage
}

func NewImportBatchController(db *db.DB, service *services.ImportBatchService, fileStorage storage.RawFileStorage) *ImportBatchController {
	return &ImportBatchController{
		DB:          db,
		service:     service,
		fileStorage: fileStorage,
	}
}

func (c *ImportBatchController) GetImportBatch(ctx context.Context, input *struct {
	UUIDInput
}) (*BodyTransporter[models.ImportBatch], error) {
	batch, err := c.service.GetImportBatch(ctx, c.DB, input.ID)
	if err != nil {
		return &BodyTransporter[models.ImportBatch]{}, err
	}

	return &BodyTransporter[models.ImportBatch]{Body: batch}, nil
}

func (c *ImportBatchController) GetImportBatchWithContent(ctx context.Context, input *struct {
	UUIDInput
}) (*BodyTransporter[models.ImportBatchWithContent], error) {
	batch, err := c.service.GetImportBatchWithContent(ctx, c.DB, input.ID)
	if err != nil {
		return &BodyTransporter[models.ImportBatchWithContent]{}, err
	}

	return &BodyTransporter[models.ImportBatchWithContent]{Body: batch}, nil
}

func (c *ImportBatchController) ListImportBatches(ctx context.Context, input *struct{}) (*BodyTransporter[[]models.ImportBatch], error) {
	batches, err := c.service.ListImportBatches(ctx, c.DB)
	if err != nil {
		return &BodyTransporter[[]models.ImportBatch]{}, err
	}

	return &BodyTransporter[[]models.ImportBatch]{Body: batches}, nil
}
func (c *ImportBatchController) ListImportBatchesWithContent(ctx context.Context, input *struct{}) (*BodyTransporter[[]models.ImportBatchWithContent], error) {
	batches, err := c.service.ListImportBatchesWithContent(ctx, c.DB)
	if err != nil {
		return &BodyTransporter[[]models.ImportBatchWithContent]{}, err
	}

	return &BodyTransporter[[]models.ImportBatchWithContent]{Body: batches}, nil
}

func (c *ImportBatchController) DeleteImportBatch(ctx context.Context, input *struct {
	UUIDInput
}) (*struct{}, error) {
	err := c.DB.WithTx(ctx, func(tx *db.Tx) error {
		return c.service.DeleteImportBatch(ctx, tx, input.ID)
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *ImportBatchController) DownloadRawFile(ctx context.Context, input *struct {
	UUIDInput
}) (*huma.StreamResponse, error) {
	batch, err := c.service.GetImportBatch(ctx, c.DB, input.ID)
	if err != nil {
		return nil, err
	}
	reader, err := c.fileStorage.Open(ctx, batch.FileKey())
	if err != nil {
		return nil, err
	}
	return &huma.StreamResponse{
		Body: func(ctx huma.Context) {
			writer := ctx.BodyWriter()
			ctx.SetHeader("Content-Type", batch.ImportedFileContentType)
			ctx.SetHeader("Content-Disposition", "attachment; filename=\""+batch.ImportedFileName+"\"")
			_, err := io.Copy(writer, reader)
			if err != nil {
				logrus.Errorf("Error streaming raw file: %v", err)
			}
			reader.Close()
		},
	}, nil
}

func (c *ImportBatchController) RegisterRoutes(r *router.Router) {

	batchesAPI := r.RouteGroup("/import-batches").WithTags([]string{"Imports"})

	router.NewSpec(batchesAPI,
		"GetImportBatch",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/{id}",
			Summary: "Get a specific import batch",
		},
		c.GetImportBatch,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(batchesAPI,
		"GetImportBatchWithContent",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/{id}/with-content",
			Summary: "Get a specific import batch with content summary",
		},
		c.GetImportBatchWithContent,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(batchesAPI,
		"ListImportBatches",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/",
			Summary: "List import batches",
		},
		c.ListImportBatches,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(batchesAPI,
		"ListImportBatchesWithContent",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/with-content",
			Summary: "List import batches with content summary",
		},
		c.ListImportBatchesWithContent,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(batchesAPI,
		"DeleteImportBatch",
		huma.Operation{
			Method:  http.MethodDelete,
			Path:    "/{id}",
			Summary: "Delete an import batch",
		},
		c.DeleteImportBatch,
	).WithAccessPolicy(auth.Role(biomedb.UserRoleMaintainer)).Register(r)

	router.NewSpec(batchesAPI,
		"DownloadRawFile",
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/{id}/raw",
			Summary: "Download the raw file for an import batch",
		},
		c.DownloadRawFile,
	).WithAccessPolicy(auth.Public()).Register(r)
}
