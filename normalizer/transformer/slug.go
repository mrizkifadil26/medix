package transformer

import (
	"regexp"
	"strings"
)

func Slugify(input string) (string, error) {
	re := regexp.MustCompile(`[^a-z0-9]+`)

	s := strings.ToLower(input)
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s, nil
}

func init() {
	GetRegistry().
		Register("slugify", Slugify)
}
