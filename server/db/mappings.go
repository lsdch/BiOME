package db

import (
	"darco/proto/models"
	"fmt"
	"reflect"
	"strings"
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

func (q UpdateQuery) Fragments(item any) []string {
	var fragments []string

	itemType := reflect.TypeOf(item)
	itemValue := reflect.ValueOf(item)
	for i := range itemType.NumField() {
		f := itemType.Field(i)
		v := itemValue.Field(i)
		jsonField := strings.Split(f.Tag.Get("json"), ",")[0]
		if jsonField == "" {
			continue
		}
		if f.Type.Implements(reflect.TypeFor[models.OptionalNullable]()) {
			value := v.Interface().(models.OptionalNullable)
			if !value.HasValue() {
				continue
			}
			fragment := fmt.Sprintf("%s := {}", jsonField)
			if !value.IsNull() {
				fragment = fmt.Sprintf("%s := %s", jsonField, q.Mappings[jsonField])
			}
			fragments = append(fragments, fragment)
		}
	}
	return fragments
}

func (q UpdateQuery) Query(item any) string {
	return fmt.Sprintf(q.Frame, strings.Join(q.Fragments(item), ",\n"))
}
