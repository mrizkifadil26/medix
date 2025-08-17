package formatter

import (
	"fmt"
	"strings"
)

func DefaultFormatter(template string, from map[string]any) (string, error) {
	result := template

	for key, val := range from {
		strVal := fmt.Sprintf("%v", val)
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, strVal)
	}

	return result, nil
}

func init() {
	GetRegistry().
		Register("default", DefaultFormatter)
}
