package extractor

import (
	"fmt"

	"github.com/mrizkifadil26/medix/normalizer/registries"
	"github.com/mrizkifadil26/medix/utils"
)

type Extractor func(string) (string, error)
type Registry struct {
	*utils.Registry[Extractor]
}

var singleton *Registry

func GetRegistry() *Registry {
	if singleton == nil {
		singleton = &Registry{
			Registry: utils.NewRegistry[Extractor](),
		}
	}

	return singleton
}

// ApplyByName applies a transformer by name to a value
func (r *Registry) Apply(
	input any, params map[string]any,
) (string, error) {
	// Convert input to string at the boundary
	strInput, ok := input.(string)
	if !ok {
		return "", fmt.Errorf("expected string input, got %T", input)
	}

	pattern, ok := params["pattern"].(string)
	if !ok || pattern == "" {
		return strInput, fmt.Errorf("pattern not provided")
	}

	fn, ok := r.Get(pattern)
	if !ok {
		return strInput, fmt.Errorf("extractor %q not found", pattern)
	}

	return fn(strInput)
}

func init() {
	registries.GetRegistry().
		Register("extract", GetRegistry())
}
