package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidationFunc func(any) error

// Validate validates struct fields using `validate` tag and optional custom validators.
func Validate(input any, custom map[string]ValidationFunc) error {
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	if t.Kind() != reflect.Struct {
		return errors.New("validation input must be a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		name := field.Name

		// Built-in tag-based validation
		tag := field.Tag.Get("validate")
		if tag != "" {
			rules := strings.Split(tag, ",")
			for _, rule := range rules {
				rule = strings.TrimSpace(rule)

				if rule == "required" {
					if isEmpty(value) {
						return fmt.Errorf("field %s is required", name)
					}
				} else if strings.HasPrefix(rule, "min=") {
					minVal := strings.TrimPrefix(rule, "min=")
					if err := checkMin(value, minVal, name); err != nil {
						return err
					}
				} else if strings.HasPrefix(rule, "max=") {
					maxVal := strings.TrimPrefix(rule, "max=")
					if err := checkMax(value, maxVal, name); err != nil {
						return err
					}
				}
			}
		}

		// Custom validator
		if custom != nil {
			if fn, ok := custom[name]; ok {
				if err := fn(value); err != nil {
					return fmt.Errorf("field %s: %w", name, err)
				}
			}
		}
	}

	return nil
}

func isEmpty(val any) bool {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	default:
		return val == nil
	}
}

func checkMin(val any, minStr string, name string) error {
	min, err := strconv.ParseFloat(minStr, 64)
	if err != nil {
		return fmt.Errorf("invalid min value on %s", name)
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Int, reflect.Int64:
		if float64(v.Int()) < min {
			return fmt.Errorf("field %s must be >= %v", name, min)
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() < min {
			return fmt.Errorf("field %s must be >= %v", name, min)
		}
	default:
		return fmt.Errorf("field %s: min not supported for type %s", name, v.Kind())
	}
	return nil
}

func checkMax(val any, maxStr string, name string) error {
	max, err := strconv.ParseFloat(maxStr, 64)
	if err != nil {
		return fmt.Errorf("invalid max value on %s", name)
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Int, reflect.Int64:
		if float64(v.Int()) > max {
			return fmt.Errorf("field %s must be <= %v", name, max)
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() > max {
			return fmt.Errorf("field %s must be <= %v", name, max)
		}
	default:
		return fmt.Errorf("field %s: max not supported for type %s", name, v.Kind())
	}
	return nil
}
