package tmdb

type SearchResult struct {
	Results []SearchItem `json:"results"`
}

type SearchItem struct {
	ID               int     `json:"id"`                // Unique TMDb ID
	Title            string  `json:"title"`             // From movie.title or tv.name
	OriginalTitle    string  `json:"original_title"`    // For fallback or search matching
	Overview         string  `json:"overview"`          // Summary
	PosterPath       string  `json:"poster_path"`       // For image thumbnail
	BackdropPath     string  `json:"backdrop_path"`     // For UI background
	ReleaseDate      string  `json:"release_date"`      // movie: release_date / tv: first_air_date
	GenreIDs         []int   `json:"genre_ids"`         // To map to genre names locally
	VoteAverage      float64 `json:"vote_average"`      // For scoring
	VoteCount        int     `json:"vote_count"`        // For credibility
	OriginalLanguage string  `json:"original_language"` // Optional filter
	Popularity       float64 `json:"popularity"`
}
