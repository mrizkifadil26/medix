package normalizer

import (
	"fmt"
	"strings"

	normalizer "github.com/mrizkifadil26/medix/normalizer/helpers"
)

// --- function types ---
type NormalizeFunc func(string) string
type ExtractFunc func(string) (string, error)
type ReplaceFunc func(string, map[string]string) (string, error)
type FormatFunc func(string, map[string]string) (string, error)

var DefaultNormalizers = map[string]NormalizeFunc{
	"unicodeFix":          normalizer.NormalizeUnicode,
	"stripExtension":      normalizer.StripExtension,
	"stripBrackets":       normalizer.StripBrackets,
	"replaceSpecialChars": normalizer.ReplaceSpecialChars,
	"collapseDashes":      normalizer.CollapseDashes,
	"toLower":             normalizer.ToLower,
	"spaceToDash":         normalizer.SpaceToDash,
	"removeKnownPrefixes": normalizer.RemoveKnownPrefixes,
	"dotToSpace":          normalizer.DotToSpace,
	"slugify":             normalizer.Slugify,
	"normalizeDashes":     normalizer.NormalizeDashes,
}

var DefaultExtractors = map[string]ExtractFunc{
	"year": normalizer.ExtractYear,
}

func DefaultReplacer(input string, params map[string]string) (string, error) {
	from, ok := params["from"]
	if !ok {
		return "", fmt.Errorf("replace: missing 'from' parameter")
	}

	to := params["to"] // `to` can be empty, which is valid (to remove `from`)
	return strings.ReplaceAll(input, from, to), nil
}

func DefaultFormatter(template string, from map[string]string) (string, error) {
	result := template

	for key, val := range from {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, val)
	}

	return result, nil
}

type OperatorRegistry struct {
	NormalizeFuncs map[string]NormalizeFunc
	ExtractFuncs   map[string]ExtractFunc
	ReplaceFunc    ReplaceFunc
	FormatFunc     FormatFunc
}

func NewOperators() *OperatorRegistry {
	return &OperatorRegistry{
		NormalizeFuncs: clone(DefaultNormalizers), // map[string]NormalizerFunc
		ExtractFuncs:   clone(DefaultExtractors),  // map[string]ExtractorFunc
		ReplaceFunc:    DefaultReplacer,           // func(string, map[string]string)
		FormatFunc:     DefaultFormatter,          // func(string, map[string]string)
	}
}

func clone[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = v
	}

	return dst
}
