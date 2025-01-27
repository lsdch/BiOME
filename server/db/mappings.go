package db

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lsdch/biome/models"
)

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
	Frame    string
	Mappings map[string]string
}

func (q UpdateQuery) structFragments(itemValue reflect.Value) []string {
	var fragments []string
	itemType := itemValue.Type()
	for i := range itemType.NumField() {

		f := itemType.Field(i)
		v := itemValue.Field(i)

		edgedbTag := f.Tag.Get("edgedb")
		if edgedbTag == "" {
			continue
		} else if edgedbTag == "$inline" {
			fragments = append(fragments, q.structFragments(v)...)
			continue
		}

		// Only consider fields that can capture not provided / null state
		if f.Type.Implements(reflect.TypeFor[models.OptionalNullable]()) {
			value := v.Interface().(models.OptionalNullable)
			// No value provided in JSON: do not update DB value
			if !value.HasValue() {
				continue
			}
			// Field is null: update DB value to empty set {}
			fragment := fmt.Sprintf("%s := {}", edgedbTag)
			// Field is not null: update DB value using provided mapping
			if !value.IsNull() {
				if mapping, ok := q.Mappings[edgedbTag]; ok {
					fragment = fmt.Sprintf("%s := %s", edgedbTag, mapping)
				} else {
					// Ignore field if mapping was not provided
					continue
				}
			}

			fragments = append(fragments, fragment)
		}
	}
	return fragments
}

func (q UpdateQuery) Fragments(item any) []string {
	var fragments []string

	itemValue := reflect.ValueOf(item)
	fragments = append(fragments, q.structFragments(itemValue)...)

	return fragments
}

func (q UpdateQuery) Query(item any) string {
	return fmt.Sprintf(q.Frame, strings.Join(q.Fragments(item), ",\n"))
}
