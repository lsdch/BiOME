package validations

import (
	"context"
	"darco/proto/db"
	"fmt"
	"strings"

	"github.com/edgedb/edgedb-go"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func checkUnknownValues[T any](db *edgedb.Client, values []T, binding BindingEdgeDB) ([]string, error) {
	query := fmt.Sprintf(`
	with values := array_unpack(<array<%s>>$0)
	select values except (
		select %s filter .%s in values
	).%s`,
		binding.TypeCast, binding.ObjectName, binding.PropertyName, binding.PropertyName,
	)
	var unknownValues []string
	err := db.Query(context.Background(), query, &unknownValues, values)
	if err != nil {
		logrus.Errorf(
			"Validation query failed with a database error while checking for the existence of %s objects having '%s' in %v. Error: %v",
			binding.ObjectName, binding.PropertyName, values, err)
	}
	return unknownValues, err
}

// Stores unknown values when checking for exist_all validation constraint.
//
// Key is the original field value formatted as string.
// Value is a slice of strings representing unknown values.
var unknownValuesMap = make(map[string][]string)

func registerNotFound(value any, notFound []string) {
	logrus.Infof("Validation failed for '%v'", value)
	unknownValuesMap[fmt.Sprintf("%v", value)] = notFound
}

func validateExistAll(fl validator.FieldLevel) bool {
	var value, _, _ = fl.ExtractType(fl.Field())
	var notFound []string
	var err error

	if slice, ok := value.Interface().([]string); ok {
		notFound, err = checkUnknownValues(
			db.Client(), slice, ParseEdgeDBBindingsOrDie(fl.Param(), "str"),
		)
	} else if slice, ok := value.Interface().([]edgedb.UUID); ok {
		notFound, err = checkUnknownValues(
			db.Client(), slice, ParseEdgeDBBindingsOrDie(fl.Param(), "uuid"),
		)
	} else {
		logrus.Fatalf("Unprocessable value for 'exist_all' validator : %+v", value.Interface())
	}

	if err != nil {
		logrus.Fatal(err)
	}

	if len(notFound) > 0 {
		registerNotFound(fl.Field().Interface(), notFound)
		return false
	}
	return true
}

// Validation constraint that checks for the existence of each item of a list within the database. Only supports string-like or UUID values.
//
// Usage: `exist_all=<module name>::<object name>.<property name>`
var ExistAllValidator = CustomValidator{
	Tag:     "exist_all",
	Handler: validateExistAll,
	Message: func(fl validator.FieldError) string {
		value := fl.Value()
		key := fmt.Sprintf("%v", value)
		notFound := unknownValuesMap[key]
		if len(notFound) == 0 {
			logrus.Errorf("ExistAllValidator: failed to retrieve non existing keys in field values: '%v'", value)
		} else {
			delete(unknownValuesMap, key)
		}
		return fmt.Sprintf("Item(s) not found: %s", strings.Join(notFound, ","))
	},
}
