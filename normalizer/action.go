package normalizer

import (
	"github.com/mrizkifadil26/medix/normalizer/extractor"
	"github.com/mrizkifadil26/medix/normalizer/formatter"
	"github.com/mrizkifadil26/medix/normalizer/replacer"
	"github.com/mrizkifadil26/medix/normalizer/transformer"
)

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
