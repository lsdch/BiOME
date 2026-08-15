package imports

import (
	"reflect"
	"sync"
)

type csvFieldConfigurator func(any)

func buildCSVConfigurator[T any]() csvFieldConfigurator {
	t := reflect.TypeOf((*T)(nil)).Elem()

	var setters []func(reflect.Value)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		sep := field.Tag.Get("csvsep")
		if sep == "" {
			continue
		}

		index := i

		setters = append(setters, func(v reflect.Value) {
			field := v.Field(index)

			if setter, ok := field.Addr().Interface().(interface {
				SetSeparator(string)
			}); ok {
				setter.SetSeparator(sep)
			}
		})
	}

	return func(v any) {
		rv := reflect.ValueOf(v).Elem()

		for _, setter := range setters {
			setter(rv)
		}
	}
}

var csvConfigurators sync.Map

func getCSVConfigurator[T any]() csvFieldConfigurator {
	t := reflect.TypeOf((*T)(nil)).Elem()

	if c, ok := csvConfigurators.Load(t); ok {
		return c.(csvFieldConfigurator)
	}

	c := buildCSVConfigurator[T]()
	csvConfigurators.Store(t, c)

	return c
}
