package normalizer

import (
	"fmt"
	"strings"
)

func Traverse(input any, selector string) ([]any, error) {
	segments := strings.Split(selector, ".")
	return traverseRecursive(input, segments)
}

func traverseRecursive(current any, segments []string) ([]any, error) {
	if len(segments) == 0 {
		// Reached the leaf node
		return []any{current}, nil
	}

	seg := segments[0]
	rest := segments[1:]

	switch seg {
	case "#":
		// Expecting an array
		arr, ok := current.([]any)
		if !ok {
			return nil, fmt.Errorf("expected array at '#', got %T", current)
		}
		var results []any
		for _, elem := range arr {
			vals, err := traverseRecursive(elem, rest)
			if err != nil {
				continue // or log
			}

			results = append(results, vals...)
		}
		return results, nil

	default:
		// Expecting an object with this key
		obj, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected object at '%s', got %T", seg, current)
		}
		next, ok := obj[seg]
		if !ok {
			return nil, fmt.Errorf("key '%s' not found", seg)
		}
		return traverseRecursive(next, rest)
	}
}
