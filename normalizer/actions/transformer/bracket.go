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

func RemoveSquareBrackets(input string) (string, error) {
	pattern := regexp.MustCompile(`\[[^\[\]]*\]`)

	if pattern == nil {
		return input, fmt.Errorf("removeSquareBrackets: pattern is nil")
	}

	return pattern.ReplaceAllString(input, ""), nil
}

func init() {
	GetRegistry().
		Register("removeBrackets", RemoveBrackets)

	GetRegistry().
		Register("removeSquareBrackets", RemoveSquareBrackets)
}
