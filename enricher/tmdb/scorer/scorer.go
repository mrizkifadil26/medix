package scorer

import (
	"strings"
)

/*
func scoreMedia(
	item MediaItem,
	expectedTitle string,
	expectedYear int,
	genreMap map[int]string,
) int {
	return scoreMediaWithConfig(item, expectedTitle, expectedYear, genreMap, DefaultConfig)
}
*/

func scoreMediaWithConfig(
	item MediaItem,
	expectedTitle string,
	expectedYear int,
	genreMap map[int]string,
	cfg *ScoreConfig,
) int {
	score := 0

	// Normalize title comparisons
	titleLower := strings.ToLower(item.GetTitle())
	originalLower := strings.ToLower(item.GetOriginalTitle())
	expectedLower := strings.ToLower(expectedTitle)

	// Title matching
	switch {
	case titleLower == expectedLower:
		score += cfg.TitleWeights.ExactMatch
	case originalLower == expectedLower:
		score += cfg.TitleWeights.OriginalMatch
	case strings.Contains(titleLower, expectedLower):
		score += cfg.TitleWeights.PartialMatch
	case strings.Contains(originalLower, expectedLower):
		score += cfg.TitleWeights.PartialOrig
	}

	// Year handling
	if item.GetReleaseYear() != 0 && expectedYear != 0 {
		diff := abs(item.GetReleaseYear() - expectedYear)
		switch diff {
		case 0:
			score += cfg.YearWeights.Exact
		case 1:
			score += cfg.YearWeights.OffBy1
		case 2:
			score += cfg.YearWeights.OffBy2
		default:
			score += cfg.YearWeights.Penalty * diff
		}
	}

	// Popularity & votes
	score += int(item.GetPopularity()) * cfg.PopularityWeight
	score += int(item.GetVoteAverage()) * cfg.VoteWeight

	// Genres
	for idx, id := range item.GetGenreIDs() {
		if g, ok := genreMap[id]; ok {
			if cfg.GenreWeights.MinorGenres[g] {
				score += cfg.GenreWeights.Minor
			} else if idx == 0 {
				score += cfg.GenreWeights.Primary
			} else {
				score += cfg.GenreWeights.Secondary
			}
		}
	}

	return score
}

func PickBestMatch[T MediaItem](
	items []T,
	expectedTitle string,
	expectedYear int,
	genreMap map[int]string,
) *T {
	return PickBestMatchWithConfig(items, expectedTitle, expectedYear, genreMap, DefaultConfig)
}

func PickBestMatchWithConfig[T MediaItem](
	items []T,
	expectedTitle string,
	expectedYear int,
	genreMap map[int]string,
	cfg *ScoreConfig,
) *T {
	var best *T
	highestScore := -1

	for _, item := range items {
		score := scoreMediaWithConfig(item, expectedTitle, expectedYear, genreMap, cfg)
		if score > highestScore {
			best = &item
			highestScore = score
		}
	}

	return best
}

func abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}
