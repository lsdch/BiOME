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

var unknownValuesMap = make(map[string][]string)

func validateExistAll(fl validator.FieldLevel) bool {
	var value, _, _ = fl.ExtractType(fl.Field())
	var not_found []string
	var err error

	if slice, ok := value.Interface().([]string); ok {
		not_found, err = checkUnknownValues(
			db.Client(), slice, ParseEdgeDBBindings(fl.Param(), "str"),
		)
	} else if slice, ok := value.Interface().([]edgedb.UUID); ok {
		not_found, err = checkUnknownValues(
			db.Client(), slice, ParseEdgeDBBindings(fl.Param(), "uuid"),
		)
	}

	if err != nil {
		return false
	}
	if len(not_found) > 0 {
		logrus.Errorf("Validation failed for '%v'", fl.Field().Interface())
		unknownValuesMap[fmt.Sprintf("%v", fl.Field().Interface())] = not_found
		return false
	}
	return true
}

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
