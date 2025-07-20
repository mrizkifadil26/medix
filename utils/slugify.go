// utils/slugify.go
//
// Package utils provides general-purpose utility functions.
//
// This file implements a slugification helper for transforming strings
// (e.g. file names or titles) into safe, URL/filename-friendly slugs.
package utils

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Slugify converts an input string (e.g. file name or title) into a slug-friendly format.
// It performs the following transformations:
//   - Converts to lowercase
//   - Removes file extensions
//   - Replaces special characters (& → "and", accented letters → plain)
//   - Collapses multiple hyphens
//   - Strips diacritics (é → e)
//   - Replaces spaces with dashes
//
// Example:
//
//	Slugify("Café – Hello World (2023).mp4") → "cafe-hello-world-2023"
func Slugify(name string) string {
	name = strings.TrimSuffix(name, filepath.Ext(name))
	name = strings.ToLower(name)

	// Replace encoded Unicode ampersand
	name = strings.ReplaceAll(name, `\u0026`, "&")

	// Normalize smart quotes and other characters
	name = strings.NewReplacer(
		"’", "'", // smart apostrophe
		"‘", "'",
		"“", `"`,
		"”", `"`,
		"³", "3",
		"½", "1-2",
	).Replace(name)

	// Strip diacritics like é → e
	name = removeDiacritics(name)

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

// removeDiacritics removes accent marks from characters (e.g., é → e).
func removeDiacritics(s string) string {
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)
	result, _, _ := transform.String(t, s)
	return result
}
