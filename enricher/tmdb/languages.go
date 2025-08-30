package tmdb

import (
	"fmt"
)

type TMDBLanguage struct {
	ISO639_1    string `json:"iso_639_1"`
	EnglishName string `json:"english_name"`
	Name        string `json:"name"`
}

type LanguageService struct {
	client *Client
	cache  *langCache
}

func NewLanguageService(client *Client, cache *langCache) *LanguageService {
	return &LanguageService{client: client, cache: cache}
}

func (s *LanguageService) Get() (map[string]string, error) {
	// Try cache first
	if m, ok := s.cache.Get("languages", "iso-639"); ok {
		return m, nil
	}

	// Fallback → remote fetch
	result, err := s.client.GetLanguages()
	if err != nil {
		return nil, err
	}

	langMap := toLanguageMap(result)

	// Store into cache
	s.cache.Put("languages", "iso-639", langMap)
	_ = s.cache.Save() // best-effort save

	return langMap, nil
}

// Resolve takes ISO-639-1 code and returns the human-readable name.
func (s *LanguageService) Resolve(code string) (string, error) {
	langMap, err := s.Get()
	if err != nil {
		return "", err
	}

	if name, ok := langMap[code]; ok {
		return name, nil
	}

	return code, nil // fallback: return code if unknown
}

func (c *Client) GetLanguages() ([]TMDBLanguage, error) {
	endpoint := fmt.Sprintf("%s/configuration/languages", c.BaseURL)

	var results []TMDBLanguage
	if err := c.doRequest(endpoint, nil, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// build map[code]name from []TMDBLanguage
func toLanguageMap(items []TMDBLanguage) map[string]string {
	m := make(map[string]string, len(items))
	for _, l := range items {
		// prefer English name if available, else fallback
		if l.EnglishName != "" {
			m[l.ISO639_1] = l.EnglishName
		} else {
			m[l.ISO639_1] = l.Name
		}
	}
	return m
}
