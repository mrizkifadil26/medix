package normalizer

import "regexp"

var bracketPattern = regexp.MustCompile(`[\[\(\{][^\[\]\(\)\{\}]{1,30}[\]\)\}]`)

func StripBrackets(s string) string {
	return bracketPattern.ReplaceAllString(s, "")
}
