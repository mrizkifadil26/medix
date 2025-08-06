package normalizer

import (
	"errors"
	"regexp"
)

func ExtractYear(input string) (string, error) {
	re := regexp.MustCompile(`\((\d{4})\)`)
	match := re.FindStringSubmatch(input)

	if len(match) != 2 {
		return "", errors.New("year not found in format (YYYY)")
	}

	return match[1], nil
}
