package normalizer

import (
	"encoding/json"
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
	current := input

	for key, opVal := range field {
		var err error
		switch key {
		case "replace":
			current, err = r.applyReplace(current, opVal)

		case "normalize":
			current, err = r.applyNormalize(current, opVal)

		case "extract":
			current, err = r.applyExtract(current, opVal)

		case "format":
			current, err = r.applyFormat(current, field, opVal)

		default:
			continue
		}

		if err != nil {
			return "", err
		}
	}

	return strings.TrimSpace(current), nil
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

func (r *OperatorRegistry) applyReplace(input string, opVal any) (string, error) {
	replaceCfg, ok := opVal.(map[string]string)
	if !ok {
		return "", fmt.Errorf("replace: expected map[string]string but got %T", opVal)
	}

	from, _ := replaceCfg["from"]
	to, _ := replaceCfg["to"]
	return r.ReplaceFunc(input, map[string]string{"from": from, "to": to})
}

func (r *OperatorRegistry) applyExtract(input string, opVal any) (string, error) {
	if name, ok := opVal.(string); ok {
		if fn, found := r.ExtractFuncs[name]; found {
			return fn(input)
		}
	}

	return input, nil
}

func (r *OperatorRegistry) applyNormalize(input string, opVal any) (string, error) {
	switch steps := opVal.(type) {
	case string:
		if fn, found := r.NormalizeFuncs[steps]; found {
			return fn(input), nil
		}
	case []any:
		for _, step := range steps {
			if name, ok := step.(string); ok {
				if fn, found := r.NormalizeFuncs[name]; found {
					input = fn(input)
				}
			}
		}
	default:
		return "", fmt.Errorf("normalize: invalid type %T", opVal)
	}

	return input, nil
}

func (r *OperatorRegistry) applyFormat(input string, field map[string]any, opVal any) (string, error) {
	formatStr, ok := opVal.(string)
	if !ok {
		return "", fmt.Errorf("format: expected string format but got %T", opVal)
	}

	fromMap, ok := field["from"].(map[string]string)
	if !ok {
		return "", fmt.Errorf("format: missing or invalid 'from' map")
	}

	return r.FormatFunc(formatStr, fromMap)
}

var ErrUnsupportedInput = fmt.Errorf("unsupported input type: must be string or []string")

func PrettyPrintJSON(field map[string]any) string {
	bytes, err := json.MarshalIndent(field, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(bytes)
}
