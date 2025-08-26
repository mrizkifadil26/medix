package enricher

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mrizkifadil26/medix/enricher/tmdb"
	"github.com/mrizkifadil26/medix/enricher/tmdb/scorer"
	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

func Enrich(
	data any,
	config *Config,
) (any, error) {
	var cache map[string]bool
	cacheFile := "cache.json"

	// initialize
	cache = make(map[string]bool)

	if cacheData, err := os.ReadFile(cacheFile); err == nil {
		json.Unmarshal(cacheData, &cache)
	}

	// client := tmdb.NewClient(config.APIKey)
	client := tmdb.NewClient("eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI3N2QxNGJiN2JkODYyY2E0ZTE4MzBiZWNiODgxNGU3NyIsIm5iZiI6MTU2MzYzMjEyOC43NzUsInN1YiI6IjVkMzMyMjAwYWU2ZjA5MDAwZTdiNWJlZiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.ZpKTx-psLaWjwS-zRvpDLO7QKmoNJnF_xubbb8vn-48")
	genreMap, err := loadGenreCache(client, "movie")
	if err != nil {
		return nil, fmt.Errorf("failed to load genres: %w", err)
	}

	// var enriched []any
	dt := data.(map[string]any)
	items := dt["items"].([]any)

	errors := make(map[string]string)
	for idx, item := range items {
		var errs []string

		entry := item.(map[string]any)
		metadata, ok := entry["metadata"].(map[string]any)
		if !ok {
			name := entry["name"].(string)
			errors[name] = "missing metadata"

			fmt.Println("No metadata found for entry:", entry)
			continue
		}

		title := metadata["title"].(string)
		// fmt.Printf("title: %v | cached: %v | enriched: %v\n", title, cache[title], entry["enriched"] != nil)

		// if cached, ok := cache[title]; ok && cached && entry["enriched"] != nil {
		// 	fmt.Println("Skipping cached and already enriched entry:", title)
		// 	continue
		// }

		if entry["enriched"] != nil {
			// fmt.Println("Skipping cached and already enriched entry:", title)
			continue
		}

		year, _ := metadata["year"].(string)
		if year == "" {
			errs = append(errs, "missing year")

			fmt.Println("Year not found in metadata for:", title)
			year = ""
		}

		if len(errs) > 0 {
			errors[title] = strings.Join(errs, ", ")
			continue
		}

		results, err := client.Search("movie", tmdb.SearchQuery{
			Query: title, // should be Name
			Year:  year,  // should be Year
		})

		if err != nil {
			errors[title] = "tmdb search failed: " + err.Error()
			fmt.Printf("Failed to search TMDb for %q: %v\n", title, err)
			continue
		}

		// fallback using alternate_title
		if len(results) == 0 && metadata["alternate_title"] != "" {
			if altVal, ok := metadata["alternate_title"]; ok {
				if altStr, ok := altVal.(string); ok && altStr != "" {
					alternateTitle, _ := metadata["alternate_title"].(string)
					fmt.Printf("No results for %q, trying alternate title %q\n", title, alternateTitle)

					results, err = client.Search("movie", tmdb.SearchQuery{
						Query: alternateTitle,
						Year:  year,
					})

					if err != nil {
						errors[title] = "tmdb search failed on alternate title: " + err.Error()
						fmt.Printf("Failed to search TMDb for alternate title %q: %v\n", alternateTitle, err)
						continue
					}

				}
			}
		}

		if len(results) == 0 {
			errors[title] = "tmdb no results found"
			fmt.Printf("No result found for %q\n", title)
			continue
		}

		yearInt, _ := strconv.Atoi(year)
		bestResult := scorer.PickBestMatch[tmdb.SearchItem](
			results, title, yearInt, genreMap) // should be derived from Year

		if bestResult == nil {
			errors[title] = "no good match found"
			fmt.Printf("No good match found for %q\n", title)
			continue
		}

		var genres []string
		for _, id := range bestResult.GenreIDs {
			if g, ok := genreMap[id]; ok {
				genres = append(genres, g)
			}
		}

		if len(results) > 0 {
			targetTitle := strings.ReplaceAll("items.#.enriched.title", "#", strconv.Itoa(idx))
			targetReleaseDate := strings.ReplaceAll("items.#.enriched.release_date", "#", strconv.Itoa(idx))
			targetGenres := strings.ReplaceAll("items.#.enriched.genres", "#", strconv.Itoa(idx))

			jsonpath.Set(data, targetTitle, bestResult.Title)
			jsonpath.Set(data, targetReleaseDate, bestResult.ReleaseDate)
			jsonpath.Set(data, targetGenres, genres)
		} else {
			fmt.Println("No results found for:", title)
		}

		cache[title] = true
	}

	cacheData, _ := json.Marshal(cache)
	os.WriteFile(cacheFile, cacheData, 0644)

	jsonpath.Set(data, "errors", errors)

	return data, nil
}

const genreCacheFile = "genres.cache.json"

func loadGenreCache(client *tmdb.Client, kind string) (map[int]string, error) {
	const genreCacheFile = "genres.cache.json"

	// Try from file
	if data, err := os.ReadFile(genreCacheFile); err == nil {
		var genres map[string]map[int]string
		if err := json.Unmarshal(data, &genres); err == nil {
			if g, ok := genres[kind]; ok {
				return g, nil
			}
		}
	}

	// If not found, fetch from TMDb
	resultGenres, err := client.GetGenres(kind)
	if err != nil {
		return nil, err
	}

	genreMap := make(map[int]string)
	for _, g := range resultGenres {
		genreMap[g.ID] = g.Name
	}

	// Save to file
	var genres map[string]map[int]string
	if data, err := os.ReadFile(genreCacheFile); err == nil {
		json.Unmarshal(data, &genres)
	}
	if genres == nil {
		genres = make(map[string]map[int]string)
	}
	genres[kind] = genreMap

	enc, _ := json.MarshalIndent(genres, "", "  ")
	_ = os.WriteFile(genreCacheFile, enc, 0644)

	return genreMap, nil
}
