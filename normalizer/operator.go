package normalizer

import (
	"fmt"
	"strings"

	normalizer "github.com/mrizkifadil26/medix/normalizer/helpers"
)

type NormalizeFunc func(string) string

var Normalizers = map[string]NormalizeFunc{
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

type ExtractFunc func(string) (string, error)

var Extractors = map[string]ExtractFunc{
	"year": normalizer.ExtractYear,
}

type ReplaceFunc func(string, map[string]string) (string, error)

func DefaultReplacer(input string, params map[string]string) (string, error) {
	from, ok := params["from"]
	if !ok {
		return "", fmt.Errorf("replace: missing 'from' parameter")
	}

	to := params["to"] // `to` can be empty, which is valid (to remove `from`)
	return strings.ReplaceAll(input, from, to), nil
}

type FormatFunc func(string, map[string]string) (string, error)

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
		NormalizeFuncs: Normalizers,      // map[string]NormalizerFunc
		ExtractFuncs:   Extractors,       // map[string]ExtractorFunc
		ReplaceFunc:    DefaultReplacer,  // func(string, map[string]string)
		FormatFunc:     DefaultFormatter, // func(string, map[string]string)
	}
}
