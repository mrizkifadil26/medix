package util

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Slugify(name string) string {
	name = strings.TrimSuffix(name, filepath.Ext(name))
	name = strings.ToLower(name)

	// Replace Unicode \u0026 with actual character
	name = strings.ReplaceAll(name, `\u0026`, "&")

	name = strings.NewReplacer(
		"’", "'", // smart apostrophe
		"‘", "'",
		"“", `"`,
		"”", `"`,
		"³", "3",
		"½", "1-2",
	).Replace(name)

	name = removeDiacritics(name) // <-- ✅ strip accents like é → e

	// Replace " - " with "-"
	name = strings.ReplaceAll(name, " - ", "-")

	// Clean up special characters
	replacer := strings.NewReplacer(
		"&", "and",
		"+", "",
		"'", "",
		",", "",
		"_", "-",
		"(", "",
		")", "",
		".", "",
		"!", "",
		"?", "",
		`"`, "",
		"/", "-",
		"\\", "-",
	)
	name = replacer.Replace(name)

	// Replace all whitespace with "-"
	name = strings.ReplaceAll(name, " ", "-")

	// Collapse multiple dashes
	name = regexp.MustCompile(`-+`).ReplaceAllString(name, "-")

	// Trim leading/trailing dashes
	name = strings.Trim(name, "-")

	return name
}

func removeDiacritics(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn = nonspacing marks
}
