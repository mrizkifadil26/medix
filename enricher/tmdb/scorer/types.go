package scorer

type TitleWeights struct {
	ExactMatch    int
	OriginalMatch int
	PartialMatch  int
	PartialOrig   int
}

type YearWeights struct {
	Exact   int
	OffBy1  int
	OffBy2  int
	Penalty int
}

type GenreWeights struct {
	Primary     int
	Secondary   int
	Minor       int
	MinorGenres map[string]bool
}

type ScoreConfig struct {
	TitleWeights     TitleWeights
	YearWeights      YearWeights
	GenreWeights     GenreWeights
	PopularityWeight int
	VoteWeight       int
}

type MediaItem interface {
	GetTitle() string
	GetOriginalTitle() string
	GetReleaseYear() int // for TV: first_air_date, for movie: release_date
	GetPopularity() float64
	GetVoteAverage() float64
	GetGenreIDs() []int
}
