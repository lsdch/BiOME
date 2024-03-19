package models

import (
	"time"

	"github.com/edgedb/edgedb-go"
)

type UserShortIdentity struct {
	edgedb.Optional
	Name  string `json:"name" edgedb:"name"`
	Alias string `json:"alias" edgedb:"alias"`
}

// Metadata attached to an item in the database that track updates of the item.
type Meta struct {
	ID          edgedb.UUID             `edgedb:"id" json:"-" swaggerignore:"true"`
	Created     time.Time               `edgedb:"created" json:"created" example:"2023-09-01T16:41:10.921097+00:00" binding:"required"`
	Modified    edgedb.OptionalDateTime `edgedb:"modified" json:"modified" example:"2023-09-02T20:39:10.218057+00:00"`
	LastUpdated time.Time               `edgedb:"lastUpdated" json:"last_updated" example:"2023-09-02T20:39:10.218057+00:00"`
	CreatedBy   UserShortIdentity       `json:"created_by" edgedb:"created_by"`
	UpdatedBy   UserShortIdentity       `json:"updated_by" edgedb:"updated_by"`
} // @name Meta
