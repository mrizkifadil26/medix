package transformer

import "strings"

func SanitizeSymbols(input string) (string, error) {
	r := strings.NewReplacer(
		"&", "and",
		"+", "",
		"'", "",
		",", "",
		"!", "",
		"?", "",
		"\"", "",
		"/", "-",
		"\\", "-",
	)

	return r.Replace(input), nil
}

func init() {
	GetRegistry().
		Register("sanitize", SanitizeSymbols)
}
