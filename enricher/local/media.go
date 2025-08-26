package local

type MediaSourceEnricher struct{}

func (s *MediaSourceEnricher) Name() string {
	return "subtitle"
}

func (s *MediaSourceEnricher) Enrich(data any) error {
	return nil
}
