package replacer

import (
	"fmt"
	"strings"
)

func DefaultReplacer(input string, params map[string]any) (string, error) {
	from, ok := params["from"].(string)
	if !ok {
		return "", fmt.Errorf("replace: missing 'from' parameter")
	}

	if from == "" {
		return input, fmt.Errorf("replace: 'from' cannot be empty")
	}

	to := params["to"].(string) // `to` can be empty, which is valid (to remove `from`)
	return strings.ReplaceAll(input, from, to), nil
}

func init() {
	GetRegistry().
		Register("default", DefaultReplacer)
}
