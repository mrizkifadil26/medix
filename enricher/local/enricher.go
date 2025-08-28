package local

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type LocalEnricher struct{}

func (s *LocalEnricher) Name() string {
	return "local"
}

var filters = []Filter{
	MediaFilter{},
	SubtitlesFilter{},
	IconFilter{},
	CollectionFilter{},
}

func (e *LocalEnricher) Enrich(
	data any,
	options map[string]string,
) (any, error) {
	var errs []error

	// Get optional filter argument from options
	filterArg := ""
	if options != nil {
		filterArg = options["filters"]
	}

	// Build allowed map
	allowed := map[string]bool{}
	if filterArg == "" {
		for _, f := range filters {
			allowed[f.Name()] = true
		}
	} else {
		for _, f := range filters {
			allowed[f.Name()] = false
		}
		for _, name := range strings.Split(filterArg, ",") {
			allowed[strings.TrimSpace(name)] = true
		}
	}

	// Get items using your Get helper
	nodes, err := jsonpath.Get(data, "items.#")
	if err != nil {
		// items key not found, nothing to enrich
		return data, nil
	}

	items, ok := nodes.([]any)
	if !ok {
		return data, fmt.Errorf("items is not an array, got %T", nodes)
	}

	for _, item := range items {
		for _, f := range filters {
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
