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

type GenreService struct {
	client *Client
	cache  *genreCache
}

func NewGenreService(client *Client, cache *genreCache) *GenreService {
	return &GenreService{client: client, cache: cache}
}

func (s *GenreService) Get(kind string) (map[int]string, error) {
	// Try cache first
	if m, ok := s.cache.Get("genres", kind); ok {
		return m, nil
	}

	// Remote fetch
	resultGenres, err := s.client.GetGenres(kind)
	if err != nil {
		return nil, err
	}

	genreMap := toGenreMap(resultGenres)

	// Store into cache
	s.cache.Put("genres", kind, genreMap)
	_ = s.cache.Save() // best-effort save

	return genreMap, nil
}

func (c *Client) GetGenres(mediaType string) ([]GenreItem, error) {
	endpoint := fmt.Sprintf("%s/genre/%s/list", c.BaseURL, mediaType)

	var result GenreResult
	if err := c.doRequest(endpoint, nil, &result); err != nil {
		return nil, err
	}

	return result.Genres, nil
}

func (s *GenreService) Resolve(kind string, ids []int) ([]string, error) {
	genreMap, err := s.Get(kind)
	if err != nil {
		return nil, err
	}

	var genres []string
	for _, id := range ids {
		if g, ok := genreMap[id]; ok {
			genres = append(genres, g)
		}
	}

	return genres, nil
}

// build map[int]string from []GenreItem
func toGenreMap(items []GenreItem) map[int]string {
	m := make(map[int]string, len(items))
	for _, g := range items {
		m[g.ID] = g.Name
	}

	return m
}
