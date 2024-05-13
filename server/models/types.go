package models

import (
	"encoding/json"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
)

// Creatable items can be inserted to the database
type Creatable[CreatedItem any] interface {
	Create(db edgedb.Executor) (CreatedItem, error)
}

// Updatable items can be updated in the database
type Updatable[ID any, Updated any] interface {
	// Returning [edgedb.UUID] because of a bug in edgedb
	// which do not return up-to-date items after `select (update ... )` statement
	Update(db edgedb.Executor, id ID) (edgedb.UUID, error)
}

// ItemFinder functions fetch an item from the database using an identifier having a generic type
type ItemFinder[ID any, Item any] func(db edgedb.Executor, id ID) (Item, error)

type Optional[T any] struct {
	edgedb.Optional
	Value T `edgedb:"$inline" json:",inline"`
}

func (o *Optional[T]) MarshalJSON() ([]byte, error) {
	if o.Missing() {
		return []byte("null"), nil
	}
	return json.Marshal(o.Value)
}

func (o *Optional[T]) Schema(r huma.Registry) *huma.Schema {
	return r.Schema(reflect.TypeOf(o.Value), true, "")
}
