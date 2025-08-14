package formatter

import (
	"fmt"
	"strings"
)

func DefaultFormatter(template string, from map[string]string) (string, error) {
	result := template

	for key, val := range from {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, val)
	}

	return result, nil
}

func init() {
	GetRegistry().
		Register("default", DefaultFormatter)
}
