package formatter

import (
	"fmt"

	"github.com/mrizkifadil26/medix/utils"
)

type Formatter func(string, map[string]string) (string, error)
type FormatterRegistry struct {
	*utils.Registry[Formatter]
}

var formatterSingleton *FormatterRegistry

func GetFormatterRegistry() *FormatterRegistry {
	if formatterSingleton == nil {
		formatterSingleton = &FormatterRegistry{
			Registry: utils.NewRegistry[Formatter](),
		}
	}

	return formatterSingleton
}

// ApplyByName applies a transformer by name to a value
func (r *FormatterRegistry) Apply(
	name, template string, params map[string]string
) (string, error) {
	fn, ok := r.Get(name)
	if !ok {
		return input, fmt.Errorf("formatter %q not found", name)
	}

	return fn(template, params)
}
