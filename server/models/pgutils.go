package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func PtrIfZero[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}

func ValueOrZero[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// TextToPg converts a string to a pgtype.Text.
func TextToPg(v string) pgtype.Text {
	return pgtype.Text{String: v, Valid: true}
}

func TextPtr(v *string) pgtype.Text {
	if v == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *v, Valid: true}
}

func Int4(v int32) pgtype.Int4 {
	if v == 0 {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: v, Valid: true}
}

func Int4Ptr(v *int32) pgtype.Int4 {
	if v == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: *v, Valid: true}
}

func TextOpt(v MaybeGet[string]) pgtype.Text {
	val, ok := v.Get()
	if !ok {
		return pgtype.Text{}
	}
	return pgtype.Text{String: val, Valid: true}
}

func Int4Opt(v MaybeGet[int32]) pgtype.Int4 {
	val, ok := v.Get()
	if !ok {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: val, Valid: true}
}

func UUIDToPg(v uuid.UUID) (u pgtype.UUID) {
	if v == uuid.Nil {
		return u
	}
	_ = u.Scan(v.String())
	return u
}

func UUIDOpt(v MaybeGet[uuid.UUID]) pgtype.UUID {
	val, ok := v.Get()
	if !ok {
		return pgtype.UUID{}
	}
	var u pgtype.UUID
	_ = u.Scan(val.String())
	return u
}
