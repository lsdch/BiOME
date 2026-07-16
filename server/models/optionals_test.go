package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestNullable(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("should return value and true when not null", func(t *testing.T) {
			n := Nullable[int]{Value: 42}
			val, ok := n.Get()
			assert.Equal(t, 42, val)
			assert.True(t, ok)
		})

		t.Run("should return zero and false when null", func(t *testing.T) {
			n := Nullable[int]{}
			n.SetNull()
			val, ok := n.Get()
			assert.Equal(t, 0, val)
			assert.False(t, ok)
		})
	})

	t.Run("SetNull", func(t *testing.T) {
		n := Nullable[string]{Value: "hello"}
		n.SetNull()
		assert.True(t, n.IsNull())
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Run("should marshal to null when IsNull", func(t *testing.T) {
			n := Nullable[int]{}
			n.SetNull()
			data, err := json.Marshal(n)
			assert.NoError(t, err)
			assert.Equal(t, "null", string(data))
		})

		t.Run("should marshal value when not null", func(t *testing.T) {
			n := Nullable[string]{Value: "test"}
			data, err := json.Marshal(n)
			assert.NoError(t, err)
			assert.Equal(t, `"test"`, string(data))
		})
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Run("should set null on null input", func(t *testing.T) {
			n := Nullable[int]{}
			err := json.Unmarshal([]byte("null"), &n)
			assert.NoError(t, err)
			assert.True(t, n.IsNull())
		})

		t.Run("should set null on empty string", func(t *testing.T) {
			n := Nullable[int]{}
			err := json.Unmarshal([]byte(`""`), &n)
			assert.NoError(t, err)
			assert.True(t, n.IsNull())
		})

		t.Run("should unmarshal valid value", func(t *testing.T) {
			n := Nullable[int]{}
			err := json.Unmarshal([]byte("42"), &n)
			assert.NoError(t, err)
			assert.False(t, n.IsNull())
			assert.Equal(t, 42, n.Value)
		})

		t.Run("should return error on invalid JSON", func(t *testing.T) {
			n := Nullable[int]{}
			err := json.Unmarshal([]byte("not_a_number"), &n)
			assert.Error(t, err) // Error is propagated
		})
	})
}

