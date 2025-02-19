package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
)

type Optional[T any] struct {
	edgedb.Optional `json:"-"`
	Value           T `edgedb:"$inline"`
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if o.Missing() {
		return json.Marshal(nil)
	}
	return json.Marshal(o.Value)
}

func (o *Optional[T]) Schema(r huma.Registry) *huma.Schema {
	contentName := huma.DefaultSchemaNamer(reflect.TypeOf(o.Value), "")
	name := fmt.Sprintf("Optional%s", contentName)
	if r.Map()[name] == nil {
		s := *r.Schema(reflect.TypeOf(o.Value), false, "")
		s.Nullable = true
		r.Map()[name] = &s
	}
	return &huma.Schema{Ref: fmt.Sprintf("#/components/schemas/%s", name)}
}

type OptionalNullable interface {
	HasValue() bool
	IsNull() bool
}

type Nullable[T any] struct {
	Null  bool
	Value T
}

func (n Nullable[T]) HasValue() bool {
	return true
}
func (n Nullable[T]) IsNull() bool {
	return n.Null
}

func (n Nullable[T]) Schema(r huma.Registry) *huma.Schema {
	s := r.Schema(reflect.TypeOf(n.Value), true, "")
	s.Nullable = true
	return s
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.Null {
		return json.Marshal(nil)
	}
	return json.Marshal(n.Value)
}

func (o *Nullable[T]) UnmarshalJSON(b []byte) error {
	if len(b) > 0 {
		if bytes.Equal(b, []byte("null")) || bytes.Equal(b, []byte("")) {
			o.Null = true
			return nil
		}
		return json.Unmarshal(b, &o.Value)
	}
	return nil
}

type OptionalInput[T any] struct {
	IsSet bool
	Value T
}

func NewOptionalInput[T any](value T) OptionalInput[T] {
	return OptionalInput[T]{
		Value: value,
		IsSet: true,
	}
}

func (o OptionalInput[T]) Get() (T, bool) {
	return o.Value, o.IsSet
}

func (o OptionalInput[T]) HasValue() bool {
	return o.IsSet
}

func (o *OptionalInput[T]) SetValue(value T) {
	o.IsSet = true
	o.Value = value
}

func (o OptionalInput[T]) IsNull() bool {
	return false
}

func (o *OptionalInput[T]) Fake(f *gofakeit.Faker) (any, error) {
	var value T
	if err := f.Struct(&value); err != nil {
		return nil, err
	}
	return OptionalInput[T]{
		IsSet: f.Bool(),
		Value: value,
	}, nil
}

var _ gofakeit.Fakeable = (*OptionalInput[any])(nil)

func (o OptionalInput[T]) Schema(r huma.Registry) *huma.Schema {
	return r.Schema(reflect.TypeOf(o.Value), true, "")
}

func (o OptionalInput[T]) MarshalJSON() ([]byte, error) {
	if !o.IsSet {
		return json.Marshal(nil)
	}
	return json.Marshal(o.Value)
}

func (o *OptionalInput[T]) UnmarshalJSON(b []byte) error {
	if len(b) > 0 && string(b) != `""` {
		o.IsSet = true
		return json.Unmarshal(b, &o.Value)
	}
	o.IsSet = false
	return nil
}

func (o OptionalInput[T]) MarshalEdgeDBStr() ([]byte, error) {
	if !o.IsSet {
		return []byte("null"), nil
	}
	return json.Marshal(o.Value)
}

// Implementation of huma.ParamWrapper interface for request parameters binding
func (o *OptionalInput[T]) Receiver() reflect.Value {
	return reflect.ValueOf(o).Elem().Field(0)
}

// Implementation of huma.ParamReactor interface for request parameters binding
func (o *OptionalInput[T]) OnParamSet(isSet bool, parsed any) {
	o.IsSet = isSet
}

// OptionalNull is a field which can be omitted from the input,
// set to `null`, or set to a value. Each state is tracked and can
// be checked for in handling code.
type OptionalNull[T any] struct {
	Null bool
	OptionalInput[T]
}

func NewOptionalNull[T any](value T) OptionalNull[T] {
	return OptionalNull[T]{
		OptionalInput: OptionalInput[T]{
			Value: value,
			IsSet: true,
		},
		Null: false,
	}
}

func (o OptionalNull[T]) IsNull() bool {
	return o.Null
}

func (o OptionalNull[T]) Schema(r huma.Registry) *huma.Schema {
	s := r.Schema(reflect.TypeOf(o.Value), true, "")
	s.Nullable = true
	return s
}

func (o OptionalNull[T]) MarshalJSON() ([]byte, error) {
	if (!o.IsSet) || o.Null {
		return json.Marshal(nil)
	}
	return json.Marshal(o.Value)
}

func (o *OptionalNull[T]) UnmarshalJSON(b []byte) error {
	if len(b) > 0 {
		o.IsSet = true
		if bytes.Equal(b, []byte("null")) || bytes.Equal(b, []byte("")) {
			o.Null = true
			return nil
		}
		return json.Unmarshal(b, &o.Value)
	}
	return nil
}

func (o *OptionalNull[T]) Fake(f *gofakeit.Faker) (any, error) {
	v, err := o.OptionalInput.Fake(f)
	if err != nil {
		return nil, err
	}
	return OptionalNull[T]{
		OptionalInput: v.(OptionalInput[T]),
	}, nil
}
