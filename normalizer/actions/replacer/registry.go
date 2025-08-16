package replacer

import (
	"fmt"

	"github.com/mrizkifadil26/medix/normalizer/registries"
	"github.com/mrizkifadil26/medix/utils"
)

type Replacer func(string, map[string]string) (string, error)
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
	name, input string, params map[string]string,
) (string, error) {
	fn, ok := r.Get(name)
	if !ok {
		return input, fmt.Errorf("replacer %q not found", name)
	}

	return fn(input, params)
}

func init() {
	registries.GetRegistry().
		Register("replace", GetRegistry())
}
