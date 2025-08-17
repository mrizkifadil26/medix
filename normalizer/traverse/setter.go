package traverse

import (
	"fmt"
	"strconv"
)

type PathSetter struct {
	Data any
}

func (p *PathSetter) Set(path []string, value any) {
	current := p.Data
	for i, key := range path[:len(path)-1] {
		switch node := current.(type) {
		case map[string]any:
			if next, ok := node[key]; ok {
				current = next
			} else {
				// decide map vs slice based on next token
				nextKey := path[i+1]
				var newNode any
				if _, err := strconv.Atoi(nextKey); err == nil {
					newNode = []any{}
				} else {
					newNode = map[string]any{}
				}
				node[key] = newNode
				current = newNode
			}
		case []any:
			idx, _ := strconv.Atoi(key)
			for len(node) <= idx {
				node = append(node, map[string]any{})
			}
			current = node[idx]
		default:
			panic(fmt.Sprintf("unsupported type %T at path segment %s", current, key))
		}
	}

	last := path[len(path)-1]
	switch node := current.(type) {
	case map[string]any:
		node[last] = value
	case []any:
		idx, _ := strconv.Atoi(last)
		for len(node) <= idx {
			node = append(node, nil)
		}
		node[idx] = value
	default:
		panic(fmt.Sprintf("unsupported type %T at last segment", current))
	}
}
