package normalizer

import (
	"fmt"
	"strings"
)

type Normalizer struct {
	Original  string
	Value     string
	Meta      map[string]string
	Operators *OperatorRegistry
}

func (n *Normalizer) SetMeta(key, val string) {
	n.Meta[key] = val
}

func (n *Normalizer) GetMeta(key string) (string, bool) {
	val, ok := n.Meta[key]
	return val, ok
}

func New(input string, operators *OperatorRegistry) *Normalizer {
	return &Normalizer{
		Original:  input,
		Value:     input,
		Meta:      make(map[string]string),
		Operators: operators,
	}
}

func (r *OperatorRegistry) ApplyOperators(input string, field map[string]any) (string, error) {
	value := input
	var err error

	// normalize
	if steps, ok := field["normalize"].([]any); ok {
		for _, step := range steps {
			if name, ok := step.(string); ok {
				if fn, found := r.NormalizeFuncs[name]; found {
					value = fn(value)
				}
			}
		}
	}

	// replace
	if replaceCfg, ok := field["replace"].(map[string]any); ok {
		from, _ := replaceCfg["from"].(string)
		to, _ := replaceCfg["to"].(string)
		value, err = r.ReplaceFunc(value, map[string]string{"from": from, "to": to})
		if err != nil {
			return "", err
		}
	}

	// extract
	if extractorName, ok := field["extract"].(string); ok {
		if fn, found := r.ExtractFuncs[extractorName]; found {
			value, err = fn(value)
			if err != nil {
				return "", err
			}
		}
	}

	// format
	if formatStr, ok := field["format"].(string); ok {
		if fromMap, ok := field["from"].(map[string]string); ok {
			value, err = r.FormatFunc(formatStr, fromMap)
			if err != nil {
				return "", err
			}
		}
	}

	return strings.TrimSpace(value), nil
}

func (n *Normalizer) Normalize(s string, steps []string) string {
	for _, step := range steps {
		if fn, ok := Normalizers[step]; ok {
			s = fn(s)
		}
	}

	return strings.TrimSpace(s) // <- always trim after normalization
}

func (n *Normalizer) Run(input any, steps []string) (any, error) {
	switch v := input.(type) {
	case string:
		return n.Normalize(v, steps), nil

	case []any:
		var result []string
		for _, val := range v {
			if s, ok := val.(string); ok {
				result = append(result, n.Normalize(s, steps))
			}
		}
		return result, nil

	default:
		return nil, ErrUnsupportedInput
	}
}

var ErrUnsupportedInput = fmt.Errorf("unsupported input type: must be string or []string")
