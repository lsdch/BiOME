package validations

import (
	"context"
	"fmt"

	"github.com/lsdch/biome/db"

	"github.com/edgedb/edgedb-go"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func checkExistence[T any](db *edgedb.Client, val T, bindings BindingEdgeDB) (exists bool) {
	query := fmt.Sprintf(
		"select exists %s filter .%s = <%s>$0",
		bindings.ObjectName, bindings.PropertyName, bindings.TypeCast,
	)
	err := db.QuerySingle(context.Background(), query, &exists, val)
	if err != nil {
		logrus.Errorf("Database query failed when validating 'exist' tag with error: %v\nQuery:%s", err, query)
		exists = false
	}
	return
}

func validateExists(fl validator.FieldLevel) bool {
	var value, kind, _ = fl.ExtractType(fl.Field())
	if val, ok := value.Interface().(edgedb.UUID); ok {
		return checkExistence(db.Client(), val, ParseEdgeDBBindingsOrDie(fl.Param(), "uuid"))
	} else if val, ok := value.Interface().(string); ok {
		return checkExistence(db.Client(), val, ParseEdgeDBBindingsOrDie(fl.Param(), "str"))
	}
	logrus.Fatalf("Unsupported type encountered while trying to validate 'exist=%s' constraint: %s", fl.Param(), kind)
	return false
}

// Validation constraint that checks for the existence of an item within the database. Only supports string-like or UUID values.
//
// Usage: `exist=<module name>::<object name>.<property name>`
var ExistValidator = CustomValidator{
	Tag:     "exist",
	Handler: validateExists,
	Message: func(fl validator.FieldError) string {
		return fmt.Sprintf("Item '%v' does not exist", fl.Value())
	},
}
