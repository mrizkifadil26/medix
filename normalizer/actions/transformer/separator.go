package transformer

import "regexp"

var separatorsRe = regexp.MustCompile(`[._\-–—]+`)
var multiSpaceRe = regexp.MustCompile(`\s+`)

// NormalizeSeparators replaces _, ., –, —, - with space and collapses multiple spaces
func NormalizeSeparators(input string) (string, error) {
	s := separatorsRe.ReplaceAllString(input, " ")
	s = multiSpaceRe.ReplaceAllString(s, " ")
	return s, nil
}

func init() {
	GetRegistry().
		Register("separator", NormalizeSeparators)
}
