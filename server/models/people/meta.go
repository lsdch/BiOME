package people

import (
	"time"

	"github.com/edgedb/edgedb-go"
)

type UserShortIdentity struct {
	edgedb.Optional
	Name  string `json:"name" edgedb:"name"`
	Alias string `json:"alias" edgedb:"alias"`
}

type OptionalUser struct {
	edgedb.Optional
	User `edgedb:"$inline" json:",inline"`
}

// Metadata attached to an item in the database that track updates of the item.
type Meta struct {
	ID            edgedb.UUID             `edgedb:"id" json:"-" swaggerignore:"true"`
	Created       time.Time               `edgedb:"created" json:"created" binding:"required"`
	Modified      edgedb.OptionalDateTime `edgedb:"modified" json:"modified"`
	LastUpdated   time.Time               `edgedb:"lastUpdated" json:"last_updated"`
	CreatedByUser OptionalUser            `json:"created_by_user" edgedb:"created_by_user"`
	UpdatedByUser OptionalUser            `json:"updated_by_user" edgedb:"modified_by_user"`
	CreatedBy     UserShortIdentity       `json:"created_by" edgedb:"created_by"`
	UpdatedBy     UserShortIdentity       `json:"updated_by" edgedb:"updated_by"`
} // @name Meta
