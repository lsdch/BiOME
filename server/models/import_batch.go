package models

import (
	"time"

	"github.com/lsdch/biome/db/biomedb"
	"github.com/oklog/ulid/v2"
)

type ImportBatch struct {
	ID          ulid.ULID        `json:"id"`
	Label       string           `json:"label"`
	Description Optional[string] `json:"description,omitempty"`
	SubmittedBy Optional[string] `json:"submitted_by,omitempty"`
	AssembledBy []string         `json:"assembled_by,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}

func ImportBatchFromDB(b biomedb.ImportBatch) ImportBatch {
	return ImportBatch{
		ID:          b.ID,
		Label:       b.Label,
		Description: NewOptionalFromPtr(b.Description),
		SubmittedBy: NewOptionalFromPtr(b.SubmittedBy),
		AssembledBy: b.AssembledBy,
		CreatedAt:   b.CreatedAt,
	}
}
