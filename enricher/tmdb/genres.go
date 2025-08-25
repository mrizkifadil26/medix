package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type GenreResult struct {
	Genres []GenreItem `json:"genres"`
}

type GenreItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Client) GetGenres(mediaType string) ([]GenreItem, error) {
	endpoint := fmt.Sprintf("%s/genre/%s/list", c.BaseURL, mediaType)

	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	// url := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	var result GenreResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return result.Genres, nil
}
