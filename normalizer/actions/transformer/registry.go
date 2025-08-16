package transformer

import (
	"fmt"

	"github.com/mrizkifadil26/medix/normalizer/registries"
	"github.com/mrizkifadil26/medix/utils"
)

type Transformer func(string) (string, error)
type Registry struct {
	*utils.Registry[Transformer]
}

var singleton *Registry

func GetRegistry() *Registry {
	if singleton == nil {
		singleton = &Registry{
			Registry: utils.NewRegistry[Transformer](),
		}
	}

	return singleton
}

// ApplyByName applies a transformer by name to a value
func (r *Registry) Apply(
	name, input string,
) (string, error) {
	fn, ok := r.Get(name)
	if !ok {
		return input, fmt.Errorf("transformer %q not found", name)
	}

	return fn(input)
}

// ApplyByName applies a transformer by name to a value
func (r *Registry) ApplyAll(
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

func init() {
	registries.GetRegistry().
		Register("transform", GetRegistry())
}
