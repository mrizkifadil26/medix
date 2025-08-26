package tmdb

type TMDbEnricher struct {
	Client *Client
}

func (t *TMDbEnricher) Name() string { return "tmdb" }

func (t *TMDbEnricher) Enrich(entry any) error {
	return nil
}
