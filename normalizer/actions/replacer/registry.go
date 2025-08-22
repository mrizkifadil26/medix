package replacer

import (
	"fmt"

	"github.com/mrizkifadil26/medix/normalizer/registries"
	"github.com/mrizkifadil26/medix/utils"
)

type Replacer func(string, map[string]any) (string, error)
type Registry struct {
	*utils.Registry[Replacer]
}

var singleton *Registry

func GetRegistry() *Registry {
	if singleton == nil {
		singleton = &Registry{
			Registry: utils.NewRegistry[Replacer](),
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

	name := "default"
	fn, ok := r.Get(name)
	if !ok {
		return strInput, fmt.Errorf("replacer %q not found", name)
	}

	return fn(strInput, params)
}

func init() {
	registries.GetRegistry().
		Register("replace", GetRegistry())
}
