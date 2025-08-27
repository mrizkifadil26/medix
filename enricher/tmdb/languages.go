package tmdb

import (
	"fmt"
)

type TMDBLanguage struct {
	ISO639_1    string `json:"iso_639_1"`
	EnglishName string `json:"english_name"`
	Name        string `json:"name"`
}

func (c *Client) GetLanguages() ([]TMDBLanguage, error) {
	endpoint := fmt.Sprintf("%s/configuration/languages", c.BaseURL)

	var results []TMDBLanguage
	if err := c.doRequest(endpoint, nil, &results); err != nil {
		return nil, err
	}

	return results, nil
}
