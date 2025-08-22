package extractor

import (
	"errors"
	"regexp"
	"strings"
)

func ExtractTitle(input string) (string, error) {
	re := regexp.MustCompile(`^(?:\d+\.)?(.*?)\.?\(?\d{4}\)?`)
	matches := re.FindStringSubmatch(input)

	if len(matches) > 1 && matches[1] != "" {
		title := strings.ReplaceAll(matches[1], ".", " ")
		return strings.TrimSpace(title), nil
	}

	return "", errors.New("title could not be extracted: no year found or invalid format")
}

func init() {
	GetRegistry().
		Register("title", ExtractTitle)
}
