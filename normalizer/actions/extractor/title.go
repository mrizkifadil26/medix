package extractor

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	reBracketYear = regexp.MustCompile(`^(.*?)\s*\((\d{4})\)`)
	rePlainYear   = regexp.MustCompile(`^(.*?)(?:\.?\d{4})`)
)

func ExtractTitle(input string) (string, error) {
	// 1. Check bracketed year first
	if match := reBracketYear.FindStringSubmatch(input); len(match) == 3 {
		title := strings.ReplaceAll(match[1], ".", " ")
		return strings.TrimSpace(title), nil
	}

	// 2. Fallback to plain year (e.g. 28.Years.Later.2025)
	if match := rePlainYear.FindStringSubmatch(input); len(match) == 2 {
		title := strings.ReplaceAll(match[1], ".", " ")
		return strings.TrimSpace(title), nil
	}

	return "", fmt.Errorf("title could not be extracted: no year found or input format is invalid, input=%q", input)
}

func init() {
	GetRegistry().
		Register("title", ExtractTitle)
}
