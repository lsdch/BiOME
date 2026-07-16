package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/types"
)

type ImportBatch struct {
	ID          types.ULID       `json:"id"`
	Label       string           `json:"label"`
	Description Optional[string] `json:"description,omitempty"`
	CreatedBy   uuid.UUID        `json:"created_by"`
	AssembledBy []string         `json:"assembled_by,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}

func ImportBatchFromDB(b biomedb.ImportBatch) ImportBatch {
	return ImportBatch{
		ID:          b.ID,
		Label:       b.Label,
		Description: NewOptionalFromPtr(b.Description),
		CreatedBy:   b.CreatedBy,
		AssembledBy: b.AssembledBy,
		CreatedAt:   b.CreatedAt,
	}
}
