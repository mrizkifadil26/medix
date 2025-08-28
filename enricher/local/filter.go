package local

import "github.com/mrizkifadil26/medix/utils/datawrapper"

type Filter interface {
	Name() string
	Apply(item datawrapper.Data, errs *[]error)
}
