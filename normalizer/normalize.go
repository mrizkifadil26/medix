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

func New(input string, operators *OperatorRegistry) *Normalizer {
	if operators == nil {
		operators = NewOperators()
	}

	return &Normalizer{
		Original:  input,
		Value:     input,
		Meta:      make(map[string]string),
		Operators: operators,
	}
}

func (n *Normalizer) SetMeta(key, val string) {
	if n.Meta == nil {
		n.Meta = make(map[string]string)
	}

	n.Meta[key] = val
}

func (n *Normalizer) GetMeta(key string) (string, bool) {
	val, ok := n.Meta[key]
	return val, ok
}

var ErrUnsupportedInput = fmt.Errorf("unsupported input type: must be string or []string")

// func (n *Normalizer) Normalize(s string, steps []string) string {
// 	for _, step := range steps {
// 		if fn, ok := Normalizers[step]; ok {
// 			s = fn(s)
// 		}
// 	}

// 	return strings.TrimSpace(s) // <- always trim after normalization
// }

// func (n *Normalizer) Run(input any, steps []string) (any, error) {
// 	switch v := input.(type) {
// 	case string:
// 		return n.Normalize(v, steps), nil

// 	case []any:
// 		var result []string
// 		for _, val := range v {
// 			if s, ok := val.(string); ok {
// 				result = append(result, n.Normalize(s, steps))
// 			}
// 		}
// 		return result, nil

// 	default:
// 		return nil, ErrUnsupportedInput
// 	}
// }

// ApplyOperators applies operations defined in the 'field' map in a deterministic order:
// replace -> normalize -> extract -> format
//
// Expected shapes (common):
//   - replace: map[string]any  { "from": "a", "to": "b" }
//   - normalize: string | []any (string names of normalizers)
//   - extract: string (extractor name)
//   - format: string (template), and field["from"] -> map[string]any for placeholders
func (r *OperatorRegistry) ApplyOperators(input string, field map[string]any) (string, error) {
	current := input
	var err error

	// deterministic order
	if val, ok := field["replace"]; ok {
		current, err = r.applyReplace(current, val)
		if err != nil {
			return "", err
		}
	}

	if val, ok := field["normalize"]; ok {
		current, err = r.applyNormalize(current, val)
		if err != nil {
			return "", err
		}
	}

	if val, ok := field["extract"]; ok {
		current, err = r.applyExtract(current, val)
		if err != nil {
			return "", err
		}
	}

	if val, ok := field["format"]; ok {
		current, err = r.applyFormat(current, field, val)
		if err != nil {
			return "", err
		}
	}

	return strings.TrimSpace(current), nil
}

// --- apply helpers ---
func (r *OperatorRegistry) applyReplace(input string, opVal any) (string, error) {
	// support map[string]string or map[string]any (common when unmarshalling JSON)
	if m, ok := toStringMap(opVal); ok {
		from := m["from"] // may be empty string
		to := m["to"]
		return r.ReplaceFunc(input, map[string]string{"from": from, "to": to})
	}

	return "", fmt.Errorf("replace: expected map[string]string or map[string]any, got %T", opVal)
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
		return "", fmt.Errorf("format: expected template string, got %T", opVal)
	}

	fromAny, ok := field["from"]
	if !ok {
		// support fallback: when source is the current result use {"value": input}
		return r.FormatFunc(formatStr, map[string]string{"value": input})
	}

	fromMap, ok := toStringMap(fromAny)
	if !ok {
		return "", fmt.Errorf("format: expected 'from' map with string values, got %T", fromAny)
	}

	// if template references {{value}} and it's not present, inject current input
	if _, exists := fromMap["value"]; !exists {
		fromMap["value"] = input
	}

	return r.FormatFunc(formatStr, fromMap)
}

// toStringMap tries to convert map[string]any or map[string]string to map[string]string
func toStringMap(v any) (map[string]string, bool) {
	switch m := v.(type) {
	case map[string]string:
		// copy to avoid mutation of caller's map
		out := make(map[string]string, len(m))
		for k, val := range m {
			out[k] = val
		}

		return out, true
	case map[string]any:
		out := make(map[string]string, len(m))
		for k, val := range m {
			if s, ok := val.(string); ok {
				out[k] = s
			} else {
				// best-effort convert non-strings to string using fmt
				out[k] = fmt.Sprintf("%v", val)
			}
		}

		return out, true
	default:
		return nil, false
	}
}

func PrettyPrintJSON(field map[string]any) string {
	bytes, err := json.MarshalIndent(field, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(bytes)
}
