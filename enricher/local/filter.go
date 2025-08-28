package local

type Filter interface {
	Name() string
	Apply(item any, errs *[]error)
}
