package tmdb

import (
	"fmt"
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

	var result GenreResult
	if err := c.doRequest(endpoint, nil, &result); err != nil {
		return nil, err
	}

	return result.Genres, nil
}
