package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ValueOrZero[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

func TextToPg(v string) pgtype.Text {
	if v == "" {
		return pgtype.Text{}
	}
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

func TextOpt(v OptionalNullable[string]) pgtype.Text {
	val, ok := v.Get()
	if !ok {
		return pgtype.Text{}
	}
	return pgtype.Text{String: val, Valid: true}
}

func Int4Opt(v OptionalNullable[int32]) pgtype.Int4 {
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
	u.Scan(v.String())
	return u
}

func UUIDOpt(v OptionalNullable[uuid.UUID]) pgtype.UUID {
	val, ok := v.Get()
	if !ok {
		return pgtype.UUID{}
	}
	var u pgtype.UUID
	_ = u.Scan(val)
	return u
}
