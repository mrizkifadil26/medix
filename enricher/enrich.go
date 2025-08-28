package enricher

import (
	"github.com/mrizkifadil26/medix/enricher/local"
	"github.com/mrizkifadil26/medix/enricher/tmdb"
)

func Enrich(
	data any,
	config *Config,
) (any, error) {
	enrichers := []Enricher{
		tmdb.TMDbEnricher{},
		local.LocalEnricher{},
	}

	return data, nil
}
