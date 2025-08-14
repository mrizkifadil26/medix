package normalizer

import (
	"github.com/mrizkifadil26/medix/normalizer/extractor"
	"github.com/mrizkifadil26/medix/normalizer/formatter"
	"github.com/mrizkifadil26/medix/normalizer/replacer"
	"github.com/mrizkifadil26/medix/normalizer/transformer"
)

//		"unicodeFix":          normalizer.NormalizeUnicode,
//		"stripExtension":      normalizer.StripExtension,
//		"stripBrackets":       normalizer.StripBrackets,
//		"replaceSpecialChars": normalizer.ReplaceSpecialChars,
//		"collapseDashes":      normalizer.CollapseDashes,
//		"toLower":             normalizer.ToLower,
//		"spaceToDash":         normalizer.SpaceToDash,
//		"removeKnownPrefixes": normalizer.RemoveKnownPrefixes,
//		"dotToSpace":          normalizer.DotToSpace,
//		"slugify":             normalizer.Slugify,
//		"normalizeDashes":     normalizer.NormalizeDashes,
//	}
// var DefaultTransformers = transformers.NewTransformerRegistry(
// 	transformers.UnicodeNormalizer(),
// 	transformers.SanitizeSymbols(),
// 	transformers.NormalizeSeparators(),
// 	transformers.Trim(),
// 	transformers.Slugify(),
// 	transformers.RemoveBrackets(),
// )

type ActionRegistry struct {
	Transformers *transformer.Registry
	Extractors   *extractor.Registry
	Replacers    *replacer.Registry
	Formatters   *formatter.Registry
}

var singleton *ActionRegistry

func GetActions() *ActionRegistry {
	if singleton == nil {
		singleton = &ActionRegistry{
			Transformers: transformer.GetRegistry(),
			Extractors:   extractor.GetRegistry(),
			Replacers:    replacer.GetRegistry(),
			Formatters:   formatter.GetRegistry(),
		}
	}

	return singleton
}
