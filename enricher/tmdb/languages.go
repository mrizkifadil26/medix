package tmdb

import (
	"fmt"

	"github.com/mrizkifadil26/medix/utils/cache"
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

// LoadLanguageMap loads language map from cache, fallback to TMDb if missing
func LoadLanguageMap(client *Client, cm *cache.Manager) (map[string]string, error) {
	// try from cache
	if raw, ok := cm.Get("languages", "iso-639"); ok {
		if sm, ok := raw.(map[string]string); ok {
			return sm, nil
		}
	}

	// fetch from TMDb
	langs, err := client.GetLanguages()
	if err != nil {
		return nil, err
	}

	langMap := toLanguageMap(langs)

	// save to cache
	cm.Put("languages", "iso-639", langMap)
	_ = cm.Save() // best-effort

	return langMap, nil
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
