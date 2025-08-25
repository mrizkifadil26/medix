package enricher

type Config struct {
	Root   string `json:"root"`   // Path to raw media entries (scanned)
	Output string `json:"output"` // Path to write enriched result
}
