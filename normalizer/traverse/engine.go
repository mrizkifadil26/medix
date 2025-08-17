package traverse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/iancoleman/orderedmap"
)

// Traversal emits every JSON node as a hit
type Hit struct {
	Path   []string // e.g., ["items", "3", "movieName"]
	Value  any
	Parent any // immediate container: map[string]any or []any
	Key    any // map key (string) or slice index (int)
	Index  int // numeric index for arrays (used for #)
}

type Engine struct {
	OnHit      func(hit Hit, ctx *TraverseContext) any // user provides handler
	maxHits    int
	hitCounter int
}

type TraverseContext struct {
	root  any
	cache map[string]any
	batch map[string]any
}

func newContext(root any) *TraverseContext {
	return &TraverseContext{
		root:  root,
		cache: make(map[string]any),
		batch: make(map[string]any),
	}
}

func (ctx *TraverseContext) SetBatch(path, value string) {
	ctx.batch[path] = value
}

func (ctx *TraverseContext) ApplyBatch() {
	for pathStr, value := range ctx.batch {
		fmt.Println(pathStr, value)
		ctx.insertByPath(pathStr, value)
	}

	ctx.batch = make(map[string]any)
}

func (ctx *TraverseContext) insertByPath(pathStr string, value any) {
	if cached, ok := ctx.cache[pathStr]; ok {
		// update existing cached node
		switch node := cached.(type) {
		case *any:
			*node = value
			return
		}
	}

	tokens := splitPath(pathStr)
	current := ctx.root
	var parent any = nil
	var parentKey any = nil

	for i, key := range tokens[:len(tokens)-1] {
		switch node := current.(type) {

		case *orderedmap.OrderedMap:
			if next, ok := node.Get(key); ok {
				parent = node
				parentKey = key
				current = next
			} else {
				var newNode any
				if _, err := strconv.Atoi(tokens[i+1]); err == nil {
					newNode = []any{}
				} else {
					newNode = orderedmap.New()
				}
				node.Set(key, newNode)
				parent = node
				parentKey = key
				current = newNode
			}

		case []any:
			idx, _ := strconv.Atoi(key)
			for len(node) <= idx {
				node = append(node, nil)
			}

			// Attach to parent if slice itself is nested
			if parent != nil {
				switch p := parent.(type) {
				case *orderedmap.OrderedMap:
					p.Set(fmt.Sprint(parentKey), node)
				case []any:
					p[parentKey.(int)] = node
				}
			}

			parent = node
			parentKey = idx

			// Ensure current node
			if node[idx] == nil {
				node[idx] = orderedmap.New()
			}
			current = node[idx]

		default:
			// fallback: treat as new OrderedMap
			newNode := orderedmap.New()

			if parent != nil {
				switch p := parent.(type) {
				case *orderedmap.OrderedMap:
					p.Set(fmt.Sprint(parentKey), newNode)
				case []any:
					p[parentKey.(int)] = newNode
				}
			}

			parent = newNode
			current = newNode
		}
	}

	// Set last token
	last := tokens[len(tokens)-1]
	switch node := current.(type) {
	case *orderedmap.OrderedMap:
		node.Set(last, value)
	case []any:
		idx, _ := strconv.Atoi(last)
		for len(node) <= idx {
			node = append(node, nil)
		}
		node[idx] = value
	default:
		panic(fmt.Sprintf("unsupported type %T at last segment", current))
	}

	// Cache final node
	ctx.cache[pathStr] = value
}

func (ctx *TraverseContext) Root() any {
	return ctx.root
}

func (e *Engine) walkInternal(
	node any,
	path []string,
	parent any,
	key any,
	ctx *TraverseContext,
) {
	if e.hitCounter >= e.maxHits {
		return
	}

	hit := Hit{
		Path:   path,
		Parent: parent,
		Value:  node,
		Key:    "",
		Index:  -1,
	}

	// assign Key/Index properly depending on parent container type
	switch parent.(type) {
	case []any:
		if i, ok := key.(int); ok {
			hit.Index = i
		}

	case *orderedmap.OrderedMap:
		if s, ok := key.(string); ok {
			hit.Key = s
		}

	case map[string]any:
		if s, ok := key.(string); ok {
			hit.Key = s
		}
	}

	// --- emit log ---
	fmt.Printf(
		"Hit path=%v key=%q idx=%d type=%T value=%v\n",
		hit.Path, hit.Key, hit.Index, hit.Value, hit.Value,
	)

	if e.OnHit != nil {
		e.OnHit(hit, ctx)
		e.hitCounter++
	}

	switch n := node.(type) {
	case *orderedmap.OrderedMap:
		for _, k := range n.Keys() {
			child, _ := n.Get(k)
			e.walkInternal(child, append(path, k), n, k, ctx)
		}

	case map[string]any:
		for k, child := range n {
			e.walkInternal(child, append(path, k), n, k, ctx) // parent = map, key = string
		}

	case []any:
		for i, child := range n {
			e.walkInternal(child, append(path, fmt.Sprintf("%d", i)), n, i, ctx)
		}
	}
}

func (e *Engine) Walk(root *orderedmap.OrderedMap, maxHits int) any {
	e.maxHits = maxHits
	e.hitCounter = 0

	ctx := newContext(root)
	e.walkInternal(root, []string{}, nil, nil, ctx)
	ctx.ApplyBatch()

	return ctx.Root()
}

func splitPath(path string) []string {
	// simple split, can extend for selector notation with #
	return strings.Split(path, ".")
}
