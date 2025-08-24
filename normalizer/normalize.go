package normalizer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mrizkifadil26/medix/normalizer/registries"
	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type Normalizer struct {
	Fields  []Field
	Meta    map[string]any
	Targets map[string]any
}

func New(config *Config) *Normalizer {
	return &Normalizer{
		Fields:  config.Fields,
		Meta:    make(map[string]any),
		Targets: make(map[string]any),
	}
}

func originalKey(key string) string { return key + ":original" }
func currentKey(key string) string  { return key + ":current" }

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

func (n *Normalizer) Normalize(data any) (any, error) {
	for _, field := range n.Fields {
		value, err := jsonpath.Get(data, field.Name)
		if err != nil {
			return nil, fmt.Errorf("field %q not found: %v", field.Name, err)
		}

		switch v := value.(type) {
		case []any:
			for i, val := range v {
				key := strings.ReplaceAll(field.Name, "#", strconv.Itoa(i))
				if err := n.processField(field, key, val, &i); err != nil {
					return nil, err
				}
			}

		default:
			if err := n.processField(field, field.Name, v, nil); err != nil {
				return nil, err
			}
		}
	}

	if err := n.applyTargets(data); err != nil {
		return nil, err
	}

	return data, nil
}

// processField handles storing vars and updating targets
func (n *Normalizer) processField(
	field Field,
	key string,
	val any,
	idx *int,
) error {
	// store original once
	if _, ok := n.GetMeta(originalKey(key)); !ok {
		n.SetMeta(originalKey(key), val)
	}

	// run actions
	_, err := n.applyActions(key, field.Actions, idx)
	if err != nil {
		return fmt.Errorf("failed on field %q: %v", key, err)
	}

	return nil
}

func (n *Normalizer) applyActions(key string, actions []Action, idx *int) (any, error) {
	registry := registries.GetRegistry()

	current, _ := n.GetMeta(originalKey(key))
	for _, action := range actions {
		input := current

		result, err := registry.Apply(action.Type, input, action.Params)
		// if err != nil {
		// 	return nil, fmt.Errorf("action %q failed: %v", action.Type, err)
		// }
		if err != nil {
			// just log and continue
			fmt.Printf("action %q failed, skipping: %v\n", action.Type, err)
			continue
		}

		// Only mutate current for Transform or Replacer
		if action.Type == "transform" || action.Type == "replace" {
			current = result
			n.SetMeta(currentKey(key), current)
		}

		// update target immediately if defined
		if action.Target != "" {
			target := action.Target
			if idx != nil {
				target = strings.ReplaceAll(target, "#", strconv.Itoa(*idx))
			}

			n.Targets[target] = result
		}
	}

	return current, nil
}

func (n *Normalizer) applyTargets(data any) error {
	for path, value := range n.Targets {
		if err := jsonpath.Set(data, path, value); err != nil {
			return fmt.Errorf("set %q failed: %v", path, err)
		}
	}

	return nil
}
