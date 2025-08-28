package tmdb

import (
	"fmt"

	"github.com/mrizkifadil26/medix/utils/cache"
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

// GetGenreMap fetches genres directly from TMDb and returns them as map[id]name.
func (c *Client) GetGenreMap(mediaType string) (map[int]string, error) {
	genres, err := c.GetGenres(mediaType)
	if err != nil {
		return nil, err
	}

	return toGenreMap(genres), nil
}

func LoadGenreMap(client *Client, cm *cache.Manager, kind string) (map[int]string, error) {
	// try from cache
	if raw, ok := cm.Get("genres", kind); ok {
		if m, ok := raw.(map[int]string); ok {
			return m, nil
		}

		// if it's stored as map[string]string, convert
		if sm, ok := raw.(map[string]string); ok {
			converted := make(map[int]string, len(sm))
			for k, v := range sm {
				var id int
				fmt.Sscanf(k, "%d", &id)
				converted[id] = v
			}

			return converted, nil
		}
	}

	// fallback: fetch remote
	resultGenres, err := client.GetGenres(kind)
	if err != nil {
		return nil, err
	}

	genreMap := make(map[int]string, len(resultGenres))
	for _, g := range resultGenres {
		genreMap[g.ID] = g.Name
	}

	// store in cache manager (lazy save later)
	cm.Put("genres", kind, genreMap)
	_ = cm.Save() // save best-effort

	return genreMap, nil
}

// build map[int]string from []GenreItem
func toGenreMap(items []GenreItem) map[int]string {
	m := make(map[int]string, len(items))
	for _, g := range items {
		m[g.ID] = g.Name
	}

	return m
}
