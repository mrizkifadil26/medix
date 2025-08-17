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
	input string, params map[string]any,
) (string, error) {
	pattern, ok := params["pattern"].(string)
	if !ok || pattern == "" {
		return input, fmt.Errorf("pattern not provided")
	}

	fn, ok := r.Get(pattern)
	if !ok {
		return input, fmt.Errorf("extractor %q not found", pattern)
	}

	return fn(input)
}

func init() {
	registries.GetRegistry().
		Register("extract", GetRegistry())
}
