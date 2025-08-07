package normalizer

import (
	"errors"
	"fmt"
	"strings"
)

func ResolvePath(json any, path string) ([]any, error) {
	tokens := tokenize(path)
	return walk(json, tokens)
}

func walk(data any, path []string) ([]any, error) {
	if len(path) == 0 {
		// Reached the leaf node
		return []any{data}, nil
	}

	token := path[0]
	rest := path[1:]

	switch d := data.(type) {
	case map[string]any:
		val, ok := d[token]
		if !ok {
			return nil, fmt.Errorf("field %q not found in object", token)
		}
		return walk(val, rest)
	case []any:
		if token == "#" {
			var results []any
			for i, item := range d {
				vals, err := walk(item, rest)
				if err != nil {
					return nil, fmt.Errorf("error in array at index %d: %w", i, err)
				}

				results = append(results, vals...)
			}

			return results, nil
		} else {
			return nil, fmt.Errorf("unexpected token %q for array (expected '#')", token)
		}

	default:
		return nil, errors.New("unexpected structure; cannot continue path")
	}
}

// case "#":
// 	// Expecting an array
// 	arr, ok := current.([]any)
// 	if !ok {
// 		return nil, fmt.Errorf("expected array at '#', got %T", current)
// 	}
// 	var results []any
// 	for _, elem := range arr {
// 		vals, err := traverseRecursive(elem, rest)
// 		if err != nil {
// 			continue // or log
// 		}

// 		results = append(results, vals...)
// 	}
// 	return results, nil

// default:
// 	// Expecting an object with this key
// 	obj, ok := current.(map[string]any)
// 	if !ok {
// 		return nil, fmt.Errorf("expected object at '%s', got %T", seg, current)
// 	}
// 	next, ok := obj[seg]
// 	if !ok {
// 		return nil, fmt.Errorf("key '%s' not found", seg)
// 	}
// 	return traverseRecursive(next, rest)
// }
// }

func tokenize(path string) []string {
	return strings.Split(path, ".")
}
