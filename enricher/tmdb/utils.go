package tmdb

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type ParamEncoder interface {
	Params() url.Values
}

func buildParams(query any) url.Values {
	values := url.Values{}
	v := indirectValue(query) // handles pointer or struct
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		key := field.Tag.Get("param")
		if key == "" {
			continue
		}

		val := v.Field(i)
		switch val.Kind() {
		case reflect.String:
			if val.String() != "" {
				values.Set(key, val.String())
			}
		case reflect.Int, reflect.Int64:
			if val.Int() != 0 {
				values.Set(key, strconv.FormatInt(val.Int(), 10))
			}
		case reflect.Bool:
			values.Set(key, strconv.FormatBool(val.Bool()))
		case reflect.Slice:
			if val.Len() > 0 {
				strs := make([]string, val.Len())
				for j := 0; j < val.Len(); j++ {
					strs[j] = valueToString(val.Index(j))
				}
				values.Set(key, strings.Join(strs, ","))
			}
		}
	}

	return values

}

// Converts a reflect.Value to string
func valueToString(val reflect.Value) string {
	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())
	default:
		return ""
	}
}

// Dereferences pointer or returns the value
func indirectValue(i any) reflect.Value {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v
}

// convert map[int]string → map[string]string (for JSON/cache safety)
func toStringMap(m map[int]string) map[string]string {
	res := make(map[string]string, len(m))
	for id, name := range m {
		res[strconv.Itoa(id)] = name
	}
	return res
}

// convert map[string]string → map[int]string
func fromStringMap(m map[string]string) map[int]string {
	res := make(map[int]string, len(m))
	for k, v := range m {
		if id, err := strconv.Atoi(k); err == nil {
			res[id] = v
		}
	}
	return res
}
