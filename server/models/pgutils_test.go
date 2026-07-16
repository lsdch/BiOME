package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValueOrZero(t *testing.T) {
	t.Run("nil pointer returns zero value", func(t *testing.T) {
		var nilPtr *string
		result := ValueOrZero(nilPtr)
		assert.Equal(t, "", result)
	})

	t.Run("non-nil pointer returns dereferenced value", func(t *testing.T) {
		value := "test"
		result := ValueOrZero(&value)
		assert.Equal(t, "test", result)
	})

	t.Run("works with int32", func(t *testing.T) {
		var nilPtr *int32
		result := ValueOrZero(nilPtr)
		assert.Equal(t, int32(0), result)

		value := int32(42)
		result = ValueOrZero(&value)
		assert.Equal(t, int32(42), result)
	})

	t.Run("works with uuid", func(t *testing.T) {
		var nilPtr *uuid.UUID
		result := ValueOrZero(nilPtr)
		assert.Equal(t, uuid.Nil, result)

		id := uuid.Must(uuid.NewV7())
		result = ValueOrZero(&id)
		assert.Equal(t, id, result)
	})
}

func TestTextToPg(t *testing.T) {
	t.Run("string returns valid pgtype.Text", func(t *testing.T) {
		result := TextToPg("hello")
		assert.True(t, result.Valid)
		assert.Equal(t, "hello", result.String)
	})

	t.Run("whitespace-only string returns valid pgtype.Text", func(t *testing.T) {
		result := TextToPg("  ")
		assert.True(t, result.Valid)
		assert.Equal(t, "  ", result.String)
	})
}

func TestTextPtr(t *testing.T) {
	t.Run("nil pointer returns invalid pgtype.Text", func(t *testing.T) {
		var nilPtr *string
		result := TextPtr(nilPtr)
		assert.False(t, result.Valid)
		assert.Equal(t, "", result.String)
	})

	t.Run("non-nil pointer returns valid pgtype.Text", func(t *testing.T) {
		value := "world"
		result := TextPtr(&value)
		assert.True(t, result.Valid)
		assert.Equal(t, "world", result.String)
	})

	t.Run("pointer to empty string returns valid pgtype.Text", func(t *testing.T) {
		value := ""
		result := TextPtr(&value)
		assert.True(t, result.Valid)
		assert.Equal(t, "", result.String)
	})
}

func TestInt4(t *testing.T) {
	t.Run("zero returns invalid pgtype.Int4", func(t *testing.T) {
		result := Int4(0)
		assert.False(t, result.Valid)
		assert.Equal(t, int32(0), result.Int32)
	})

	t.Run("non-zero returns valid pgtype.Int4", func(t *testing.T) {
		result := Int4(42)
		assert.True(t, result.Valid)
		assert.Equal(t, int32(42), result.Int32)
	})

	t.Run("negative value returns valid pgtype.Int4", func(t *testing.T) {
		result := Int4(-10)
		assert.True(t, result.Valid)
		assert.Equal(t, int32(-10), result.Int32)
	})
}

func TestInt4Ptr(t *testing.T) {
	t.Run("nil pointer returns invalid pgtype.Int4", func(t *testing.T) {
		var nilPtr *int32
		result := Int4Ptr(nilPtr)
		assert.False(t, result.Valid)
		assert.Equal(t, int32(0), result.Int32)
	})

	t.Run("non-nil pointer returns valid pgtype.Int4", func(t *testing.T) {
		value := int32(100)
		result := Int4Ptr(&value)
		assert.True(t, result.Valid)
		assert.Equal(t, int32(100), result.Int32)
	})

	t.Run("pointer to zero returns valid pgtype.Int4", func(t *testing.T) {
		value := int32(0)
		result := Int4Ptr(&value)
		assert.True(t, result.Valid)
		assert.Equal(t, int32(0), result.Int32)
	})
}

func TestTextOpt(t *testing.T) {

	t.Run("OptionalNull with value returns valid pgtype.Text", func(t *testing.T) {
		opt := NewOptionalNull("hello")
		result := TextOpt(opt)
		assert.True(t, result.Valid)
		assert.Equal(t, "hello", result.String)
	})

	t.Run("OptionalNull with null value returns invalid pgtype.Text", func(t *testing.T) {
		var opt OptionalNull[string]
		opt.SetNull()
		result := TextOpt(opt)
		assert.False(t, result.Valid)
	})

	t.Run("OptionalNull not set returns invalid pgtype.Text", func(t *testing.T) {
		var opt OptionalNull[string]
		result := TextOpt(opt)
		assert.False(t, result.Valid)
	})

	t.Run("Optional with value returns valid pgtype.Text", func(t *testing.T) {
		opt := NewOptional("optional")
		result := TextOpt(opt)
		assert.True(t, result.Valid)
		assert.Equal(t, "optional", result.String)
	})

	t.Run("Optional not set returns invalid pgtype.Text", func(t *testing.T) {
		var opt Optional[string]
		result := TextOpt(opt)
		assert.False(t, result.Valid)
	})
}

