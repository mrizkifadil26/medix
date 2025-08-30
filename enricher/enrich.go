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
		apiKey, _ := cfgMap["api_key"].(string)
		return tmdb.NewTMDbEnricher(&tmdb.Config{APIKey: apiKey}), nil

	case "local":
		cfgMap, ok := eCfg.Config.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid local config")
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
