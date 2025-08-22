package traverse

import (
	"fmt"
	"strconv"

	"github.com/mrizkifadil26/medix/utils"
)

// Traversal emits every JSON node as a hit
// type Hit struct {
// 	Path   []string // e.g., ["items", "3", "movieName"]
// 	Key    string   // map key (string) or slice index (int)
// 	Value  any
// 	Parent any // immediate container: map[string]any or []any
// }

type trieNode struct {
	children map[string]*trieNode
	terminal bool
}

type Engine struct {
	root any
}

func NewRoot(root any) *Engine {
	return &Engine{
		root: root,
	}
}

func (e *Engine) GetRoot() any {
	return e.root
}

func (e *Engine) Get(path string) (any, error) {
	selector := CompileSelector(path)

	return getPath(e.root, selector)
}

// func (e *Engine) GetAll(paths []string) (any, error) {
// 	set := NewSelectorSet(paths)
// 	results := make(map[string]any)

// 	return results, nil
// }

func (e *Engine) Set(path string, value any) error {
	selector := CompileSelector(path)

	return setPath(e.root, selector, value)
}

func getPath(node any, sel Selector) (any, error) {
	return descend(node, sel, 0)
}

func descend(node any, sel Selector, pos int) (any, error) {
	if pos >= len(sel.Tokens) {
		return node, nil
	}

	token := sel.Tokens[pos]

	// check if this position is a wildcard (#)
	if isHash(sel, pos) {
		// must be an array
		arr, ok := node.([]any)
		if !ok {
			return nil, fmt.Errorf("cannot use '#' at position %d on type %T", pos, node)
		}

		results := make([]any, 0, len(arr))
		for _, elem := range arr {
			val, err := descend(elem, sel, pos+1)
			if err == nil {
				results = append(results, val)
			}
		}
		if len(results) == 1 {
			return results[0], nil
		}
		return results, nil
	}

	// normal key / index lookup
	switch n := node.(type) {
	case *utils.OrderedMap[string, any]:
		val, ok := n.Get(token)
		if !ok {
			return nil, fmt.Errorf("key %q not found", token)
		}
		return descend(val, sel, pos+1)

	case map[string]any:
		val, ok := n[token]
		if !ok {
			return nil, fmt.Errorf("key %q not found", token)
		}
		return descend(val, sel, pos+1)

	case []any:
		idx, err := strconv.Atoi(token)
		if err != nil || idx < 0 || idx >= len(n) {
			return nil, fmt.Errorf("invalid index %q", token)
		}
		return descend(n[idx], sel, pos+1)

	default:
		return nil, fmt.Errorf("cannot descend into type %T with %q", node, token)
	}
}

func setPath(node any, sel Selector, value any) error {
	return setDescend(node, sel, 0, value)
}

func setDescend(node any, sel Selector, pos int, value any) error {
	if pos >= len(sel.Tokens) {
		return fmt.Errorf("empty path")
	}

	token := sel.Tokens[pos]
	rest := sel.Tokens[pos+1:]

	// wildcard (#) must be applied on []any
	if isHash(sel, pos) {
		fmt.Println("is it hash?")
		arr, ok := node.([]any)
		if !ok {
			return fmt.Errorf("cannot use '#' at position %d on type %T", pos, node)
		}

		for i := range arr {
			if len(rest) == 0 {
				arr[i] = value
			} else {
				if err := setDescend(arr[i], sel, pos+1, value); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// normal key / index
	switch n := node.(type) {
	case *utils.OrderedMap[string, any]:
		val, ok := n.Get(token)
		if !ok {
			// create map if missing
			if len(rest) == 0 {
				val = value
			} else {
				val = utils.NewOrderedMap[string, any]()
			}

			n.Set(token, val)
		} else if len(rest) == 0 {
			n.Set(token, value)
			return nil
		}

		if len(rest) > 0 {
			return setDescend(val, sel, pos+1, value)
		}

		return nil

	case map[string]any:
		val, ok := n[token]
		if !ok {
			if len(rest) == 0 {
				val = value
			} else {
				val = make(map[string]any)
			}

			n[token] = val
		} else if len(rest) == 0 {
			n[token] = value
			return nil
		}

		if len(rest) > 0 {
			return setDescend(val, sel, pos+1, value)
		}

		return nil

	case []any:
		idx, err := strconv.Atoi(token)
		if err != nil {
			return fmt.Errorf("invalid index %q", token)
		}

		// grow slice if needed
		for idx >= len(n) {
			// n = append(n, utils.NewOrderedMap[string, any]())
			n = append(n, nil)
		}

		elem := n[idx]
		if elem == nil {
			if len(rest) == 0 {
				// leaf, will set value later
				elem = nil
			} else {
				// decide type based on parent
				switch node.(type) {
				case *utils.OrderedMap[string, any]:
					elem = utils.NewOrderedMap[string, any]()
				case map[string]any:
					elem = make(map[string]any)
				default:
					return fmt.Errorf("cannot determine type to create for slice element at index %d", idx)
				}
			}

			n[idx] = elem
		}

		if len(rest) == 0 {
			n[idx] = value
			return nil
		}

		switch elem.(type) {
		case *utils.OrderedMap[string, any], map[string]any:
			return setDescend(elem, sel, pos+1, value)
		default:
			return fmt.Errorf("cannot descend into type %T at index %d", elem, idx)
		}

	default:
		return fmt.Errorf("cannot descend into type %T with %q", node, token)
	}
}

func isHash(sel Selector, pos int) bool {
	for _, h := range sel.HashIndex {
		if h == pos {
			return true
		}
	}

	return false
}

func expandHashPaths(node any, sel Selector, pos int, basePath []string) [][]string {
	if pos >= len(sel.Tokens) {
		return [][]string{basePath}
	}

	token := sel.Tokens[pos]

	if isHash(sel, pos) {
		arr, ok := node.([]any)
		if !ok {
			// # on non-array -> error
			return nil
		}

		var allPaths [][]string
		for i := range arr {
			newBase := append(basePath, strconv.Itoa(i))
			paths := expandHashPaths(arr[i], sel, pos+1, newBase)
			allPaths = append(allPaths, paths...)
		}
		return allPaths
	}

	// normal token
	switch n := node.(type) {
	case *utils.OrderedMap[string, any]:
		val, ok := n.Get(token)
		if !ok {
			return nil
		}
		return expandHashPaths(val, sel, pos+1, append(basePath, token))
	case map[string]any:
		val, ok := n[token]
		if !ok {
			return nil
		}
		return expandHashPaths(val, sel, pos+1, append(basePath, token))
	case []any:
		idx, err := strconv.Atoi(token)
		if err != nil || idx < 0 || idx >= len(n) {
			return nil
		}
		return expandHashPaths(n[idx], sel, pos+1, append(basePath, token))
	default:
		return nil
	}
}
