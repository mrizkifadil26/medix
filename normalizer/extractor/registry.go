package extractor

import (
	"fmt"

	"github.com/mrizkifadil26/medix/utils"
)

type Extractor func(string) (string, error)
type ExtractorRegistry struct {
	*utils.Registry[Extractor]
}

var extractorSingleton *ExtractorRegistry

func GetExtractorRegistry() *ExtractorRegistry {
	if extractorSingleton == nil {
		extractorSingleton = &ExtractorRegistry{
			Registry: utils.NewRegistry[Extractor](),
		}
	}

	return extractorSingleton
}

// ApplyByName applies a transformer by name to a value
func (r *ExtractorRegistry) Apply(
	name, input string,
) (string, error) {
	fn, ok := r.Get(name)
	if !ok {
		return input, fmt.Errorf("transformer %q not found", name)
	}

	return fn(input)
}

// ApplyByName applies a transformer by name to a value
func (r *ExtractorRegistry) ApplyAll(
	names []string,
	input string,
) (string, error) {
	var err error
	for _, name := range names {
		fn, ok := r.Get(name)
		if !ok {
			return "", fmt.Errorf("transformer %q not found", name)
		}

		input, err = fn(input)
		if err != nil {
			return "", nil
		}
	}

	return input, nil
}
