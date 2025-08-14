package transformer

import (
	"regexp"
	"strings"
)

var dashRe = regexp.MustCompile(`[._\-–—]{1,}`)

func NormalizeDashes(s string) string {
	s = dashRe.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func SpaceToDash(s string) string {
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	return strings.ReplaceAll(s, " ", "-")
}

func CollapseDashes(s string) string {
	s = strings.ReplaceAll(s, "_", "-")
	s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
