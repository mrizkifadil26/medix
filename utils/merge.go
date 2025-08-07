package utils

import (
	"fmt"
	"reflect"
)

type MergeOptions struct {
	Overwrite bool // true = overwrite base, false = fill only
	Recursive bool // true = deep merge, false = shallow
}

type MergeTypeError struct {
	Reason string
}

func (e *MergeTypeError) Error() string {
	return "merge type error: " + e.Reason
}

func MergeDefault[T any](base, patch T) (T, error) {
	return Merge(base, patch, MergeOptions{})
}

func MergeOverwrite[T any](base, patch T) (T, error) {
	return Merge(base, patch, MergeOptions{Overwrite: true})
}

func MergeDeep[T any](base, patch T) (T, error) {
	return Merge(base, patch, MergeOptions{Recursive: true})
}

func MergeDeepOverwrite[T any](base, patch T) (T, error) {
	return Merge(base, patch, MergeOptions{Recursive: true, Overwrite: true})
}

// Merge performs a shallow merge of non-zero fields from override into base.
func Merge[T any](base T, override T, opts MergeOptions) (T, error) {
	baseVal := reflect.ValueOf(&base).Elem()
	overrideVal := reflect.ValueOf(override)
	if overrideVal.Kind() == reflect.Ptr {
		overrideVal = overrideVal.Elem()
	}

	if baseVal.Kind() != reflect.Struct || overrideVal.Kind() != reflect.Struct {
		return base, fmt.Errorf(
			"merge error: both base and override must be structs, got base: %s, override: %s",
			baseVal.Kind(), overrideVal.Kind(),
		)
	}

	if opts.Recursive {
		if err := mergeRecursive(baseVal, overrideVal, opts.Overwrite); err != nil {
			return base, err
		}
	} else {
		if err := mergeShallow(baseVal, overrideVal, opts.Overwrite); err != nil {
			return base, err
		}
	}

	return base, nil
}

func mergeShallow(dst, src reflect.Value, overwrite bool) error {
	dstType := dst.Type()

	for i := 0; i < dst.NumField(); i++ {
		dstField := dst.Field(i)
		if !dstField.CanSet() {
			continue
		}

		fieldName := dstType.Field(i).Name
		srcField := src.FieldByName(fieldName)
		if !srcField.IsValid() {
			continue
		}

		if overwrite {
			if !isZeroValue(srcField) {
				dstField.Set(srcField)
			}
		} else {
			if isZeroValue(dstField) && !isZeroValue(srcField) {
				dstField.Set(srcField)
			}
		}
	}
	return nil
}

func mergeRecursive(dst, src reflect.Value, overwrite bool) error {
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
			if err := mergeRecursive(dstField, srcField, overwrite); err != nil {
				return err
			}

		case reflect.Ptr:
			if srcField.IsNil() {
				continue
			}

			if dstField.IsNil() {
				if srcField.Elem().Kind() == reflect.Struct {
					// Allocate struct pointer if patch has a struct
					newStruct := reflect.New(srcField.Type().Elem())
					dstField.Set(newStruct)
				} else {
					// Scalar pointer (*string, *bool, etc.)
					if !overwrite {
						dstField.Set(srcField) // base is nil, patch has value
					} else if !isZeroValue(srcField) {
						dstField.Set(srcField)
					}
					continue
				}
			}

			if dstField.Elem().Kind() == reflect.Struct {
				// Recurse into struct pointer
				if err := mergeRecursive(dstField.Elem(), srcField.Elem(), overwrite); err != nil {
					return err
				}
			} else {
				// Scalar pointer merge
				if !overwrite {
					if isZeroValue(dstField) && !isZeroValue(srcField) {
						dstField.Set(srcField)
					}
				} else {
					if !isZeroValue(srcField) {
						dstField.Set(srcField)
					}
				}
			}

		case reflect.Slice, reflect.Map:
			if !overwrite {
				if isZeroValue(dstField) && !isZeroValue(srcField) {
					dstField.Set(srcField)
				}
			} else {
				if !isZeroValue(srcField) {
					dstField.Set(srcField)
				}
			}

		default:
			if !overwrite {
				if isZeroValue(dstField) && !isZeroValue(srcField) {
					dstField.Set(srcField)
				}
			} else {
				if !isZeroValue(srcField) {
					dstField.Set(srcField)
				}
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
	default:
		zero := reflect.Zero(v.Type()).Interface()
		return reflect.DeepEqual(v.Interface(), zero)
	}
}
