package envar

import (
	"fmt"
	"reflect"
)

// TODO: Improve!!!
func loadInStruct(envMap map[string]string, target any) (err error) {
	value := reflect.ValueOf(target)
	if value.Kind() != reflect.Pointer || value.Elem().Kind() != reflect.Struct {
		err = fmt.Errorf("target must be a pointer to a struct")
	}

	structValue := value.Elem()

	for key, value := range envMap {
		field := structValue.FieldByNameFunc(func(s string) bool {
			// TODO: Improve field identification?
			return s == key
		})
		if !field.IsValid() {
			continue
		}

		if !field.CanSet() {
			err = fmt.Errorf("field %q is not settable", key)
			return
		}

		convertedValue := reflect.ValueOf(value)
		if convertedValue.Type() != field.Type() {
			convertedValue = convertedValue.Convert(field.Type())
		}

		field.Set(convertedValue)
	}
	return
}
