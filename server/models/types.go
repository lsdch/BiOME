package models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

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

type FieldMapping struct {
	Field string
	Value string
}

func (f FieldMapping) String(isNull bool) string {
	if isNull {
		return fmt.Sprintf("%s := {}", f.Field)
	} else {
		return fmt.Sprintf("%s := %s", f.Field, f.Value)
	}
}

type UpdateQuery struct {
	Frame string
	Set   map[OptionalNullable]FieldMapping
}

func (q UpdateQuery) Fragments() []string {
	var fragments []string
	for opt, fragment := range q.Set {
		if opt.HasValue() {
			fragments = append(fragments, fragment.String(opt.IsNull()))
		}
	}
	return fragments
}

func (q UpdateQuery) Query() string {
	return fmt.Sprintf(q.Frame, strings.Join(q.Fragments(), ",\n"))
}

type Optional[T any] struct {
	edgedb.Optional `json:"-"`
	Value           T `edgedb:"$inline"`
}

func (o *Optional[T]) MarshalJSON() ([]byte, error) {
	if o.Missing() {
		return []byte("null"), nil
	}
	return json.Marshal(o.Value)
}

func (o *Optional[T]) Schema(r huma.Registry) *huma.Schema {
	contentName := huma.DefaultSchemaNamer(reflect.TypeOf(o.Value), "")
	name := fmt.Sprintf("Optional%s", contentName)
	if r.Map()[name] == nil {
		s := r.Schema(reflect.TypeOf(o.Value), false, "")
		s.Nullable = true
		r.Map()[name] = s
	}
	return &huma.Schema{Ref: fmt.Sprintf("#/components/schemas/%s", name)}
}
