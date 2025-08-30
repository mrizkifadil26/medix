package core

// The shared interface
type Enricher interface {
	Name() string
	Enrich(data any, options map[string]string) (any, error)
}
