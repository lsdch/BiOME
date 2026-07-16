package models

import (
	"time"

	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/types"
)

type Dataset struct {
	ID          types.ULID       `json:"id"`
	Label       string           `json:"label"`
	Slug        string           `json:"slug"`
	Description Optional[string] `json:"description,omitempty"`
	Pinned      bool             `json:"pinned"`
	CreatedAt   time.Time        `json:"created_at"`
}

func DatasetFromDB(d biomedb.Dataset) Dataset {
	return Dataset{
		ID:          d.ID,
		Label:       d.Label,
		Slug:        d.Slug,
		Description: NewOptionalFromPtr(d.Description),
		Pinned:      d.Pinned,
		CreatedAt:   d.CreatedAt,
	}
}
