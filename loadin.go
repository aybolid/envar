package envar

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func loadInStruct(envMap map[string]string, target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		return errors.New("target must be a pointer to a struct")
	}

	targetStruct := targetValue.Elem()

	for key, value := range envMap {
		field := targetStruct.FieldByName(key)
		if !field.IsValid() || !field.CanSet() {
			continue
		}

		if err := setFieldValue(field, value); err != nil {
			return fmt.Errorf("setting field %q: %v", key, err)
		}
	}

	return nil
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		val, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(val))
	case reflect.Slice:
		elemType := field.Type().Elem()
		switch elemType.Kind() {
		case reflect.String:
			field.Set(reflect.ValueOf(strings.Split(value, ",")))
		case reflect.Int:
			strValues := strings.Split(value, ",")
			slice := reflect.MakeSlice(field.Type(), len(strValues), len(strValues))
			for i, str := range strValues {
				val, err := strconv.Atoi(str)
				if err != nil {
					return err
				}
				slice.Index(i).SetInt(int64(val))
			}
			field.Set(slice)
		default:
			return errors.New("unsupported slice element type")
		}
	default:
		return errors.New("unsupported field type")
	}
	return nil
}
