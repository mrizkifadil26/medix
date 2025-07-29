package enricher

import (
	"fmt"

	"github.com/mrizkifadil26/medix/enricher/tmdb"
	"github.com/mrizkifadil26/medix/model"
)

// TODO: fix this
func Enrich(
	entries []model.MediaEntry,
	config *Config,
) (any, error) {
	client := tmdb.NewClient(config.APIKey)
	// scorer := NewScorer(config.Scoring)

	var enriched []any
	for _, entry := range entries {
		// normalize title
		// normalized := NormalizeTitle(entry.Title)

		results, err := client.Search("movie", tmdb.SearchQuery{
			Query:       entry.Name,
			PrimaryYear: string(entry.ContentType), // should be Year
		})
		if err != nil {
			fmt.Printf("⚠️  Failed to search TMDb for %q: %v\n", entry.Name, err)
			enriched = append(enriched, entry)
			continue
		}

		// best := scorer.BestMatch(entry, results)
		best := tmdb.PickBestMovieMatch(results, entry.Name, 2010) // should be derived from Year
		if best == nil {
			fmt.Printf("❌ No good match found for %q (%d)\n", entry.Name, 2010)
			enriched = append(enriched, entry)
			continue
		}

		// entry.TMDB = best.ToTMDBMeta()
		// enriched = append(enriched, entry)
		// fmt.Printf("✅ Enriched: %s → %s (%d)\n", entry.Name, best.Title(), best.Year())
	}

	return enriched, nil
}
