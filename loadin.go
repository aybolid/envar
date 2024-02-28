package envar

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func loadInStruct(envMap map[string]string, target any) (err error) {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Pointer || targetValue.Elem().Kind() != reflect.Struct {
		err = fmt.Errorf("target must be a pointer to a struct")
	}

	targetStruct := targetValue.Elem()

	for key, value := range envMap {
		field := targetStruct.FieldByNameFunc(func(s string) bool {
			return s == key
		})
		if !field.IsValid() {
			continue
		}

		if !field.CanSet() {
			err = fmt.Errorf("field %q is not settable", key)
			return
		}

		if reflect.ValueOf(value).Type() != field.Type() {
			converted, err := convertValue(value, field.Type())
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(converted))
		} else {
			field.Set(reflect.ValueOf(value))
		}
	}

	return
}

func convertValue(value string, fieldType reflect.Type) (converted any, err error) {
	switch fieldType.String() {
	case "int":
		converted, err = strconv.Atoi(value)
		return
	case "[]string":
		converted = strings.Split(value, ",")
		return
	case "[]int":
		convertedSlice := []int{}
		for _, element := range strings.Split(value, ",") {
			parsed, err := strconv.Atoi(element)
			if err != nil {
				return nil, err
			}
			convertedSlice = append(convertedSlice, parsed)
		}
		converted = convertedSlice
		fmt.Println(converted)
		return
	default:
		err = errors.New("unsupported field type")
		return
	}
}
