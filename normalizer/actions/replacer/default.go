package replacer

import (
	"fmt"
	"strings"
)

func DefaultReplacer(input string, params map[string]string) (string, error) {
	from, ok := params["from"]
	if !ok {
		return "", fmt.Errorf("replace: missing 'from' parameter")
	}

	if from == "" {
		return input, fmt.Errorf("replace: 'from' cannot be empty")
	}

	to := params["to"] // `to` can be empty, which is valid (to remove `from`)
	return strings.ReplaceAll(input, from, to), nil
}
