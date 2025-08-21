package registries

import (
	"fmt"
)

type ActionTypeRegistry interface {
	Apply(input string, params map[string]any) (string, error)
}

type ActionRegistry struct {
	registries map[string]ActionTypeRegistry
}

var singleton *ActionRegistry

func GetRegistry() *ActionRegistry {
	if singleton == nil {
		singleton = &ActionRegistry{
			registries: make(map[string]ActionTypeRegistry),
		}
	}

	return singleton
}

func (r *ActionRegistry) Register(actionType string, registry ActionTypeRegistry) {
	r.registries[actionType] = registry
}

func (r *ActionRegistry) Get(actionType string) (ActionTypeRegistry, bool) {
	fn, ok := r.registries[actionType]
	return fn, ok
}

func (r *ActionRegistry) All() map[string]ActionTypeRegistry {
	return r.registries
}

func (r *ActionRegistry) Apply(
	actionType string,
	input string,
	params map[string]any,
) (any, error) {
	reg, ok := r.Get(actionType)
	if !ok {
		return nil, fmt.Errorf("action type %q not found", actionType)
	}

	return reg.Apply(input, params)
}
