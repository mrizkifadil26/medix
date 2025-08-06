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
	var err error

	// var appliedKeys = make(map[string]bool) // track what keys we actually apply
	fmt.Println(PrettyPrintJSON(field))
	// var operatorOrder = []string{"replace", "normalize", "extract", "format"}
	for key, opVal := range field {
		// fmt.Println("Key: ", key)
		switch key {
		case "replace":
			// fmt.Println("Type of: ", reflect.TypeOf(field))
			// fmt.Println(PrettyPrintJSON(field))
			// fmt.Println(field["replace"])
			replaceCfg, ok := opVal.(map[string]string)
			if !ok {
				return "", fmt.Errorf("invalid type for 'replace': expected map[string]string")
			}

			// fmt.Println(PrettyPrintJSON())
			// if replaceCfg, ok := field["replace"].(map[string]any); ok {
			from, _ := replaceCfg["from"]
			to, _ := replaceCfg["to"]

			fmt.Println("REPLACE:", from, "=>", to, "| before:", current)
			current, err = r.ReplaceFunc(current, map[string]string{"from": from, "to": to})
			if err != nil {
				fmt.Println("REPLACE ERROR:", err)
				return "", err
			}

			fmt.Println("REPLACE after:", current)
			// appliedKeys["replace"] = true
			// }

		case "normalize":
			if steps, ok := field["normalize"]; ok {
				switch s := steps.(type) {
				case string:
					if fn, found := r.NormalizeFuncs[s]; found {
						current = fn(current)
						// appliedKeys["normalize"] = true

					}
				case []any:
					for _, step := range s {
						if name, ok := step.(string); ok {
							if fn, found := r.NormalizeFuncs[name]; found {
								current = fn(current)
								// appliedKeys["normalize"] = true
							}
						}
					}
				default:
					return "", fmt.Errorf("normalize: invalid type %T", steps)
				}
			}

		case "extract":
			if extractorName, ok := field["extract"].(string); ok {
				if fn, found := r.ExtractFuncs[extractorName]; found {
					current, err = fn(current)
					if err != nil {
						return "", err
					}

					// appliedKeys["extract"] = true
				}
			}

		case "format":
			formatStr, ok := opVal.(string)
			if !ok {
				return "", fmt.Errorf("format: expected string format but got %T", opVal)
			}

			fromMap, ok := field["from"].(map[string]string)
			if !ok {
				return "", fmt.Errorf("format: missing or invalid 'from' map")
			}

			current, err = r.FormatFunc(formatStr, fromMap)
			if err != nil {
				return "", fmt.Errorf("format failed: %w", err)
			}

			// if formatStr, ok := field["format"].(string); ok {
			// 	if fromMap, ok := field["from"].(map[string]string); ok {
			// 		current, err = r.FormatFunc(formatStr, fromMap)
			// 		if err != nil {
			// 			return "", err
			// 		}

			// 		// appliedKeys["format"] = true
			// 	}
			// }
		}
	}

	// check if all keys in `field` were successfully applied
	// for key := range field {
	// 	if !appliedKeys[key] {
	// 		return "", fmt.Errorf("field %q was not successfully applied", key)
	// 	}
	// }

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

var ErrUnsupportedInput = fmt.Errorf("unsupported input type: must be string or []string")

func PrettyPrintJSON(field map[string]any) string {
	bytes, err := json.MarshalIndent(field, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(bytes)
}
