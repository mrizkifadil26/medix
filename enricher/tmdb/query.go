package tmdb

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"

	"github.com/mrizkifadil26/medix/validation"
)

type SearchQuery struct {
	Query       string `param:"query" validate:"required"`
	Language    string `param:"language"`
	Region      string `param:"region"`
	Year        string `param:"year"`
	PrimaryYear string `param:"primary_release_year"`
	Page        int    `param:"page"`
}

// ToParams converts struct to URL values, skipping empty
func (q SearchQuery) ToParams() url.Values {
	values := url.Values{}
	v := reflect.ValueOf(q)
	t := reflect.TypeOf(q)

	for i := 0; i < v.NumField(); i++ {
		val := v.Field(i)
		field := t.Field(i)
		key := field.Tag.Get("param")

		if key == "" {
			continue
		}

		switch val.Kind() {
		case reflect.String:
			if val.String() != "" {
				values.Set(key, val.String())
			}

		case reflect.Int, reflect.Int64:
			if val.Int() != 0 {
				values.Set(key, strconv.FormatInt(val.Int(), 10))
			}

		}
	}

	return values
}

// Validate uses global validator package
func (q SearchQuery) Validate() error {
	return validation.Validate(q, map[string]validation.ValidationFunc{
		"Year": func(v any) error {
			s, ok := v.(string)
			if !ok || s == "" {
				return nil
			}

			if len(s) != 4 || !isNumeric(s) {
				return errors.New("must be 4-digit number")
			}

			return nil
		},
		"PrimaryYear": func(v any) error {
			s, ok := v.(string)
			if !ok || s == "" {
				return nil
			}

			if len(s) != 4 || !isNumeric(s) {
				return errors.New("must be 4-digit number")
			}

			return nil
		},
	})
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isValidYear(s string) bool {
	if len(s) != 4 {
		return false
	}

	_, err := strconv.Atoi(s)
	return err == nil
}
