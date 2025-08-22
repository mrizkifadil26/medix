package formatter

import (
	"fmt"

	"github.com/mrizkifadil26/medix/normalizer/registries"
	"github.com/mrizkifadil26/medix/utils"
)

type Formatter func(any, string) (string, error)
type Registry struct {
	*utils.Registry[Formatter]
}

var singleton *Registry

func GetRegistry() *Registry {
	if singleton == nil {
		singleton = &Registry{
			Registry: utils.NewRegistry[Formatter](),
		}
	}

	return singleton
}

// ApplyByName applies a transformer by name to a value
func (r *Registry) Apply(
	input any, params map[string]any,
) (string, error) {
	template, ok := params["template"].(string)
	if !ok || template == "" {
		return "", fmt.Errorf("input not provided")
	}

	name := "default"
	fn, ok := r.Get(name)
	if !ok {
		return "", fmt.Errorf("formatter %q not found", name)
	}

	return fn(input, template)
}

func init() {
	registries.GetRegistry().
		Register("format", GetRegistry())
}
