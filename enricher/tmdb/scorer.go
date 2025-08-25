package tmdb

import (
	"strconv"
	"strings"
)

type Scorer struct {
	Weights map[string]float64
}

func NewScorer(weights map[string]float64) *Scorer {
	return &Scorer{Weights: weights}
}

func scoreMovie(item SearchItem, expectedTitle string, expectedYear int) int {
	score := 0

	// Normalize title comparisons
	titleLower := strings.ToLower(item.Title)
	originalLower := strings.ToLower(item.OriginalTitle)
	expectedLower := strings.ToLower(expectedTitle)

	// Title matching
	if titleLower == expectedLower {
		score += 100 // exact match
	} else if originalLower == expectedLower {
		score += 90 // original title matches
	} else if strings.Contains(titleLower, expectedLower) {
		score += 70
	} else if strings.Contains(originalLower, expectedLower) {
		score += 60
	}

	// Year handling
	if item.ReleaseDate != "" && expectedYear != 0 {
		yearStr := strings.Split(item.ReleaseDate, "-")[0]
		if y, err := strconv.Atoi(yearStr); err == nil {
			diff := abs(y - expectedYear)
			switch diff {
			case 0:
				score += 50
			case 1:
				score += 30
			case 2:
				score += 10
			default:
				score -= 20 * diff // too far away = penalty
			}
		}
	}

	// Bonus: popularity and vote average as tie-breakers
	score += int(item.Popularity)      // 0-100 scale
	score += int(item.VoteAverage * 2) // 0-20 max

	// log.Printf("  📊 Final Score (%s): %d", item.OriginalTitle, score)

	return score
}

func PickBestMovieMatch(items []SearchItem, expectedTitle string, expectedYear int) *SearchItem {
	var best *SearchItem
	highestScore := -1

	for _, item := range items {
		score := scoreMovie(item, expectedTitle, expectedYear)
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
