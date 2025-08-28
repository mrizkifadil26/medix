package enricher

type Enricher interface {
	Name() string
	Enrich(data any, options map[string]string) (any, error)
}
