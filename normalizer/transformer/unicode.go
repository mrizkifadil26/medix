package transformer

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// removeDiacritics strips accent marks and other modifiers.
func removeDiacritics(s string) string {
	t := transform.Chain(
		norm.NFD,                           // Decompose into base + diacritic
		runes.Remove(runes.In(unicode.Mn)), // Remove diacritics
		norm.NFC,                           // Recompose
	)
	result, _, _ := transform.String(t, s)
	return result
}

// UnicodeNormalizer fixes fancy quotes, dashes, fractions, and removes accents.
func UnicodeNormalizer(input string) (string, error) {
	replacements := map[string]string{
		// Quotes
		"‘": "'", "’": "'",
		"“": `"`, "”": `"`,
		"«": `"`, "»": `"`,
		"‚": "'", "„": `"`, "‹": "'", "›": "'",

		// Dashes & separators
		"–": "-", // En dash
		"—": "-", // Em dash
		"−": "-", // Minus
		"•": "-", // Bullet
		"·": "-", // Middle dot

		// Fractions
		"½": "1-2", "¼": "1-4", "¾": "3-4",
		"⅓": "1-3", "⅔": "2-3",
		"⅛": "1-8", "⅜": "3-8", "⅝": "5-8", "⅞": "7-8",

		// Superscripts
		"¹": "1", "²": "2", "³": "3",

		// Misc
		"©": "c", "®": "r", "™": "tm",
	}

	s := input
	for old, new := range replacements {
		s = strings.ReplaceAll(s, old, new)
	}

	return removeDiacritics(s), nil
}

func init() {
	GetTransformerRegistry().
		Register("unicode", UnicodeNormalizer)
}
