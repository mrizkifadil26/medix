package local

type IconEnricher struct{}

func (s *IconEnricher) Name() string {
	return "subtitle"
}

func (s *IconEnricher) Enrich(data any) error {
	return nil
}
