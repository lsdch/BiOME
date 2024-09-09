package models

import (
	"github.com/edgedb/edgedb-go"
)

// Creatable items can be inserted to the database
type Creatable[CreatedItem any] interface {
	Create(db edgedb.Executor) (CreatedItem, error)
}

// Updatable items can be updated in the database
type Updatable[ID any, UpdatedID any] interface {
	// Returning [edgedb.UUID] because of a bug in edgedb
	// which do not return up-to-date items after `select (update ... )` statement
	Update(db edgedb.Executor, id ID) (UpdatedID, error)
}

// ItemFinder functions fetch an item from the database using an identifier having a generic type
type ItemFinder[ID any, Item any] func(db edgedb.Executor, id ID) (Item, error)
