package tmdb

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/mrizkifadil26/medix/utils/validation"
)

type SearchResult struct {
	Results []SearchItem `json:"results"`
}

type SearchItem struct {
	ID               int     `json:"id"`                // Unique TMDb ID
	Title            string  `json:"title"`             // From movie.title or tv.name
	OriginalTitle    string  `json:"original_title"`    // For fallback or search matching
	Overview         string  `json:"overview"`          // Summary
	PosterPath       string  `json:"poster_path"`       // For image thumbnail
	BackdropPath     string  `json:"backdrop_path"`     // For UI background
	ReleaseDate      string  `json:"release_date"`      // movie: release_date / tv: first_air_date
	GenreIDs         []int   `json:"genre_ids"`         // To map to genre names locally
	VoteAverage      float64 `json:"vote_average"`      // For scoring
	VoteCount        int     `json:"vote_count"`        // For credibility
	OriginalLanguage string  `json:"original_language"` // Optional filter
	Popularity       float64 `json:"popularity"`
}

func (m SearchItem) GetTitle() string {
	return m.Title
}

func (m SearchItem) GetOriginalTitle() string {
	return m.OriginalTitle
}

func (m SearchItem) GetReleaseYear() int {
	if m.ReleaseDate == "" {
		return 0
	}
	yearStr := strings.Split(m.ReleaseDate, "-")[0]
	if y, err := strconv.Atoi(yearStr); err == nil {
		return y
	}
	return 0
}

func (m SearchItem) GetPopularity() float64 {
	return m.Popularity
}

func (m SearchItem) GetVoteAverage() float64 {
	return m.VoteAverage
}

func (m SearchItem) GetGenreIDs() []int {
	return m.GenreIDs
}

type SearchQuery struct {
	Query       string `param:"query" validate:"required"`
	Language    string `param:"language"`
	Region      string `param:"region"`
	Year        string `param:"year"`
	PrimaryYear string `param:"primary_release_year"`
	Page        int    `param:"page"`
}

// ToParams converts struct to URL values, skipping empty
func (q SearchQuery) Params() url.Values {
	return buildParams(q)
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

// SearchMovie is a helper that searches for movies by title and release year.
func (c *Client) SearchMovie(title, year string) ([]SearchItem, error) {
	return c.Search("movie", SearchQuery{
		Query:       title,
		PrimaryYear: year,
	})
}

func (c *Client) Search(mediaType string, q SearchQuery) ([]SearchItem, error) {
	// Validate query
	if err := q.Validate(); err != nil {
		return nil, fmt.Errorf("invalid search query: %w", err)
	}

	endpoint := fmt.Sprintf("%s/search/%s", c.BaseURL, mediaType)
	params := q.Params()

	var result SearchResult
	if err := c.doRequest(endpoint, params, &result); err != nil {
		return nil, err
	}

	return result.Results, nil
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
