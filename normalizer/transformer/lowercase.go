package transformer

import "strings"

func Lowercase(input string) (string, error) {
	return strings.ToLower(input), nil
}

func init() {
	GetRegistry().
		Register("lowercase", Lowercase)
}
