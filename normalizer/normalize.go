package normalizer

import (
	"github.com/mrizkifadil26/medix/normalizer/registries"
	"github.com/mrizkifadil26/medix/normalizer/traverse"
)

type compiledField struct {
	Selector     traverse.Selector
	TargetChains map[string][]Action
}

type Normalizer struct {
	Fields []compiledField
}

func New(config *Config) *Normalizer {
	var fields []compiledField
	for _, f := range config.Fields {
		tc := make(map[string][]Action)
		for _, action := range f.Actions {
			tc[action.Target] = append(
				tc[action.Target],
				action,
			)

			fields = append(fields, compiledField{
				Selector:     traverse.CompileSelector(f.Name),
				TargetChains: tc,
			})
		}
	}

	return &Normalizer{Fields: fields}
}

func (n *Normalizer) Normalize(data any) (any, error) {
	registry := registries.GetRegistry()

	// engine := traverse.Engine{}
	// engine.OnHit = func(hit traverse.Hit) any {
	// 	path := hit.Path
	// 	for _, field := range n.Fields {
	// 		if !field.Selector.Match(path) {
	// 			continue
	// 		}

	// 		for _, actions := range field.TargetChains {
	// 			value := fmt.Sprint(hit.Value)

	// 			for _, action := range actions {
	// 				result, err := registry.Apply(action.Type, value, action.Params)
	// 				if err != nil {
	// 					break // skip this chain on error
	// 				}

	// 				value = result.(string)
	// 			}

	// fmt.Printf("%v - %v\n", target, value)
	// 		}
	// 	}

	// 	return nil
	// }

	// err := engine.Walk()
	// if err != nil {
	// 	return nil, err
	// }

	return data, nil
}
