package enricher

import (
	"fmt"

	"github.com/mrizkifadil26/medix/enricher/core"
	"github.com/mrizkifadil26/medix/enricher/local"
	"github.com/mrizkifadil26/medix/enricher/tmdb"
)

func Enrich(
	data any,
	config *Config,
) (any, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	// Build enrichers from registry
	var enrichers []core.Enricher
	for _, eCfg := range config.Options.Enrichers {
		enricher, err := buildEnricher(eCfg)
		if err != nil {
			return nil, err
		}

		enrichers = append(enrichers, enricher)
	}

	params := map[string]string{
		"concurrency": fmt.Sprint(config.Options.Concurrency),
	}

	// Run each enricher sequentially
	var err error
	for _, enricher := range enrichers {
		data, err = enricher.Enrich(data, params)
		if err != nil {
			return nil, fmt.Errorf("%s enricher failed: %w", enricher.Name(), err)
		}
	}

	return data, nil
}

func buildEnricher(eCfg EnricherConfig) (core.Enricher, error) {
	switch eCfg.Name {
	case "tmdb":
		cfgMap, ok := eCfg.Config.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid tmdb config")
		}

		cfg := &tmdb.Config{}
		if apiKey, ok := cfgMap["api_key"].(string); ok {
			cfg.APIKey = apiKey
		}

		// Optional fetch flags
		if fetchCredits, ok := cfgMap["fetch_credits"].(bool); ok {
			cfg.FetchCredits = fetchCredits
		}

		return tmdb.NewTMDbEnricher(cfg), nil

	case "local":
		cfgMap, ok := eCfg.Config.(map[string]interface{})
		if !ok {
			// no config provided, just use default (all filters)
			return local.NewLocalEnricher(&local.Config{}), nil
		}

		var cfg local.Config
		if f, exists := cfgMap["filters"]; exists {
			for _, v := range f.([]any) {
				if s, ok := v.(string); ok {
					cfg.Filters = append(cfg.Filters, s)
				}
			}
		}
		return local.NewLocalEnricher(&cfg), nil

	default:
		return nil, fmt.Errorf("unknown enricher: %s", eCfg.Name)
	}
}
