package normalizer

import (
	"fmt"
)

type Normalizer struct {
	Original string
	Value    string
	Meta     map[string]string
	Actions  *ActionRegistry
}

func New(input string, actions *ActionRegistry) *Normalizer {
	return &Normalizer{
		Original: input,
		Value:    input,
		Meta:     make(map[string]string),
		Actions:  GetActions(),
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