func TestInt4Opt(t *testing.T) {

	t.Run("OptionalNull with value returns valid pgtype.Int4", func(t *testing.T) {
		opt := NewOptionalNull(int32(100))
		result := Int4Opt(opt)
		assert.True(t, result.Valid)
		assert.Equal(t, int32(100), result.Int32)
	})

	t.Run("OptionalNull with null value returns invalid pgtype.Int4", func(t *testing.T) {
		var opt OptionalNull[int32]
		opt.SetNull()
		result := Int4Opt(opt)
		assert.False(t, result.Valid)
	})

	t.Run("OptionalNull not set returns invalid pgtype.Int4", func(t *testing.T) {
		var opt OptionalNull[int32]
		result := Int4Opt(opt)
		assert.False(t, result.Valid)
	})

	t.Run("Optional with value returns valid pgtype.Int4", func(t *testing.T) {
		opt := NewOptional(int32(200))
		result := Int4Opt(opt)
		assert.True(t, result.Valid)
		assert.Equal(t, int32(200), result.Int32)
	})

	t.Run("Optional not set returns invalid pgtype.Int4", func(t *testing.T) {
		var opt Optional[int32]
		result := Int4Opt(opt)
		assert.False(t, result.Valid)
	})
}

func TestUUIDToPg(t *testing.T) {
	t.Run("uuid.Nil returns invalid pgtype.UUID", func(t *testing.T) {
		result := UUIDToPg(uuid.Nil)
		assert.False(t, result.Valid)
	})

	t.Run("valid UUID returns valid pgtype.UUID", func(t *testing.T) {
		id := uuid.Must(uuid.NewV7())
		result := UUIDToPg(id)
		assert.True(t, result.Valid)
		assert.Equal(t, id.String(), result.String())
	})

	t.Run("UUID from string returns valid pgtype.UUID", func(t *testing.T) {
		id := uuid.Must(uuid.Parse("123e4567-e89b-12d3-a456-426614174000"))
		result := UUIDToPg(id)
		assert.True(t, result.Valid)
		assert.Equal(t, id.String(), result.String())
	})
}

func TestUUIDOpt(t *testing.T) {

	t.Run("OptionalNull with value returns valid pgtype.UUID", func(t *testing.T) {
		id := uuid.Must(uuid.NewV7())
		opt := NewOptionalNull(id)
		result := UUIDOpt(opt)
		assert.True(t, result.Valid)
		assert.Equal(t, id.String(), result.String())
	})

	t.Run("OptionalNull with null value returns invalid pgtype.UUID", func(t *testing.T) {
		var opt OptionalNull[uuid.UUID]
		opt.SetNull()
		result := UUIDOpt(opt)
		assert.False(t, result.Valid)
	})

	t.Run("OptionalNull not set returns invalid pgtype.UUID", func(t *testing.T) {
		var opt OptionalNull[uuid.UUID]
		result := UUIDOpt(opt)
		assert.False(t, result.Valid)
	})

	t.Run("Optional with value returns valid pgtype.UUID", func(t *testing.T) {
		id := uuid.Must(uuid.NewV7())
		opt := NewOptional(id)
		result := UUIDOpt(opt)
		assert.True(t, result.Valid)
		assert.Equal(t, id.String(), result.String())
	})

	t.Run("Optional not set returns invalid pgtype.UUID", func(t *testing.T) {
		var opt Optional[uuid.UUID]
		result := UUIDOpt(opt)
		assert.False(t, result.Valid)
	})
}

// Test edge cases and consistency
func TestEdgeCases(t *testing.T) {
	t.Run("TextToPg and TextPtr consistency with empty values", func(t *testing.T) {
		textResult := TextToPg("")
		emptyStr := ""
		ptrResult := TextPtr(&emptyStr)

		// Both should be valid since empty string is a valid value, not nil
		assert.True(t, textResult.Valid)
		assert.True(t, ptrResult.Valid)
		assert.Equal(t, textResult.String, ptrResult.String)
	})

	t.Run("Int4 zero vs nil pointer behavior", func(t *testing.T) {
		// Int4(0) returns invalid because 0 is treated as "not set"
		zeroResult := Int4(0)
		assert.False(t, zeroResult.Valid)

		// But Int4Ptr(&zero) returns valid because pointer is not nil
		zero := int32(0)
		ptrResult := Int4Ptr(&zero)
		assert.True(t, ptrResult.Valid)
	})

	t.Run("ValueOrZero preserves type default values", func(t *testing.T) {
		// For string, zero value is ""
		var nilStr *string
		assert.Equal(t, "", ValueOrZero(nilStr))

		// For int32, zero value is 0
		var nilInt *int32
		assert.Equal(t, int32(0), ValueOrZero(nilInt))

		// For uuid, zero value is uuid.Nil
		var nilUUID *uuid.UUID
		assert.Equal(t, uuid.Nil, ValueOrZero(nilUUID))
	})
}
