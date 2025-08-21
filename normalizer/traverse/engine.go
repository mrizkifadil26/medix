package traverse

// Traversal emits every JSON node as a hit
type Hit struct {
	// Path   []string // e.g., ["items", "3", "movieName"]
	// Key    string   // map key (string) or slice index (int)
	// Value  any
	// Parent any // immediate container: map[string]any or []any
}

type Engine struct {
	Root any
	// OnHit func(hit Hit) any
}

// func (e *Engine) walk(
// 	node any,
// 	parent any,
// 	key string,
// 	path []string,
// 	cb func(hit Hit) any,
// ) any {
// 	if node == nil && parent != nil {
// 		// If parent exists, create a map by default
// 		node = utils.NewOrderedMap[string, any]()
// 	}

// 	switch n := node.(type) {
// 	case *utils.OrderedMap[string, any]:
// 		for _, key := range n.Keys() {
// 			value, _ := n.Get(key)
// 			newPath := append(path, key)

// 			newVal := e.walk(value, n, key, newPath, cb)
// 			if newVal != nil {
// 				n.Set(key, newVal)
// 			}
// 		}

// 	case []any:
// 		for i, elem := range n {
// 			idx := strconv.Itoa(i)
// 			newPath := append(path, idx)

// 			newVal := e.walk(elem, key, idx, newPath, cb)
// 			if newVal != nil {
// 				n[i] = newVal
// 			}
// 		}
// 	}

// 	// call cb for current node
// 	return cb(Hit{
// 		Path:   path,
// 		Key:    key,
// 		Value:  node,
// 		Parent: parent,
// 	})
// }

// func (e *Engine) Walk() error {
// 	if e.OnHit == nil {
// 		return e.walk(e.Root, nil, "", nil, func(hit Hit) any { return nil })
// 	}

// 	return e.walk(e.Root, nil, "", nil, e.OnHit)
// }

func (e *Engine) Get(path string) (any, error) {
	selector := CompileSelector(path)
	return result, nil
}

func (e *Engine) Set(path string, value any) error {
	selector := CompileSelector(path)
	return nil
}
