package scorer

var DefaultConfig = &ScoreConfig{
	TitleWeights: TitleWeights{
		ExactMatch:    100,
		OriginalMatch: 90,
		PartialMatch:  70,
		PartialOrig:   60,
	},
	YearWeights: YearWeights{
		Exact:   50,
		OffBy1:  30,
		OffBy2:  10,
		Penalty: -20,
	},
	GenreWeights: GenreWeights{
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
	PopularityWeight: 1,
	VoteWeight:       2,
}
