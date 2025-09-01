package extractor

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	reBracketYear = regexp.MustCompile(`^\((\d{4})\)$`) // (YYYY)
	reYear        = regexp.MustCompile(`^\d{4}$`)       // plain year
	reAKA         = regexp.MustCompile(`(?i)^AKA$`)     // case-insensitive
)

type Titles struct {
	Main string
	AKA  string // optional, omit if not present
	Year string // optional
}

func ExtractTitles(input string) (Titles, error) {
	var titles Titles

	words := strings.Fields(strings.TrimSpace(input))
	akaIndex, yearIndex := -1, -1

	for i, w := range words {
		if akaIndex == -1 && reAKA.MatchString(w) {
			akaIndex = i
		}

		if match := reBracketYear.FindStringSubmatch(w); len(match) == 2 {
			yearIndex = i
			titles.Year = match[1]
		} else if reYear.MatchString(w) && yearIndex == -1 {
			yearIndex = i
			titles.Year = w
		}
	}

	// Main title: from start to AKA (exclusive) or year (exclusive)
	if akaIndex != -1 {
		titles.Main = strings.Join(words[:akaIndex], " ")
	} else if yearIndex != -1 {
		titles.Main = strings.Join(words[:yearIndex], " ")
	} else {
		titles.Main = strings.Join(words, " ")
	}

	// Alternate title (AKA): optional, only if AKA exists
	if akaIndex != -1 && yearIndex != -1 && akaIndex < yearIndex {
		titles.AKA = strings.Join(words[akaIndex+1:yearIndex], " ")
	} else if akaIndex != -1 {
		titles.AKA = strings.Join(words[akaIndex+1:], " ")
	} else {
		titles.AKA = "" // omit if not present
	}

	if titles.Main == "" {
		return titles, fmt.Errorf("main title could not be extracted from input=%q", input)
	}

	return titles, nil
}

func init() {
	GetRegistry().
		Register("title", func(input string) (string, error) {
			t, err := ExtractTitles(input)
			if err != nil {
				return "", err
			}

			return t.Main, nil
		})

	GetRegistry().
		Register("alternateTitle", func(input string) (string, error) {
			t, err := ExtractTitles(input)
			if err != nil {
				return "", err
			}

			return t.AKA, nil // may be empty
		})
}
