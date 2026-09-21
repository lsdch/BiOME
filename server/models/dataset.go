package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/types"
)

type Dataset struct {
	ID          types.ULID       `json:"id"`
	Label       string           `json:"label"`
	Slug        string           `json:"slug"`
	Description Optional[string] `json:"description,omitzero"`
	Pinned      bool             `json:"pinned"`
	OwnerID     uuid.UUID        `json:"owner_id"`
	CreatedAt   time.Time        `json:"created_at"`
}

func DatasetFromDB(d biomedb.Dataset) Dataset {
	return Dataset{
		ID:          d.ID,
		Label:       d.Label,
		Slug:        d.Slug,
		Description: NewOptionalFromPtr(d.Description),
		Pinned:      d.Pinned,
		OwnerID:     d.OwnerID,
		CreatedAt:   d.CreatedAt,
	}
}

func (d Dataset) WithMaintainers(maintainers []User) DatasetWithMaintainers {
	return DatasetWithMaintainers{
		Dataset:     d,
		Maintainers: maintainers,
	}
}

type DatasetWithMaintainers struct {
	Dataset
	Maintainers []User `json:"maintainers"`
}

func (d DatasetWithMaintainers) WithSummary(occurrenceCount, samplingCount, importBatchCount int64) DatasetWithSummary {
	return DatasetWithSummary{
		DatasetWithMaintainers: d,
		OccurrenceCount:        occurrenceCount,
		SamplingCount:          samplingCount,
		ImportBatchCount:       importBatchCount,
	}
}

type DatasetWithSummary struct {
	DatasetWithMaintainers
	OccurrenceCount  int64 `json:"occurrence_count"`
	SamplingCount    int64 `json:"sampling_count"`
	ImportBatchCount int64 `json:"import_batch_count"`
}

type DatasetInput struct {
	Label       string           `json:"label"`
	Description Optional[string] `json:"description,omitzero"`
	Pinned      bool             `json:"pinned"`
	Public      bool             `json:"public"`
	Maintainers []uuid.UUID      `json:"maintainers,omitempty"`
}

func (d DatasetInput) ToParams(ownerID uuid.UUID) biomedb.CreateDatasetParams {
	return biomedb.CreateDatasetParams{
		ULID:        types.MakeULID(),
		Label:       d.Label,
		Slug:        slug.Make(d.Label),
		Description: d.Description.ToPtr(),
		Pinned:      d.Pinned,
		IsPublic:    d.Public,
		OwnerID:     ownerID,
	}
}
