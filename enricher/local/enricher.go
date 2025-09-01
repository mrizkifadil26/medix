package local

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type LocalEnricher struct {
	config *Config
}

var allFilters = []Filter{
	MediaFilter{},
	SubtitlesFilter{},
	IconFilter{},
	CollectionFilter{},
}

func NewLocalEnricher(cfg *Config) *LocalEnricher {
	return &LocalEnricher{config: cfg}
}

func (s *LocalEnricher) Name() string {
	return "local"
}

func (e *LocalEnricher) Enrich(
	data any,
	options map[string]string,
) (any, error) {
	var errs []error

	// Build allowed map
	allowed := map[string]bool{}
	if len(e.config.Filters) == 0 {
		for _, f := range allFilters {
			allowed[f.Name()] = true
		}
	} else {
		for _, f := range allFilters {
			allowed[f.Name()] = false
		}
		for _, name := range e.config.Filters {
			allowed[strings.TrimSpace(name)] = true
		}
	}

	// Get items using your Get helper
	nodes, err := jsonpath.Get(data, "items.#")
	if err != nil {
		return data, nil
	}

	items, ok := nodes.([]any)
	if !ok {
		return data, fmt.Errorf("items is not an array, got %T", nodes)
	}

	for _, item := range items {
		for _, f := range allFilters {
			if allowed[f.Name()] {
				f.Apply(item, &errs)
			}
		}
	}

	if len(errs) > 0 {
		return data, errors.Join(errs...)
	}

	return data, nil
}
