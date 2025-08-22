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

func (r *Registry) Apply(
	input any, params map[string]any,
) (string, error) {
	// Convert input to string at the boundary
	strInput, ok := input.(string)
	if !ok {
		return "", fmt.Errorf("expected string input, got %T", input)
	}

	// Get methods
	methodsVal, ok := params["methods"]
	if !ok {
		return strInput, fmt.Errorf("methods not provided")
	}

	var methods []string
	switch v := methodsVal.(type) {
	case string:
		methods = []string{v}
	case []any:
		for _, item := range v {
			s, ok := item.(string)
			if !ok {
				return strInput, fmt.Errorf("methods contains non-string value: %v", item)
			}

			methods = append(methods, s)
		}
	default:
		return strInput, fmt.Errorf("invalid methods type: %T", methodsVal)
	}

	return r.applyAll(strInput, methods)
}

func (r *Registry) applyAll(
	input string,
	methods []string,
) (string, error) {
	for _, name := range methods {
		// fmt.Println("transapply", name)
		fn, ok := r.Get(name)
		if !ok {
			return "", fmt.Errorf("transformer %q not found", name)
		}

		var err error
		input, err = fn(input)
		if err != nil {
			return input, fmt.Errorf("error applying transformer %q: %w", name, err)
		}
	}

	return input, nil
}

func init() {
	registries.GetRegistry().
		Register("transform", GetRegistry())
}
