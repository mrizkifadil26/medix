package normalizer

import "strings"

func ReplaceSpecialChars(s string) string {
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

	return r.Replace(s)
}