func TestOptional(t *testing.T) {
	t.Run("Constructors", func(t *testing.T) {
		t.Run("NewOptional", func(t *testing.T) {
			opt := NewOptional(42)
			assert.True(t, opt.IsSet)
			assert.Equal(t, 42, opt.Value)
		})

		t.Run("NewOptionalFromPtr", func(t *testing.T) {
			t.Run("should create set optional from non-nil pointer", func(t *testing.T) {
				val := 42
				opt := NewOptionalFromPtr(&val)
				assert.True(t, opt.IsSet)
				assert.Equal(t, 42, opt.Value)
			})

			t.Run("should create unset optional from nil pointer", func(t *testing.T) {
				opt := NewOptionalFromPtr[int](nil)
				assert.False(t, opt.IsSet)
				assert.Equal(t, 0, opt.Value)
			})
		})

		t.Run("NewOptionalFromUUID", func(t *testing.T) {
			t.Run("should create set optional from valid UUID", func(t *testing.T) {
				uuidBytes := [16]byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
				pgUUID := pgtype.UUID{Bytes: uuidBytes, Valid: true}
				opt := NewOptionalFromUUID(pgUUID)
				assert.True(t, opt.IsSet)
				assert.Equal(t, uuid.UUID(uuidBytes[:]), opt.Value)
			})

			t.Run("should create unset optional from invalid UUID", func(t *testing.T) {
				pgUUID := pgtype.UUID{Valid: false}
				opt := NewOptionalFromUUID(pgUUID)
				assert.False(t, opt.IsSet)
			})
		})

		t.Run("NewOptionalFromTimestamp", func(t *testing.T) {
			t.Run("should create set optional from valid timestamp", func(t *testing.T) {
				now := time.Now()
				pgTimestamp := pgtype.Timestamptz{Time: now, Valid: true}
				opt := NewOptionalFromTimestamp(pgTimestamp)
				assert.True(t, opt.IsSet)
				assert.Equal(t, now, opt.Value)
			})

			t.Run("should create unset optional from invalid timestamp", func(t *testing.T) {
				pgTimestamp := pgtype.Timestamptz{Valid: false}
				opt := NewOptionalFromTimestamp(pgTimestamp)
				assert.False(t, opt.IsSet)
			})
		})

	})

	t.Run("Get", func(t *testing.T) {
		t.Run("should return value and true when set", func(t *testing.T) {
			opt := Optional[int]{Value: 42, IsSet: true}
			val, ok := opt.Get()
			assert.Equal(t, 42, val)
			assert.True(t, ok)
		})

		t.Run("should return zero and false when unset", func(t *testing.T) {
			opt := Optional[int]{}
			val, ok := opt.Get()
			assert.Equal(t, 0, val)
			assert.False(t, ok)
		})
	})

	t.Run("ToPtr", func(t *testing.T) {
		t.Run("should return pointer when set", func(t *testing.T) {
			opt := Optional[int]{Value: 42, IsSet: true}
			ptr := opt.ToPtr()
			assert.NotNil(t, ptr)
			assert.Equal(t, 42, *ptr)
		})

		t.Run("should return nil when unset", func(t *testing.T) {
			opt := Optional[int]{}
			ptr := opt.ToPtr()
			assert.Nil(t, ptr)
		})
	})

	t.Run("GetWithDefault", func(t *testing.T) {
		t.Run("should return value when set", func(t *testing.T) {
			opt := Optional[int]{Value: 42, IsSet: true}
			assert.Equal(t, 42, opt.GetWithDefault(100))
		})

		t.Run("should return default when unset", func(t *testing.T) {
			opt := Optional[int]{}
			assert.Equal(t, 100, opt.GetWithDefault(100))
		})
	})

	t.Run("SetValue/Clear", func(t *testing.T) {
		t.Run("SetValue should set value and return pointer", func(t *testing.T) {
			opt := Optional[int]{}
			result := opt.SetValue(42)
			assert.True(t, opt.IsSet)
			assert.Equal(t, 42, opt.Value)
			assert.Same(t, &opt, result) // Returns pointer to self
		})

		t.Run("Clear should unset value and return pointer", func(t *testing.T) {
			opt := Optional[int]{Value: 42, IsSet: true}
			result := opt.Clear()
			assert.False(t, opt.IsSet)
			assert.Equal(t, 0, opt.Value)
			assert.Same(t, &opt, result)
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Run("should marshal to null when unset", func(t *testing.T) {
			opt := Optional[int]{}
			data, err := json.Marshal(opt)
			assert.NoError(t, err)
			assert.Equal(t, "null", string(data))
		})

		t.Run("should marshal value when set", func(t *testing.T) {
			opt := Optional[string]{Value: "test", IsSet: true}
			data, err := json.Marshal(opt)
			assert.NoError(t, err)
			assert.Equal(t, `"test"`, string(data))
		})
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Run("should set value on valid JSON", func(t *testing.T) {
			opt := Optional[int]{}
			err := json.Unmarshal([]byte("42"), &opt)
			assert.NoError(t, err)
			assert.True(t, opt.IsSet)
			assert.Equal(t, 42, opt.Value)
		})

		t.Run("should unset on empty string", func(t *testing.T) {
			opt := Optional[int]{}
			err := json.Unmarshal([]byte(`""`), &opt)
			assert.NoError(t, err)
			assert.False(t, opt.IsSet)
		})

		t.Run("should return error on invalid JSON", func(t *testing.T) {
			opt := Optional[int]{}
			err := json.Unmarshal([]byte("not_a_number"), &opt)
			assert.Error(t, err) // Error is propagated
		})
	})

	t.Run("UnmarshalYAML", func(t *testing.T) {
		t.Run("should set value on valid YAML", func(t *testing.T) {
			opt := Optional[int]{}
			err := yaml.Unmarshal([]byte("42"), &opt)
			assert.NoError(t, err)
			assert.True(t, opt.IsSet)
			assert.Equal(t, 42, opt.Value)
		})

		t.Run("should unset on null YAML", func(t *testing.T) {
			opt := Optional[int]{}
			err := yaml.Unmarshal([]byte("null"), &opt)
			assert.NoError(t, err)
			assert.False(t, opt.IsSet)
		})

		t.Run("should unset on empty YAML", func(t *testing.T) {
			opt := Optional[int]{}
			err := yaml.Unmarshal([]byte(""), &opt)
			assert.NoError(t, err)
			assert.False(t, opt.IsSet)
		})
	})

	t.Run("Missing", func(t *testing.T) {
		t.Run("should return true when unset", func(t *testing.T) {
			opt := Optional[int]{}
			assert.True(t, opt.Missing())
		})

		t.Run("should return false when set", func(t *testing.T) {
			opt := Optional[int]{IsSet: true}
			assert.False(t, opt.Missing())
		})
	})
}

func TestOptionalNull(t *testing.T) {
	t.Run("IsNull", func(t *testing.T) {
		t.Run("should return true when set and null", func(t *testing.T) {
			opt := OptionalNull[int]{}
			opt.SetNull()
			assert.True(t, opt.IsNull())
		})

		t.Run("should return false when unset", func(t *testing.T) {
			opt := OptionalNull[int]{}
			assert.False(t, opt.IsNull())
		})

		t.Run("should return false when set and not null", func(t *testing.T) {
			opt := NewOptionalNull(42)
			assert.False(t, opt.IsNull())
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("should return value when set and not null", func(t *testing.T) {
			opt := NewOptionalNull(42)
			val, ok := opt.Get()
			assert.Equal(t, 42, val)
			assert.True(t, ok)
		})

		t.Run("should return zero and false when unset", func(t *testing.T) {
			opt := OptionalNull[int]{}
			val, ok := opt.Get()
			assert.Equal(t, 0, val)
			assert.False(t, ok)
		})

		t.Run("should return zero and false when set and null", func(t *testing.T) {
			opt := OptionalNull[int]{}
			opt.SetNull()
			val, ok := opt.Get()
			assert.Equal(t, 0, val)
			assert.False(t, ok)
		})
	})

	t.Run("SetNull", func(t *testing.T) {
		opt := NewOptionalNull(42)
		opt.SetNull()
		assert.True(t, opt.IsSet)
		assert.True(t, opt.IsNull())
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Run("should marshal to null when unset", func(t *testing.T) {
			opt := OptionalNull[int]{}
			data, err := json.Marshal(opt)
			assert.NoError(t, err)
			assert.Equal(t, "null", string(data))
		})

		t.Run("should marshal to null when set and null", func(t *testing.T) {
			opt := OptionalNull[int]{}
			opt.SetNull()
			data, err := json.Marshal(opt)
			assert.NoError(t, err)
			assert.Equal(t, "null", string(data))
		})

		t.Run("should marshal value when set and not null", func(t *testing.T) {
			opt := NewOptionalNull("test")
			data, err := json.Marshal(opt)
			assert.NoError(t, err)
			assert.Equal(t, `"test"`, string(data))
		})
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Run("should set and null on null input", func(t *testing.T) {
			opt := OptionalNull[int]{}
			err := json.Unmarshal([]byte("null"), &opt)
			assert.NoError(t, err)
			assert.True(t, opt.IsSet)
			assert.True(t, opt.IsNull())
		})

		t.Run("should set and not null on empty string", func(t *testing.T) {
			opt := OptionalNull[string]{}
			err := json.Unmarshal([]byte(`""`), &opt)
			assert.NoError(t, err)
			assert.True(t, opt.IsSet)
			assert.False(t, opt.IsNull())
		})

		t.Run("should set and unmarshal value on valid JSON", func(t *testing.T) {
			opt := OptionalNull[int]{}
			err := json.Unmarshal([]byte("42"), &opt)
			assert.NoError(t, err)
			assert.True(t, opt.IsSet)
			assert.False(t, opt.IsNull())
			assert.Equal(t, 42, opt.Value)
		})

		t.Run("should return error on invalid JSON", func(t *testing.T) {
			opt := OptionalNull[int]{}
			err := json.Unmarshal([]byte("not_a_number"), &opt)
			assert.Error(t, err) // Error is propagated
		})
	})

	t.Run("Three states", func(t *testing.T) {
		// Test the three distinct states: absent, null, and non-null
		t.Run("absent state", func(t *testing.T) {
			opt := OptionalNull[int]{}
			assert.False(t, opt.IsSet)
			assert.False(t, opt.IsNull())
			val, ok := opt.Get()
			assert.Equal(t, 0, val)
			assert.False(t, ok)
		})

		t.Run("null state", func(t *testing.T) {
			opt := OptionalNull[int]{}
			opt.SetNull()
			assert.True(t, opt.IsSet)
			assert.True(t, opt.IsNull())
			val, ok := opt.Get()
			assert.Equal(t, 0, val)
			assert.False(t, ok)
		})

		t.Run("non-null state", func(t *testing.T) {
			opt := NewOptionalNull(42)
			assert.True(t, opt.IsSet)
			assert.False(t, opt.IsNull())
			val, ok := opt.Get()
			assert.Equal(t, 42, val)
			assert.True(t, ok)
		})
	})
}

func TestFake(t *testing.T) {
	t.Run("Optional", func(t *testing.T) {
		f := gofakeit.New(0)
		opt, err := NewOptional(0).Fake(f)
		assert.NoError(t, err)
		optTyped, ok := opt.(Optional[int])
		assert.True(t, ok)
		// IsSet is randomized by Fake
		if optTyped.IsSet {
			assert.NotEqual(t, 0, optTyped.Value) // Fake generates non-zero values
		}
	})

	t.Run("OptionalNull", func(t *testing.T) {
		f := gofakeit.New(0)
		opt := NewOptionalNull(0)
		optFake, err := opt.Fake(f)
		assert.NoError(t, err)
		optTyped, ok := optFake.(OptionalNull[int])
		assert.True(t, ok)
		// IsSet is randomized by Fake
		if optTyped.IsSet {
			if !optTyped.IsNull() {
				assert.NotEqual(t, 0, optTyped.Value)
			}
		}
	})
}
