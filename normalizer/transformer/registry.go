package transformer

import (
	"fmt"

	"github.com/mrizkifadil26/medix/utils"
)

type TransformerRegistry struct {
	*utils.Registry[Transformer]
}

var transformerSingleton *TransformerRegistry

func GetTransformerRegistry() *TransformerRegistry {
	if transformerSingleton == nil {
		transformerSingleton = &TransformerRegistry{
			Registry: utils.NewRegistry[Transformer](),
		}
	}

	return transformerSingleton
}

// ApplyByName applies a transformer by name to a value
func (r *TransformerRegistry) Apply(
	name, input string,
) (string, error) {
	fn, ok := r.Get(name)
	if !ok {
		return input, fmt.Errorf("transformer %q not found", name)
	}

	return fn(input)
}

// ApplyByName applies a transformer by name to a value
func (r *TransformerRegistry) ApplyAll(
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
