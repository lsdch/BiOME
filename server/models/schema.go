package models

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

func UseSchema[T any](r huma.Registry, schema *huma.Schema) *huma.Schema {
	// parentSchema := r.Schema(reflect.TypeFor[T](), true, "")
	parentSchema := huma.SchemaFromType(r, reflect.TypeFor[T]())
	for name, s := range schema.Properties {
		*s = *parentSchema.Properties[name]
	}
	return schema
}
