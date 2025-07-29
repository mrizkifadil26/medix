package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// SearchMovie is a helper that searches for movies by title and release year.
func (c *Client) SearchMovie(title string, year int) ([]SearchItem, error) {
	return c.Search("movie", SearchQuery{
		Query:       title,
		PrimaryYear: strconv.Itoa(year),
	})
}

func (c *Client) Search(mediaType string, q SearchQuery) ([]SearchItem, error) {
	// Validate query
	if err := q.Validate(); err != nil {
		return nil, fmt.Errorf("invalid search query: %w", err)
	}

	endpoint := fmt.Sprintf("%s/search/%s", c.BaseURL, mediaType)

	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	reqURL.RawQuery = q.ToParams().Encode()

	// url := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		// return nil, fmt.Errorf("request failed: %w", err)
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
		// return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return result.Results, nil
}
