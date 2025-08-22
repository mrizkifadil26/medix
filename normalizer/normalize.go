package normalizer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mrizkifadil26/medix/normalizer/registries"
	"github.com/mrizkifadil26/medix/normalizer/traverse"
)

type Normalizer struct {
	Fields []Field
	Meta   map[string]any
}

func New(config *Config) *Normalizer {
	return &Normalizer{
		Fields: config.Fields,
	}
}

// GetMeta retrieves a value from the Meta map by key
func (n *Normalizer) GetMeta(key string) (any, bool) {
	if n.Meta == nil {
		return nil, false
	}

	val, ok := n.Meta[key]
	return val, ok
}

// SetMeta sets a value in the Meta map by key
func (n *Normalizer) SetMeta(key string, value any) {
	if n.Meta == nil {
		n.Meta = map[string]any{}
	}

	n.Meta[key] = value
}

func (n *Normalizer) setIntermediate(key string, value any) {
	if n.Meta == nil {
		n.Meta = map[string]any{}
	}
	interm, ok := n.Meta["intermediate"].(map[string]any)
	if !ok || interm == nil {
		interm = map[string]any{}
		n.Meta["intermediate"] = interm
	}
	interm[key] = value
}

func (n *Normalizer) getIntermediate(key string) (any, bool) {
	if n.Meta == nil {
		return nil, false
	}
	interm, ok := n.Meta["intermediate"].(map[string]any)
	if !ok || interm == nil {
		return nil, false
	}
	val, exists := interm[key]
	return val, exists
}

func (n *Normalizer) Normalize(data any) (any, error) {
	registry := registries.GetRegistry()

	engine := traverse.NewRoot(data)
	for _, field := range n.Fields {
		value, err := engine.Get(field.Name)
		if err != nil {
			return nil, fmt.Errorf("value not found for field %q", field.Name)
		}

		switch v := value.(type) {
		case []any:
			for i, val := range v {
				original := val

				for _, action := range field.Actions {
					if action.Type == "transform" {
						key := field.Name
						if strings.Contains(field.Name, "#") {
							key = strings.ReplaceAll(field.Name, "#", strconv.Itoa(i))
						}
						if cached, ok := n.getIntermediate(key); ok {
							val = cached
						}

						result, err := registry.Apply(action.Type, val, action.Params)
						if err != nil {
							return nil, fmt.Errorf("error transforming array: %v", err)
						}
						val = result

						n.setIntermediate(key, val)
					} else {
						val = original

						// Non-transform actions just apply normally
						result, err := registry.Apply(action.Type, val, action.Params)
						if err != nil {
							return nil, fmt.Errorf("error applying action %q: %v", action.Type, err)
						}
						val = result
					}

					// Skip set if target is empty
					if action.Target != "" {
						target := strings.ReplaceAll(action.Target, "#", strconv.Itoa(i))
						if err := engine.Set(target, val); err != nil {
							return nil, fmt.Errorf("error setting value for field %q: %v", field.Name, err)
						}
					}
				}
			}

		default:
			key := field.Name
			for _, action := range field.Actions {
				// Use cached intermediate
				if cached, ok := n.getIntermediate(key); ok {
					v = cached
				}

				result, err := registry.Apply(
					action.Type, v, action.Params)

				if err != nil {
					return nil, fmt.Errorf("error in transforming prim: %v", err)
				}
				v = result

				n.setIntermediate(key, v)

				// Skip set if target is empty
				if action.Target != "" {
					if err := engine.Set(action.Target, v); err != nil {
						return nil, fmt.Errorf("error setting value for field %q: %v", field.Name, err)
					}
				}
			}
		}
	}

	return engine.GetRoot(), nil
}
