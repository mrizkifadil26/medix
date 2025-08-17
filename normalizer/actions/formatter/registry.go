package formatter

import (
	"fmt"

	"github.com/mrizkifadil26/medix/normalizer/registries"
	"github.com/mrizkifadil26/medix/utils"
)

type Formatter func(string, map[string]any) (string, error)
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
	input string, params map[string]any,
) (string, error) {
	template, ok := params["template"].(string)
	if !ok || template == "" {
		return input, fmt.Errorf("template not provided")
	}

	name := "default"
	fn, ok := r.Get(name)
	if !ok {
		return template, fmt.Errorf("formatter %q not found", name)
	}

	return fn(template, params)
}

func init() {
	registries.GetRegistry().
		Register("format", GetRegistry())
}
