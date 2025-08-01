package utils

import (
	"fmt"
	"reflect"
)

// Merge performs a shallow merge of non-zero fields from override into base.
func Merge[T any](base T, override T) (T, error) {
	baseVal := reflect.ValueOf(&base).Elem()
	overrideVal := reflect.ValueOf(override)
	if overrideVal.Kind() == reflect.Ptr {
		overrideVal = overrideVal.Elem()
	}

	if baseVal.Kind() != reflect.Struct || overrideVal.Kind() != reflect.Struct {
		return base, fmt.Errorf("both base and override must be structs")
	}

	baseType := baseVal.Type()

	for i := 0; i < baseVal.NumField(); i++ {
		field := baseVal.Field(i)
		if !field.CanSet() {
			continue
		}
		ovrField := overrideVal.FieldByName(baseType.Field(i).Name)
		if !ovrField.IsValid() {
			continue
		}
		if !isZeroValue(ovrField) {
			field.Set(ovrField)
		}
	}
	return base, nil
}

// MergeDeep performs a deep recursive merge of override into base.
func MergeDeep[T any](base T, override T) (T, error) {
	baseVal := reflect.ValueOf(&base).Elem()
	overrideVal := reflect.ValueOf(override)
	if overrideVal.Kind() == reflect.Ptr {
		overrideVal = overrideVal.Elem()
	}

	if baseVal.Kind() != reflect.Struct || overrideVal.Kind() != reflect.Struct {
		return base, fmt.Errorf("both base and override must be structs")
	}

	err := mergeDeepRecursive(baseVal, overrideVal)
	if err != nil {
		return base, err
	}
	return base, nil
}

// mergeDeepRecursive does field-by-field recursive merging.
func mergeDeepRecursive(dst, src reflect.Value) error {
	for i := 0; i < dst.NumField(); i++ {
		dstField := dst.Field(i)
		if !dstField.CanSet() {
			continue
		}

		srcField := src.Field(i)
		if !srcField.IsValid() {
			continue
		}

		switch dstField.Kind() {
		case reflect.Struct:
			err := mergeDeepRecursive(dstField, srcField)
			if err != nil {
				return err
			}

		case reflect.Ptr:
			if srcField.IsNil() {
				continue
			}
			if dstField.IsNil() {
				dstField.Set(reflect.New(dstField.Type().Elem()))
			}
			if dstField.Elem().Kind() == reflect.Struct {
				err := mergeDeepRecursive(dstField.Elem(), srcField.Elem())
				if err != nil {
					return err
				}
			} else {
				dstField.Set(srcField)
			}

		default:
			if !isZeroValue(srcField) {
				dstField.Set(srcField)
			}
		}
	}
	return nil
}

// isZeroValue checks if a reflect.Value is zero (default-initialized).
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Interface, reflect.Ptr, reflect.Func:
		return v.IsNil()
	}
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
