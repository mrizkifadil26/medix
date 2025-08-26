package local

type SubtitleEnricher struct{}

func (s *SubtitleEnricher) Name() string {
	return "subtitle"
}

func (s *SubtitleEnricher) Enrich(data any) error {
	return nil
}
