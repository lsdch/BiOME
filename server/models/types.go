package models

import "github.com/geldata/gel-go/geltypes"

// Persistable items can make changes in the database,
// e.g. when inserting a record
type Persistable[Persisted any] interface {
	Save(db geltypes.Executor) (Persisted, error)
}

// PersistableWithID items can make changes in the database,
// provided an identifier, e.g. when updating a record using its UUID or code
type PersistableWithID[ID any, Persisted any] interface {
	Save(db geltypes.Executor, id ID) (Persisted, error)
}

// ItemFinder functions fetch an item from the database using an identifier having a generic type
type ItemFinder[ID any, Item any] func(db geltypes.Executor, id ID) (Item, error)

type FindOrCreate[ID any, Input Persistable[Item], Item any] struct {
	Use   []ID    `json:"use"`
	Input []Input `json:"create"`
}
