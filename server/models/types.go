package models

import (
	"github.com/edgedb/edgedb-go"
)

// Persistable items can make changes in the database,
// e.g. when inserting a record
type Persistable[Persisted any] interface {
	Save(db edgedb.Executor) (Persisted, error)
}

// PersistableWithID items can make changes in the database,
// provided an identifier, e.g. when updating a record using its UUID or code
type PersistableWithID[ID any, Persisted any] interface {
	Save(db edgedb.Executor, id ID) (Persisted, error)
}

// ItemFinder functions fetch an item from the database using an identifier having a generic type
type ItemFinder[ID any, Item any] func(db edgedb.Executor, id ID) (Item, error)
