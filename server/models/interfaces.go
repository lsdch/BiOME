package models

import (
	"github.com/edgedb/edgedb-go"
)

type Creatable[CreatedItem any] interface {
	Create(db *edgedb.Client) (CreatedItem, error)
}

type Updatable[ID any, Updated any] interface {
	Update(db *edgedb.Client, id ID) (edgedb.UUID, error)
	// Custom validations that can not be run through Gin bindings.
	// This is primarily intended to check for unique constraints or related objects existence, by returning validations.InputValidationError
	// Validate(db *edgedb.Client, id ID) error
}

type ItemFinder[ID any, Item any] func(db *edgedb.Client, id ID) (Item, error)

// type ItemUpdateInit[ID any, Item Updatable[ID, Updated], Updated any] func(db *edgedb.Client, id ID) (Item, error)
