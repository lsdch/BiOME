package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

func OptionalInputDecodeHook(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	// Only apply to OptionalInput[T] structs
	if t.Kind() != reflect.Struct || !strings.HasPrefix(t.String(), "models.Optional") {
		return data, nil
	}

	out := reflect.New(t).Elem()
	valField := out.FieldByName("Value")
	isSetField := out.FieldByName("IsSet")

	if !valField.IsValid() || !isSetField.IsValid() {
		return data, nil
	}

	dataVal := reflect.ValueOf(data)

	// Case 1: scalar value (string, int, etc.)
	if dataVal.Type().AssignableTo(valField.Type()) {
		valField.Set(dataVal)
		isSetField.SetBool(true)
		return out.Interface(), nil
	}

	// Case 2: map/object
	if m, ok := data.(map[string]interface{}); ok {
		err := mapstructure.Decode(m, out.Addr().Interface())
		if err != nil {
			return nil, err
		}
		return out.Interface(), nil
	}

	return nil, fmt.Errorf("cannot decode %v into OptionalInput", t)
}
