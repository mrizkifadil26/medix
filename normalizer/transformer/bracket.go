package transformer

import (
	"fmt"
	"regexp"
)

func RemoveBrackets(input string) (string, error) {
	pattern := regexp.MustCompile(`[\[\(\{][^\[\]\(\)\{\}]{1,30}[\]\)\}]`)

	if pattern == nil {
		return input, fmt.Errorf("removeBrackets: pattern is nil")
	}

	return pattern.ReplaceAllString(input, ""), nil
}

func init() {
	GetTransformerRegistry().
		Register("removeBrackets", RemoveBrackets)
}
