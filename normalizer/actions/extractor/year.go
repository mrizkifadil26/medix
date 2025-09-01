package extractor

import (
	"fmt"
	"regexp"
)

func ExtractYear(input string) (string, error) {
	// Regex pattern for years between 1900-2100
	yearPattern := `(19[1-9]\d|20\d{2}|2100)`

	// 1. Try with brackets first
	reBrackets := regexp.MustCompile(`\(` + yearPattern + `\)`)
	if match := reBrackets.FindStringSubmatch(input); len(match) == 2 {
		return match[1], nil
	}

	// 2. Fallback: four consecutive digits within range
	reDigits := regexp.MustCompile(`\b` + yearPattern + `\b`)
	if match := reDigits.FindStringSubmatch(input); len(match) == 2 {
		return match[1], nil
	}

	return "", fmt.Errorf("year not found in input (1900-2100): %q", input)
}

func init() {
	GetRegistry().
		Register("year", ExtractYear)
}
