package tmdb

type QueryInput struct {
	Slug           string
	Title          string
	Year           string
	AlternateTitle string
}

type EnrichedData struct {
	TMDBID        int      `json:"tmdb_id,omitempty"`
	MatchedTitle  string   `json:"matched_title,omitempty"`
	OriginalTitle string   `json:"original_title,omitempty"`
	ReleaseDate   string   `json:"release_date,omitempty"`
	Genres        []string `json:"genres,omitempty"`
	Language      string   `json:"language,omitempty"`
	PosterPath    string   `json:"poster_path,omitempty"`
	Overview      string   `json:"overview,omitempty"`
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
	Enriched EnrichedData
}
