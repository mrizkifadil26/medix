package extractor

import (
	"fmt"
	"regexp"
)

var reBrackets = regexp.MustCompile(`\((\d{4})\)`)

// ExtractBracketYear extracts a 4-digit year only if it's inside brackets, e.g. "Movie (2019)".
func ExtractBracketYear(input string) (string, error) {
	if match := reBrackets.FindStringSubmatch(input); len(match) == 2 {
		return match[1], nil
	}
	return "", fmt.Errorf("bracketed year not found in input: %q", input)
}

func init() {
	GetRegistry().
		Register("bracketYear", ExtractBracketYear)
}
