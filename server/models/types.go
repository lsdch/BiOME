package models

import (
	"github.com/edgedb/edgedb-go"
)

// Creatable items can be inserted to the database
type Creatable[CreatedItem any] interface {
	Create(db edgedb.Executor) (CreatedItem, error)
}

// Updatable items can be updated in the database
type Updatable[ID any, Updated any] interface {
	Save(db edgedb.Executor, id ID) (Updated, error)
}

// ItemFinder functions fetch an item from the database using an identifier having a generic type
type ItemFinder[ID any, Item any] func(db edgedb.Executor, id ID) (Item, error)
