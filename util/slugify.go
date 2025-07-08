package util

import (
	"path/filepath"
	"regexp"
	"strings"
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
