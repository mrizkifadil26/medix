package enricher

type Config struct {
	InputFile   string `json:"input"`       // Path to raw media entries (scanned)
	OutputFile  string `json:"output"`      // Path to write enriched result
	APIKey      string `json:"tmdb_key"`    // TMDb API key
	WeightsFile string `json:"weights"`     // Optional: scoring weights JSON
	OnlyKind    string `json:"only_kind"`   // "movie" or "tv" (optional filter)
	Concurrency int    `json:"concurrency"` // Optional: default 1
}
