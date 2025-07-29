package normalizer

import (
	"regexp"
	"strings"
)

func Slugify(s string) string {
	s = NormalizeUnicode(s)
	s = strings.ToLower(s)
	s = NormalizeDashes(s)
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
