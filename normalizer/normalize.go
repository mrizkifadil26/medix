package normalizer

import (
	"fmt"
	"strings"
)

type Normalizer struct{}

func New() *Normalizer {
	return &Normalizer{}
}

type NormalizerFunc func(string) string

var Normalizers = map[string]NormalizerFunc{
	"unicodeFix":          NormalizeUnicode,
	"stripExtension":      StripExtension,
	"stripBrackets":       StripBrackets,
	"replaceSpecialChars": ReplaceSpecialChars,
	"collapseDashes":      CollapseDashes,
	"toLower":             ToLower,
	"spaceToDash":         SpaceToDash,
	"removeKnownPrefixes": RemoveKnownPrefixes,
	"dotToSpace":          DotToSpace,
	"slugify":             Slugify,
	"normalizeDashes":     NormalizeDashes,
}

func (n *Normalizer) Normalize(s string, steps []string) string {
	for _, step := range steps {
		if fn, ok := Normalizers[step]; ok {
			s = fn(s)
		}
	}

	return strings.TrimSpace(s) // <- always trim after normalization
}

func (n *Normalizer) Run(input any, steps []string) (any, error) {
	switch v := input.(type) {
	case string:
		return n.Normalize(v, steps), nil

	case []any:
		var result []string
		for _, val := range v {
			if s, ok := val.(string); ok {
				result = append(result, n.Normalize(s, steps))
			}
		}
		return result, nil

	default:
		return nil, ErrUnsupportedInput
	}
}

var ErrUnsupportedInput = fmt.Errorf("unsupported input type: must be string or []string")
