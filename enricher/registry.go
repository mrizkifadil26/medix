package enricher

import (
	"github.com/mrizkifadil26/medix/enricher/core"
	"github.com/mrizkifadil26/medix/utils"
)

var enricherRegistry = utils.NewRegistry[core.Enricher]()

func Register(e core.Enricher) {
	enricherRegistry.Register(e.Name(), e)
}

func Get(name string) (core.Enricher, bool) {
	return enricherRegistry.Get(name)
}

func All() map[string]core.Enricher {
	return enricherRegistry.All()
}
