package models

import (
	"bytes"
	"encoding/json"
	"reflect"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/danielgtaylor/huma/v2"
	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

type MaybeGet[T any] interface {
	Get() (T, bool)
}

type Null struct {
	isNull bool
}

func (n Null) IsNull() bool {
	return n.isNull
}
func (n *Null) SetNull() {
	n.isNull = true
}

type Nullable[T any] struct {
	Null
	Value T
}

func (n Nullable[T]) Get() (T, bool) {
	return n.Value, !n.IsNull()
}

func (n Nullable[T]) Schema(r huma.Registry) *huma.Schema {
	schemaRef := r.Schema(reflect.TypeOf(n.Value), true, "")
	s := r.Schema(reflect.TypeOf(n.Value), false, "")
	schemaRef.Type = s.Type
	schemaRef.Nullable = true
	return schemaRef
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.IsNull() {
		return json.Marshal(nil)
	}
	return json.Marshal(n.Value)
}

func (o *Nullable[T]) UnmarshalJSON(b []byte) error {
	if len(b) > 0 {
		if bytes.Equal(b, []byte("null")) || bytes.Equal(b, []byte(`""`)) {
			o.SetNull()
			return nil
		}
		return json.Unmarshal(b, &o.Value)
	}
	return nil
}

type Optional[T any] struct {
	Value T
	IsSet bool
}

func NewOptional[T any](value T) Optional[T] {
	return Optional[T]{
		Value: value,
		IsSet: true,
	}
}

func NewOptionalFromPtr[T any](ptr *T) Optional[T] {
	if ptr == nil {
		var zero T
		return Optional[T]{
			Value: zero,
			IsSet: false,
		}
	}
	return Optional[T]{
		Value: *ptr,
		IsSet: true,
	}
}

func NewOptionalFromUUID(u pgtype.UUID) Optional[uuid.UUID] {
	if !u.Valid {
		var zero uuid.UUID
		return Optional[uuid.UUID]{
			Value: zero,
			IsSet: false,
		}
	}
	// ignore error since we know the UUID is valid and has 16 bytes
	UUID, _ := uuid.FromBytes(u.Bytes[:16])
	return Optional[uuid.UUID]{
		Value: UUID,
		IsSet: true,
	}
}

func MapOptional[T any, U any](opt Optional[T], f func(T) U) Optional[U] {
	if !opt.IsSet {
		return Optional[U]{}
	}
	return Optional[U]{
		Value: f(opt.Value),
		IsSet: true,
	}
}

func NewOptionalFromTimestamp(t pgtype.Timestamptz) Optional[time.Time] {
	if !t.Valid {
		var zero time.Time
		return Optional[time.Time]{
			Value: zero,
			IsSet: false,
		}
	}
	return Optional[time.Time]{
		Value: t.Time,
		IsSet: true,
	}
}

func (o Optional[T]) Get() (T, bool) {
	return o.Value, o.IsSet
}

func (o Optional[T]) ToPtr() *T {
	if o.IsSet {
		return &o.Value
	}
	return nil
}

func (o Optional[T]) GetWithDefault(value T) T {
	if o.IsSet {
		return o.Value
	}
	return value
}

func (o Optional[T]) HasValue() bool {
	return o.IsSet
}

func (o *Optional[T]) SetValue(value T) *Optional[T] {
	o.IsSet = true
	o.Value = value
	return o
}

func (o *Optional[T]) Clear() *Optional[T] {
	o.IsSet = false
	var zero T
	o.Value = zero
	return o
}

func (o Optional[T]) IsZero() bool {
	return !o.IsSet
}

func (o Optional[T]) Fake(f *gofakeit.Faker) (any, error) {
	var value T
	if err := f.Struct(&value); err != nil {
		return nil, err
	}
	return Optional[T]{
		IsSet: f.Bool(),
		Value: value,
	}, nil
}

var _ gofakeit.Fakeable = (*Optional[any])(nil)

func (o Optional[T]) Schema(r huma.Registry) *huma.Schema {
	return r.Schema(reflect.TypeOf(o.Value), true, "")
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.IsSet {
		return json.Marshal(nil)
	}
	return json.Marshal(o.Value)
}

func (o *Optional[T]) UnmarshalJSON(b []byte) error {
	if len(b) > 0 && string(b) != `""` {
		o.IsSet = true
		return json.Unmarshal(b, &o.Value)
	}
	o.IsSet = false
	return nil
}

func (o *Optional[T]) UnmarshalYAML(b []byte) error {
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		o.IsSet = false
		var zero T
		o.Value = zero
		return nil
	}

	o.IsSet = true
	return yaml.Unmarshal(b, &o.Value)
}

func (o Optional[T]) Missing() bool {
	return !o.IsSet
}

// Implementation of huma.ParamWrapper interface for request parameters binding
func (o *Optional[T]) Receiver() reflect.Value {
	return reflect.ValueOf(o).Elem().Field(0)
}

// Implementation of huma.ParamReactor interface for request parameters binding
func (o *Optional[T]) OnParamSet(isSet bool, parsed any) {
	o.IsSet = isSet
}

// OptionalNull is a field which can be omitted from the input,
// set to `null`, or set to a value. Each state is tracked and can
// be checked for in handling code.
type OptionalNull[T any] struct {
	Null
	Optional[T]
}

func NewOptionalNull[T any](value T) OptionalNull[T] {
	return OptionalNull[T]{
		Optional: Optional[T]{
			Value: value,
			IsSet: true,
		},
		Null: Null{isNull: false},
	}
}

func (o OptionalNull[T]) IsNull() bool {
	return o.IsSet && o.Null.IsNull()
}

func (o OptionalNull[T]) Get() (T, bool) {
	if o.IsSet && !o.Null.IsNull() {
		return o.Value, true
	}
	var zero T
	return zero, false
}

func (o OptionalNull[T]) Schema(r huma.Registry) *huma.Schema {
	schemaRef := r.Schema(reflect.TypeOf(o.Value), true, "")
	s := r.Schema(reflect.TypeOf(o.Value), false, "")

	if s.Ref != "" {
		sDeref := r.SchemaFromRef(s.Ref)
		if sDeref == nil {
			logrus.Errorf("Failed to follow schema reference: %s", s.Ref)
			return schemaRef
		}
		schemaRef.Type = sDeref.Type
	} else {
		schemaRef.Type = s.Type
	}
	schemaRef.Nullable = true
	return schemaRef
}

func (o *OptionalNull[T]) SetNull() {
	o.IsSet = true
	o.Null.isNull = true
}

func (o OptionalNull[T]) MarshalJSON() ([]byte, error) {
	// Due to JSON marshalling rules, we cannot distinguish between an omitted field and a field set to null.
	if (!o.IsSet) || o.Null.IsNull() {
		return json.Marshal(nil)
	}
	return json.Marshal(o.Value)
}

func (o *OptionalNull[T]) UnmarshalJSON(b []byte) error {
	if len(b) > 0 {
		o.IsSet = true
		if bytes.Equal(b, []byte("null")) || bytes.Equal(b, []byte("")) {
			o.SetNull()
			return nil
		}
		return json.Unmarshal(b, &o.Value)
	}
	return nil
}

func (o OptionalNull[T]) Fake(f *gofakeit.Faker) (any, error) {
	v, err := o.Optional.Fake(f)
	if err != nil {
		return nil, err
	}
	return OptionalNull[T]{
		Optional: v.(Optional[T]),
	}, nil
}
