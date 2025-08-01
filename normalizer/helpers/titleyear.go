package normalizer

import (
	"regexp"
	"strconv"
	"strings"
)

var yearRegex = regexp.MustCompile(`(?i)(.*?)\s*\((\d{4})\)$`)

// ExtractTitleYear extracts the title and year from a string like "Interstellar (2014)".
// Returns (title, year) or (original input, 0) if no match.
func ExtractTitleYear(name string) (string, int) {
	matches := yearRegex.FindStringSubmatch(name)
	if len(matches) != 3 {
		return strings.TrimSpace(name), 0
	}

	title := strings.TrimSpace(matches[1])
	year, err := strconv.Atoi(matches[2])
	if err != nil {
		return title, 0
	}

	return title, year
}
