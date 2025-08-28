package tmdb

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/mrizkifadil26/medix/enricher/tmdb/scorer"
	"github.com/mrizkifadil26/medix/utils/cache"
	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

const (
	tmdbCache = "tmdb.cache.json"
	dataCache = "data.cache.json"
	apiKey    = "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI3N2QxNGJiN2JkODYyY2E0ZTE4MzBiZWNiODgxNGU3NyIsIm5iZiI6MTU2MzYzMjEyOC43NzUsInN1YiI6IjVkMzMyMjAwYWU2ZjA5MDAwZTdiNWJlZiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.ZpKTx-psLaWjwS-zRvpDLO7QKmoNJnF_xubbb8vn-48"
)

type QueryInput struct {
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

	// TMDb result (written to JSON under .enriched)
	Enriched EnrichedData
}

type TMDbEnricher struct {
	client    *Client
	tmdbCache *cache.Manager
	dataCache *cache.Manager

	genres map[int]string
	langs  map[string]string
}

func (t *TMDbEnricher) Name() string { return "tmdb" }

func (t *TMDbEnricher) Enrich(
	data any,
	options map[string]string,
) (any, error) {
	client := NewClient(apiKey)

	tmdbCM := cache.NewManager("tmdb.cache.json")
	if err := tmdbCM.Load(); err != nil {
		log.Fatalf("failed to load TMDb cache: %v", err)
	}

	dataCM := cache.NewManager("data.cache.json")
	if err := dataCM.Load(); err != nil {
		log.Printf("warning: failed to load data cache: %v", err)
	}

	e := &TMDbEnricher{
		client:    client,
		tmdbCache: tmdbCM,
		dataCache: dataCM,
	}

	// load genres for movies
	genreMap, err := LoadGenreMap(client, tmdbCM, "movie")
	if err != nil {
		log.Fatalf("failed to load genres: %v", err)
	}
	e.genres = genreMap

	// load languages
	langMap, err := LoadLanguageMap(client, tmdbCM)
	if err != nil {
		log.Fatalf("failed to load languages: %v", err)
	}
	e.langs = langMap

	// --- Build query inputs ---
	queries, err := extractQueries(data)
	if err != nil {
		return nil, err
	}

	// item_count (from JSON)
	itemCount, err := jsonpath.Get(data, "item_count")
	if err != nil {
		panic("item count not found: " + err.Error())
	}

	count, _ := itemCount.(float64)
	progress := &Progress{total: int32(count)}

	// --- Parallel enrichment ---
	c := options["concurrency"]
	concurrency, _ := strconv.Atoi(c)
	if concurrency <= 0 {
		concurrency = 5 // default
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)

	var mu sync.Mutex
	var jsonmu sync.RWMutex
	errorsMap := make(map[string]string)

	for idx, _ := range queries {
		wg.Add(1)

		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{} // acquire slot
			defer func() { <-sem }()

			// Get slug for caching
			slugNode, _ := jsonpath.Get(data, fmt.Sprintf("items.%d.slug", idx))
			slug, _ := slugNode.(string)

			var q QueryInput
			if slug != "" {
				if cached, ok := e.dataCache.Get("movie", slug); ok {
					q, _ = cached.(QueryInput)
				} else {
					q = queries[idx]
					e.dataCache.Put("movie", slug, q)
				}
			} else {
				q = queries[idx]
			}

			enriched := e.enrichItem(q, idx, slug)

			// Set JSON for enriched data
			jsonmu.Lock()
			jsonpath.Set(data, fmt.Sprintf("items.%d.enriched", idx), enriched.Enriched)
			jsonmu.Unlock()

			// Collect errors
			if enriched.Error != "" {
				mu.Lock()
				errorsMap[enriched.Title] = enriched.Error
				mu.Unlock()
			}

			display := enriched.Title
			if enriched.Year != "" {
				display = fmt.Sprintf("%s (%s)", enriched.Title, enriched.Year)
			}

			progress.Inc(display, enriched.Error)
		}(idx)
	}

	wg.Wait()

	jsonpath.Set(data, "errors", errorsMap)

	// Save caches
	if err := tmdbCM.Save(); err != nil {
		log.Printf("failed to save TMDb cache: %v", err)
	}

	if err := dataCM.Save(); err != nil {
		log.Printf("failed to save data cache: %v", err)
	}

	return data, nil
}

func (e *TMDbEnricher) enrichItem(
	query QueryInput,
	idx int,
	slug string,
) *EnrichedItem {
	title := query.Title
	year := query.Year
	alt := query.AlternateTitle

	item := &EnrichedItem{
		Index:          idx,
		Title:          title,
		Year:           year,
		AlternateTitle: alt,
	}

	if title == "" {
		item.Error = "missing title"
		return item
	}

	// Try cache first
	if slug != "" {
		if cached, ok := e.tmdbCache.Get("movie", slug); ok {
			if cachedItem, ok := cached.(EnrichedData); ok {
				item.Enriched = cachedItem
				return item
			}
		}
	}

	search := SearchQuery{Query: title}
	if year != "" {
		search.Year = year
	}

	// Search TMDb
	results, err := e.client.Search("movie", search)
	if err != nil {
		item.Error = fmt.Sprintf("tmdb search failed: %v", err)
		return item
	}

	// Fallback with alternate_title
	if len(results) == 0 && alt != "" {
		altQuery := SearchQuery{Query: alt}
		if year != "" {
			altQuery.Year = year
		}

		results, err = e.client.Search("movie", altQuery)
		if err != nil {
			item.Error = fmt.Sprintf("tmdb alt search failed: %v", err)
			return item
		}
	}

	if len(results) == 0 {
		if year == "" {
			item.Error = "missing year, no results found"
		} else {
			item.Error = "tmdb no results found"
		}

		return item
	}

	var yearInt int
	if year != "" {
		yearInt, _ = strconv.Atoi(year)
	}

	best := scorer.PickBestMatch(results, title, yearInt, e.genres)
	if best == nil {
		item.Error = "no good match found"
		return item
	}

	// Map genre IDs → names
	var genres []string
	for _, id := range best.GenreIDs {
		if g, ok := e.genres[id]; ok {
			genres = append(genres, g)
		}
	}

	langName := e.langs[best.OriginalLanguage]

	// Fill enriched data
	item.Enriched = EnrichedData{
		TMDBID:        best.ID,
		MatchedTitle:  best.Title,
		OriginalTitle: best.OriginalTitle,
		ReleaseDate:   best.ReleaseDate,
		Genres:        genres,
		Language:      langName,
		PosterPath:    best.PosterPath,
		Overview:      best.Overview,
	}

	return item
}

func extractQueries(data any) ([]QueryInput, error) {
	nodes, err := jsonpath.Get(data, "items.#.metadata")
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	arr, ok := nodes.([]any)
	if !ok {
		return nil, fmt.Errorf("expected array of metadata, got %T", nodes)
	}

	var queries []QueryInput
	for _, node := range arr {
		q := QueryInput{}
		if v, err := jsonpath.Get(node, "title"); err == nil {
			if s, ok := v.(string); ok {
				q.Title = s
			}
		}

		if v, err := jsonpath.Get(node, "year"); err == nil {
			if s, ok := v.(string); ok {
				q.Year = s
			}
		}

		if v, err := jsonpath.Get(node, "alternate_title"); err == nil {
			if s, ok := v.(string); ok {
				q.AlternateTitle = s
			}
		}

		// only append if title exists
		if q.Title != "" {
			queries = append(queries, q)
		}
	}

	return queries, nil
}
