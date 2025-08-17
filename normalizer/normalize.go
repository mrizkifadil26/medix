package normalizer

import (
	"fmt"
	"strings"

	"github.com/iancoleman/orderedmap"
	"github.com/mrizkifadil26/medix/normalizer/registries"
	"github.com/mrizkifadil26/medix/normalizer/traverse"
)

type Normalizer struct {
	Original string
	Value    string
	Meta     map[string]string
	Actions  *registries.ActionRegistry
}

type CompiledField struct {
	Selector traverse.Selector
	Actions  []Action
}

func CompileField(
	field Field,
) CompiledField {
	compiledActions := make([]Action, len(field.Actions))
	for i, action := range field.Actions {
		compiledActions[i] = CompileAction(action)
	}

	selectorFunc := traverse.CompileSelector(field.Name)

	return CompiledField{
		Selector: selectorFunc,
		Actions:  compiledActions,
	}
}

func CompileAction(action Action) Action {
	// Precompute token positions of `#` for fast resolution during traversal
	tokens := strings.Split(action.Target, ".")
	hashIdx := []int{}
	for i, tok := range tokens {
		if tok == "#" {
			hashIdx = append(hashIdx, i)
		}
	}

	// action.HashIndex = hashIdx
	return action
}

func New(input string) *Normalizer {
	return &Normalizer{
		Original: input,
		Value:    input,
		Meta:     make(map[string]string),
		// Actions:  GetRegistry(),
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

func (n *Normalizer) Normalize(
	data *orderedmap.OrderedMap,
	config *Config,
) (any, error) {
	registry := registries.GetRegistry()

	compiled := make([]CompiledField, len(config.Fields))
	for i, field := range config.Fields {
		compiled[i] = CompileField(field)
	}

	engine := traverse.Engine{}
	engine.OnHit = func(hit traverse.Hit, ctx *traverse.TraverseContext) any {
		path := hit.Path
		// fmt.Printf("%v %T %v\n", hit.Path, hit.Value, hit.Value)
		// val := fmt.Sprintf("%v", hit.Value)

		for _, field := range compiled {
			matchField := field.Selector.Match(path)
			// fmt.Println("field: ", field, "compiled: ", compiled)
			// fmt.Println()
			if matchField {
				strVal := fmt.Sprint(hit.Value)
				newVal := strVal

				for _, action := range field.Actions {
					result, err := registry.Apply(
						action.Type, newVal, action.Params,
					)

					if err != nil {
						panic(err)
					}

					newVal = result.(string)
				}
			}
		}

		return nil
	}

	engine.Walk(data, 22)

	return data, nil
}

// -------------------
// Helper: set value
// -------------------
// func setValue(
// 	hit traverse.TraversalHit,
// 	target string, value string,
// 	setter *traverse.PathSetter,
// ) {
// 	if target == "" {
// 		// mutate in-place
// 		switch parent := hit.Parent.(type) {
// 		case map[string]any:
// 			parent[hit.Key.(string)] = value
// 		case []any:
// 			parent[hit.Key.(int)] = value
// 		}
// 	} else {
// 		// resolve # in target path
// 		tokens := strings.Split(target, ".")
// 		for i, t := range tokens {
// 			if t == "#" {
// 				tokens[i] = strconv.Itoa(hit.Index)
// 			}
// 		}

// 		fmt.Println(hit)
// 		fmt.Println(tokens)

// 		setter.Set(tokens, value)
// 	}
// }
