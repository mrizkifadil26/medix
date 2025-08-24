package jsonpath

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mrizkifadil26/medix/utils"
)

// Get retrieves a value by path
func Get(root any, path string) (any, error) {
	node := root
	tokens := strings.Split(path, ".")

	for i, token := range tokens {
		isLast := i == len(tokens)-1

		if token == "#" {
			arr, ok := node.([]any)
			if !ok {
				return nil, fmt.Errorf("cannot use '#' at %d on type %T", i, node)
			}
			results := make([]any, 0, len(arr))
			for _, elem := range arr {
				val, err := Get(elem, strings.Join(tokens[i+1:], "."))
				if err == nil {
					results = append(results, val)
				}
			}
			if len(results) == 1 {
				return results[0], nil
			}
			return results, nil
		}

		switch n := node.(type) {
		case map[string]any:
			val, ok := n[token]
			if !ok {
				return nil, fmt.Errorf("key %q not found", token)
			}
			node = val
		case *utils.OrderedMap[string, any]:
			val, ok := n.Get(token)
			if !ok {
				return nil, fmt.Errorf("key %q not found", token)
			}
			node = val
		case []any:
			idx, err := strconv.Atoi(token)
			if err != nil || idx < 0 || idx >= len(n) {
				return nil, fmt.Errorf("invalid index %q", token)
			}
			node = n[idx]
		default:
			if !isLast {
				return nil, fmt.Errorf("cannot descend into type %T with %q", node, token)
			}
		}
	}

	return node, nil
}

// Set assigns a value by path
func Set(root any, path string, value any) error {
	node := root
	tokens := strings.Split(path, ".")

	for i, token := range tokens {
		isLast := i == len(tokens)-1

		switch n := node.(type) {
		case map[string]any:
			if isLast {
				n[token] = value
				return nil
			}
			if _, ok := n[token]; !ok {
				n[token] = make(map[string]any)
			}
			node = n[token]

		case *utils.OrderedMap[string, any]:
			if isLast {
				n.Set(token, value)
				return nil
			}
			if val, ok := n.Get(token); ok {
				node = val
			} else {
				newMap := utils.NewOrderedMap[string, any]()
				n.Set(token, newMap)
				node = newMap
			}

		case []any:
			idx, err := strconv.Atoi(token)
			if err != nil || idx < 0 {
				return fmt.Errorf("invalid index %q", token)
			}
			for idx >= len(n) {
				n = append(n, nil)
			}
			if isLast {
				n[idx] = value
				return nil
			}
			if n[idx] == nil {
				n[idx] = make(map[string]any)
			}
			node = n[idx]

		default:
			return fmt.Errorf("cannot descend into type %T with %q", node, token)
		}
	}

	return nil
}
