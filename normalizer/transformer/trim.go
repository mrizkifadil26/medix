package transformer

import "strings"

// Trim removes leading/trailing whitespace
func Trim(input string) (string, error) {
	return strings.TrimSpace(input), nil
}

func init() {
	GetRegistry().
		Register("trim", Trim)
}
