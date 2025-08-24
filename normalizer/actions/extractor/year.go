package extractor

import (
	"fmt"
	"regexp"
)

func ExtractYear(input string) (string, error) {
	// 1. Try with brackets first
	reBrackets := regexp.MustCompile(`\((\d{4})\)`)
	if match := reBrackets.FindStringSubmatch(input); len(match) == 2 {
		return match[1], nil
	}

	// 2. Fallback: four consecutive digits
	reDigits := regexp.MustCompile(`\b(\d{4})\b`)
	if match := reDigits.FindStringSubmatch(input); len(match) == 2 {
		return match[1], nil
	}

	return "", fmt.Errorf("year not found in input: %q", input)
}

func init() {
	GetRegistry().
		Register("year", ExtractYear)
}
