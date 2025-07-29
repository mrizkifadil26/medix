package tmdb_test

import (
	"net/url"
	"testing"

	"github.com/mrizkifadil26/medix/enricher/tmdb"
	"github.com/stretchr/testify/assert"
)

func TestToParams(t *testing.T) {
	query := tmdb.SearchQuery{
		Query:       "Batman",
		Language:    "en-US",
		Region:      "US",
		Year:        "2005",
		PrimaryYear: "2005",
		Page:        2,
	}

	params := query.ToParams()

	expected := url.Values{}
	expected.Set("query", "Batman")
	expected.Set("language", "en-US")
	expected.Set("region", "US")
	expected.Set("year", "2005")
	expected.Set("primary_release_year", "2005")
	expected.Set("page", "2")

	assert.Equal(t, expected, params)
}

func TestToParams_OmitsEmptyFields(t *testing.T) {
	query := tmdb.SearchQuery{
		Query: "Interstellar",
		Page:  1,
	}

	params := query.ToParams()

	assert.Equal(t, "Interstellar", params.Get("query"))
	assert.Equal(t, "1", params.Get("page"))
	assert.Empty(t, params.Get("language"))
	assert.Empty(t, params.Get("region"))
	assert.Empty(t, params.Get("year"))
}

func TestValidate_ValidQueryOnly(t *testing.T) {
	q := tmdb.SearchQuery{Query: "Oppenheimer"}
	err := q.Validate()
	assert.NoError(t, err)
}

func TestValidate_ValidAllFields(t *testing.T) {
	q := tmdb.SearchQuery{
		Query:       "Tenet",
		Year:        "2020",
		PrimaryYear: "2020",
	}
	err := q.Validate()
	assert.NoError(t, err)
}

func TestValidate_MissingQuery(t *testing.T) {
	q := tmdb.SearchQuery{
		Year: "2020",
	}

	err := q.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "field Query is required")
}

func TestValidate_InvalidYearFormat(t *testing.T) {
	q := tmdb.SearchQuery{
		Query: "Matrix",
		Year:  "20A0",
	}
	err := q.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "field Year: must be 4-digit number")
}

func TestValidate_InvalidPrimaryYearFormat(t *testing.T) {
	q := tmdb.SearchQuery{
		Query:       "Matrix",
		PrimaryYear: "abcd",
	}
	err := q.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "field PrimaryYear: must be 4-digit number")
}
