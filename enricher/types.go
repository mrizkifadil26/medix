package enricher

import "github.com/mrizkifadil26/medix/utils/datawrapper"

type Enricher interface {
	Name() string
	Enrich(data datawrapper.Data) (any, error)
}
