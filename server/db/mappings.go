package db

import (
	"darco/proto/models"
	"fmt"
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
	Frame string
	Set   map[models.OptionalNullable]FieldMapping
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
