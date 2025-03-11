package validations

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lsdch/biome/db"

	"github.com/sirupsen/logrus"
)

// Specifies a field in Gel and how to cast a value for its type
type GelBinding struct {
	ObjectName   string // e.g. 'people::Organisation'
	PropertyName string // e.g. 'code'
	TypeCast     string // e.g. 'str'
}

// Validates that the value is unique
// with respect to the specified Gel field binding.
func (b *GelBinding) UniqueQuery(value any) bool {
	query := fmt.Sprintf("select exists (select %s filter .%s = <%s>$0)", b.ObjectName, b.PropertyName, b.TypeCast)
	var exists bool
	err := db.Client().QuerySingle(context.Background(), query, &exists, value)
	if err != nil {
		logrus.Errorf("Unique validation query failed: %v with query %s", err, query)
	}
	return !exists
}

/*
Parses the path to an object property in Gel,
and creates a binding specification with a typecast so that it can be queried
with a parameter.

See https://docs.geldata.com/reference/stdlib/type#operator::cast

Expected format is <module name>::<object name>.<property name>

Example: path = "people::User.email" ; typeCast = "str"
*/
func ParseGelBindings(path string, typeCast string) (*GelBinding, error) {
	var gelBinding = strings.SplitN(path, ".", 2)
	if len(gelBinding) != 2 {
		err := fmt.Sprintf("Invalid field specification for Gel binding. Expected format like '{module name}::{object name}.{property name}', got %s", path)
		return nil, errors.New(err)
	}
	return &GelBinding{
		ObjectName:   gelBinding[0],
		PropertyName: gelBinding[1],
		TypeCast:     typeCast,
	}, nil
}

func ParseGelBindingsOrDie(path string, typeCast string) GelBinding {
	bindings, err := ParseGelBindings(path, typeCast)
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	return *bindings
}
