package enricher

type Enricher interface {
	Name() string
	Enrich(entry any) error
}
