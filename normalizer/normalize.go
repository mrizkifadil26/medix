package normalizer

import (
	"fmt"
	"strings"

	normalizer "github.com/mrizkifadil26/medix/normalizer/helpers"
)

type Normalizer struct{}

func New() *Normalizer {
	return &Normalizer{}
}

type NormalizerFunc func(string) string

var Normalizers = map[string]NormalizerFunc{
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
