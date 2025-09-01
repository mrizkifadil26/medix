package tmdb

type QueryInput struct {
	Slug           string
	Title          string
	Year           string
	AlternateTitle string
}

type EnrichedItem struct {
	Slug           string
	Index          int
	Title          string
	Year           string
	AlternateTitle string
	Error          string
	Source         string // "cache" or "remote"

	// TMDb result (written to JSON under .enriched)
	Enriched *EnrichedData
}

type EnrichedData struct {
	TMDbID        int      `json:"tmdb_id,omitempty"`
	Title         string   `json:"title,omitempty"`
	OriginalTitle string   `json:"original_title,omitempty"`
	ReleaseDate   string   `json:"release_date,omitempty"`
	Genres        []string `json:"genres,omitempty"`
	Language      string   `json:"language,omitempty"`
	PosterPath    string   `json:"poster_path,omitempty"`
	Overview      string   `json:"overview,omitempty"`

	Credits         *TMDbCredits `json:"credits,omitempty"`
	Recommendations []TMDbMain   `json:"recommendations,omitempty"`
	Similar         []TMDbMain   `json:"similar,omitempty"`
}

type TMDbMain struct {
	TMDbID        int      `json:"tmdb_id,omitempty"`
	Title         string   `json:"title,omitempty"`
	OriginalTitle string   `json:"original_title,omitempty"`
	ReleaseDate   string   `json:"release_date,omitempty"`
	Genres        []string `json:"genres,omitempty"`
	Language      string   `json:"language,omitempty"`
	PosterPath    string   `json:"poster_path,omitempty"`
	Overview      string   `json:"overview,omitempty"`
}

type TMDbCredits struct {
	Actors    []Person `json:"actors,omitempty"`
	Directors []Person `json:"directors,omitempty"`
	Producers []Person `json:"producers,omitempty"`
}

type Person struct {
	ID   int    `json:"id"`   // TMDb person ID (or internal ID if local)
	Name string `json:"name"` // human-readable
	Role string `json:"role"` // optional: "actor", "director", "producer", etc.
}
