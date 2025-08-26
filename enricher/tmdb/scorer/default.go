package scorer

var DefaultConfig = &ScoreConfig{
	TitleWeights: struct {
		ExactMatch    int
		OriginalMatch int
		PartialMatch  int
		PartialOrig   int
	}{
		ExactMatch:    100,
		OriginalMatch: 90,
		PartialMatch:  70,
		PartialOrig:   60,
	},
	YearWeights: struct {
		Exact   int
		OffBy1  int
		OffBy2  int
		Penalty int
	}{
		Exact:   50,
		OffBy1:  30,
		OffBy2:  10,
		Penalty: -20,
	},
	PopularityWeight: 1,
	VoteWeight:       2,
	GenreWeights: struct {
		Primary     int
		Secondary   int
		Minor       int
		MinorGenres map[string]bool
	}{
		Primary:   20,
		Secondary: 10,
		Minor:     2,
		MinorGenres: map[string]bool{
			"Music":       true,
			"Documentary": true,
			"Short":       true,
			"TV Movie":    true,
		},
	},
}
