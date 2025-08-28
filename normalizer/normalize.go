package normalizer

import (
	"fmt"
	"log"
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
func (n *Normalizer) getMeta(key string) (any, bool) {
	if n.Meta == nil {
		return nil, false
	}

	val, ok := n.Meta[key]
	return val, ok
}

// SetMeta sets a value in the Meta map by key
func (n *Normalizer) setMeta(key string, value any) {
	if n.Meta == nil {
		n.Meta = map[string]any{}
	}

	n.Meta[key] = value
}

func (n *Normalizer) ensureOriginal(key string, val any) any {
	if existing, ok := n.getMeta(originalKey(key)); ok {
		return existing
	}

	n.setMeta(originalKey(key), val)
	return val
}

func (n *Normalizer) setCurrent(key string, val any) {
	n.setMeta(currentKey(key), val)
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
				if err := n.processField(data, field, key, val, i); err != nil {
					return nil, err
				}
			}

		default:
			if err := n.processField(data, field, field.Name, v, -1); err != nil {
				return nil, err
			}
		}
	}

	return data, nil
}

// processField handles storing vars and updating targets
func (n *Normalizer) processField(
	data any,
	field Field,
	key string,
	val any,
	idx int,
) error {
	// store original once
	orig := n.ensureOriginal(key, val)

	// run actions
	current := orig
	for _, action := range field.Actions {
		result, err := n.applyAction(data, action, current, idx)
		if err != nil {
			log.Printf("field %q action %q skipped: %v", key, action.Type, err)
			continue
		}

		current = result
	}

	n.setCurrent(key, current)
	return nil
}

func (n *Normalizer) applyAction(
	data any,
	action Action,
	input any,
	idx int,
) (any, error) {
	registry := registries.GetRegistry()

	result, err := registry.Apply(action.Type, input, action.Params)
	if err != nil {
		return nil, err
	}

	// update target immediately if defined
	if action.Target != "" {
		// omit if result is "empty"
		if result != nil && result != "" {

			target := action.Target
			if idx >= 0 {
				target = strings.ReplaceAll(target, "#", strconv.Itoa(idx))
			}

			if err := jsonpath.Set(data, target, result); err != nil {
				return nil, fmt.Errorf("set %q failed: %v", target, err)
			}

			n.Targets[target] = result
		}

		// Only mutate current for Transform or Replacer
		if action.Type == "transform" || action.Type == "replace" {
			return result, nil
		}
	}

	return input, nil
}
