package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
)

type ImportBatchStatus = biomedb.ImportBatchStatus

type ImportBatch struct {
	ID               uuid.UUID           `json:"id"`
	Label            string              `json:"label"`
	Description      Optional[string]    `json:"description,omitzero"`
	Status           ImportBatchStatus   `json:"status"`
	CreatedBy        uuid.UUID           `json:"created_by"`
	AssembledBy      []string            `json:"assembled_by,omitempty"`
	CreatedAt        time.Time           `json:"created_at"`
	CompletedAt      Optional[time.Time] `json:"completed_at,omitzero"`
	CompletedBy      Optional[uuid.UUID] `json:"completed_by,omitzero"`
	TaxonomicScope   int32               `json:"taxonomic_scope"`
	ImportedFileName string              `json:"imported_file_name"`
	ImportedFileSize int64               `json:"imported_file_size"`
	ImportedFileHash string              `json:"imported_file_hash"`
}

func ImportBatchFromDB(b biomedb.ImportBatch) ImportBatch {
	return ImportBatch{
		ID:               b.ID,
		Label:            b.Label,
		Description:      NewOptionalFromPtr(b.Description),
		Status:           b.Status,
		CreatedBy:        b.CreatedBy,
		AssembledBy:      b.AssembledBy,
		CreatedAt:        b.CreatedAt,
		CompletedAt:      NewOptionalFromTimestamp(b.CompletedAt),
		CompletedBy:      NewOptionalFromUUID(b.CompletedBy),
		TaxonomicScope:   b.TaxonomicScope,
		ImportedFileName: b.ImportedFileName,
		ImportedFileSize: b.ImportedFileSize,
		ImportedFileHash: b.ImportedFileHash,
	}
}

func (b ImportBatch) WithContent(occurrenceCount, samplingCount int64, createdBy, completedBy User) ImportBatchWithContent {
	return ImportBatchWithContent{
		ImportBatch:     b,
		OccurrenceCount: occurrenceCount,
		SamplingCount:   samplingCount,
		CreatedByUser:   createdBy,
		CompletedByUser: completedBy,
	}
}

type ImportBatchWithContent struct {
	ImportBatch
	OccurrenceCount int64 `json:"occurrence_count"`
	SamplingCount   int64 `json:"sampling_count"`
	CreatedByUser   User  `json:"created_by_user"`
	CompletedByUser User  `json:"completed_by_user"`
}

type ImportBatchInput struct {
	Label          string           `json:"label" form:"label" required:"true"`
	Description    Optional[string] `json:"description,omitzero" form:"description"`
	AssembledBy    []string         `json:"assembled_by,omitempty" form:"assembled_by"`
	TaxonomicScope int32            `json:"taxonomic_scope" form:"taxonomic_scope" required:"true"`
}

func (i ImportBatchInput) ToParams(userID uuid.UUID, file FileMetadata) biomedb.InitImportBatchParams {
	return biomedb.InitImportBatchParams{
		Label:            i.Label,
		Description:      i.Description.ToPtr(),
		AssembledBy:      i.AssembledBy,
		TaxonomicScope:   i.TaxonomicScope,
		CreatedBy:        userID,
		ImportedFileName: file.Name,
		ImportedFileSize: file.Size,
		ImportedFileHash: file.Hash,
	}
}

type FileMetadata struct {
	Name string `json:"name" form:"name" required:"true"`
	Size int64  `json:"size" form:"size" required:"true"`
	Hash string `json:"hash" form:"hash" required:"true"`
}

type ImportBatchWithFileInput struct {
	ImportBatchInput
	File FileMetadata `json:"file" form:"file" required:"true"`
}

func (i ImportBatchWithFileInput) ToParams(userID uuid.UUID) biomedb.InitImportBatchParams {
	return i.ImportBatchInput.ToParams(userID, i.File)
}
