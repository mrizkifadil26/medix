package enricher

type EnricherConfig struct {
	Name    string      `json:"name"`
	Config  interface{} `json:"config,omitempty"` // per-enricher config
	Filters []string    `json:"filters,omitempty"`
}

type Options struct {
	Concurrency int              `json:"concurrency"`
	Enrichers   []EnricherConfig `json:"enrichers"`
}

type Config struct {
	Root    string  `json:"root"`   // Path to raw media entries (scanned)
	Output  string  `json:"output"` // Path to write enriched result
	Options Options `json:"options"`
}
