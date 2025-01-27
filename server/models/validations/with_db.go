package validations

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lsdch/biome/db"

	"github.com/sirupsen/logrus"
)

// Specifies a field in EdgeDB and how to cast a value for its type
type BindingEdgeDB struct {
	ObjectName   string // e.g. 'people::Institution'
	PropertyName string // e.g. 'code'
	TypeCast     string // e.g. 'str'
}

// Validates that the value is unique
// with respect to the specified EdgeDB field binding.
func (b *BindingEdgeDB) UniqueQuery(value any) bool {
	query := fmt.Sprintf("select exists (select %s filter .%s = <%s>$0)", b.ObjectName, b.PropertyName, b.TypeCast)
	var exists bool
	err := db.Client().QuerySingle(context.Background(), query, &exists, value)
	if err != nil {
		logrus.Errorf("Unique validation query failed: %v with query %s", err, query)
	}
	return !exists
}

/*
Parses the path to an object property in EdgeDB,
and creates a binding specification with a typecast so that it can be queried
with a parameter.

See https://www.edgedb.com/docs/stdlib/type#operator::cast.

Expected format is <module name>::<object name>.<property name>

Example: path = "people::User.email" ; typeCast = "str"
*/
func ParseEdgeDBBindings(path string, typeCast string) (*BindingEdgeDB, error) {
	var edgeDBBinding = strings.SplitN(path, ".", 2)
	if len(edgeDBBinding) != 2 {
		err := fmt.Sprintf("Invalid field specification for EdgeDB binding. Expected format like '{module name}::{object name}.{property name}', got %s", path)
		return nil, errors.New(err)
	}
	return &BindingEdgeDB{
		ObjectName:   edgeDBBinding[0],
		PropertyName: edgeDBBinding[1],
		TypeCast:     typeCast,
	}, nil
}

func ParseEdgeDBBindingsOrDie(path string, typeCast string) BindingEdgeDB {
	bindings, err := ParseEdgeDBBindings(path, typeCast)
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	return *bindings
}
