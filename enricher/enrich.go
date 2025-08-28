package enricher

import (
	"fmt"

	"github.com/mrizkifadil26/medix/enricher/local"
	"github.com/mrizkifadil26/medix/enricher/tmdb"
)

func Enrich(
	data any,
	config *Config,
) (any, error) {
	enrichers := []Enricher{
		&tmdb.TMDbEnricher{},
		&local.LocalEnricher{},
	}

	params := map[string]string{
		"concurrency": fmt.Sprint(config.Concurrency), // example
		// add other config → params mapping here
	}

	var err error
	for _, enricher := range enrichers {
		data, err = enricher.Enrich(data, params)
		if err != nil {
			return nil, fmt.Errorf("%s enricher failed: %w", enricher.Name(), err)
		}
	}

	return data, nil
}
