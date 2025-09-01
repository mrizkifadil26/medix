package validation_test

import (
	"errors"
	"testing"

	"github.com/mrizkifadil26/medix/utils/validation"
)

type SearchQuery struct {
	Query string `validate:"required"`
	Year  string `validate:"required"`
	Page  int    `validate:"min=1,max=100"`
}

func TestValidateSuccess(t *testing.T) {
	q := SearchQuery{
		Query: "Matrix",
		Year:  "1999",
		Page:  2,
	}

	err := validation.Validate(q, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMissingQuery(t *testing.T) {
	q := SearchQuery{
		Year: "1999",
		Page: 1,
	}

	err := validation.Validate(q, nil)
	if err == nil || err.Error() != "field Query is required" {
		t.Errorf("expected missing Query error, got: %v", err)
	}
}

func TestMinMax(t *testing.T) {
	q := SearchQuery{
		Query: "Test",
		Year:  "2000",
		Page:  0,
	}

	err := validation.Validate(q, nil)
	if err == nil || err.Error() != "field Page must be >= 1" {
		t.Errorf("expected min error on Page, got: %v", err)
	}
}

func TestCustomValidator(t *testing.T) {
	q := SearchQuery{
		Query: "Test",
		Year:  "20",
		Page:  5,
	}

	custom := map[string]validation.ValidationFunc{
		"Year": func(v any) error {
			year := v.(string)
			if len(year) != 4 {
				return errors.New("must be 4-digit")
			}

			return nil
		},
	}

	err := validation.Validate(q, custom)
	if err == nil || err.Error() != "field Year: must be 4-digit" {
		t.Errorf("expected custom error, got: %v", err)
	}
}
