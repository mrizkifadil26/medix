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
}

func New(config *Config) *Normalizer {
	return &Normalizer{
		Fields: config.Fields,
	}
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
				strVal, ok := val.(string)
				if !ok {
					strVal = "" // fallback for non-string values
				}

				for _, action := range field.Actions {
					result, err := registry.Apply(
						action.Type, strVal, action.Params)

					if err != nil {
						if action.Type == "extract" {
							continue
						}

						return nil, fmt.Errorf("error in transforming arr: %v", err)
					}

					resultStr, ok := result.(string)
					if !ok {
						resultStr = ""
					}

					target := strings.ReplaceAll(action.Target, "#", strconv.Itoa(i))
					if err := engine.Set(target, resultStr); err != nil {
						return nil, fmt.Errorf("error while saving")
					}
				}
			}

		default:
			for _, action := range field.Actions {
				newVal := v.(string)
				result, err := registry.Apply(
					action.Type, newVal, action.Params)

				if err != nil {
					return nil, fmt.Errorf("error in transforming prim: %v", err)
				}

				result = result.(string)
				if err := engine.Set(action.Target, result); err != nil {
					return nil, fmt.Errorf("error while saving")
				}
			}
		}
	}

	return engine.GetRoot(), nil
}
