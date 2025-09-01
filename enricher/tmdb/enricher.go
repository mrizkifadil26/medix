package tmdb

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/mrizkifadil26/medix/enricher/tmdb/scorer"
	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type TMDbEnricher struct {
	client    *Client
	dataCache *dataCache

	genreService  *GenreService
	langService   *LanguageService
	creditService *CreditService
	config        *Config
}

func NewTMDbEnricher(cfg *Config) *TMDbEnricher {
	client := NewClient(cfg.APIKey)
	dataCache, genreCache, langCache := newCaches()
	creditsCache := newCreditsCache()

	// Load caches at startup (best effort)
	_ = dataCache.Load()
	_ = genreCache.Load()
	_ = langCache.Load()
	_ = creditsCache.Load()

	return &TMDbEnricher{
		client:        client,
		dataCache:     dataCache,
		genreService:  NewGenreService(client, genreCache),
		langService:   NewLanguageService(client, langCache),
		creditService: NewCreditService(client, creditsCache),
		config:        cfg,
	}
}

func (t *TMDbEnricher) Name() string { return "tmdb" }

func (t *TMDbEnricher) Enrich(
	data any,
	options map[string]string,
) (any, error) {
	queries, err := extractQueries(data)
	if err != nil {
		return nil, err
	}

	itemCount, err := jsonpath.Get(data, "item_count")
	if err != nil {
		return nil, fmt.Errorf("item count not found: %w", err)
	}

	total := int32(asInt(itemCount))
	progress := &Progress{total: total}

	concurrency := parseConcurrency(options, 5)

	var (
		wg       sync.WaitGroup
		sem      = make(chan struct{}, concurrency)
		errorsMu sync.Mutex
		jsonMu   sync.RWMutex
		errors   = make(map[string]string)
	)

	for idx := range queries {
		wg.Add(1)

		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{} // acquire slot
			defer func() { <-sem }()

			result := t.enrichItem(queries[idx], idx)

			// Set JSON for enriched data
			// Only set enriched if not nil
			if result.Enriched != nil {
				jsonMu.Lock()
				// _ = jsonpath.Set(data, fmt.Sprintf("items.%d.enriched", idx), result.Enriched)
				_ = jsonpath.Set(data, fmt.Sprintf("items.%d.enriched", idx), result.Enriched)
				jsonMu.Unlock()
			}

			// Collect errors
			if result.Error != "" {
				errorsMu.Lock()
				errors[result.Title] = result.Error
				errorsMu.Unlock()
			}

			display := result.Title
			if result.Year != "" {
				display = fmt.Sprintf("%s (%s)", result.Title, result.Year)
			}

			progress.Inc(display, result.Error, result.Source)
		}(idx)
	}

	wg.Wait()

	_ = jsonpath.Set(data, "errors", errors)

	// Save data cache (best-effort)
	_ = t.dataCache.Save()

	return data, nil
}

func (e *TMDbEnricher) enrichItem(
	query QueryInput,
	idx int,
) *EnrichedItem {
	title := query.Title
	year := query.Year
	alt := query.AlternateTitle
	slug := query.Slug

	item := &EnrichedItem{
		Slug:           slug,
		Index:          idx,
		Title:          title,
		Year:           year,
		AlternateTitle: alt,
	}

	if title == "" {
		item.Error = "missing title"
		return item
	}

	// --- Try cache first ---
	if cached, ok := e.dataCache.Get("movie", slug); ok && cached != nil {
		cached.Source = "cache"
		return cached
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

	// Step 1: get genre map for scoring
	genreMap, err := e.genreService.Get("movie")
	if err != nil {
		item.Error = fmt.Sprintf("failed to fetch genre map: %v", err)
		return item
	}

	best := scorer.PickBestMatch(results, title, yearInt, genreMap)
	if best == nil {
		item.Error = "no good match found"
		return item
	}

	// Map genre IDs → names
	genres, err := e.genreService.Resolve("movie", best.GenreIDs)
	if err != nil {
		item.Error = fmt.Sprintf("failed to resolve genres: %v", err)
		return item
	}

	langName, err := e.langService.Resolve(best.OriginalLanguage)
	if err != nil {
		item.Error = fmt.Sprintf("failed to resolve language: %v", err)
		return item
	}

	// Fill enriched data
	item.Enriched = &EnrichedData{
		TMDbID:        best.ID,
		Title:         best.Title,
		OriginalTitle: best.OriginalTitle,
		ReleaseDate:   best.ReleaseDate,
		Genres:        genres,
		Language:      langName,
		PosterPath:    best.PosterPath,
		Overview:      best.Overview,
	}

	if e.config.FetchCredits {
		if credits, err := e.creditService.FetchCredits(best.ID); err == nil {
			item.Enriched.Credits = &credits
		}
	}

	item.Source = "remote" // mark it came from remote fetch
	if slug != "" {
		e.dataCache.Put("movie", slug, item)
	}

	return item
}

func extractQueries(data any) ([]QueryInput, error) {
	items, err := jsonpath.Get(data, "items.#")
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	arr, ok := items.([]any)
	if !ok {
		return nil, fmt.Errorf("expected array of items, got %T", items)
	}

	var queries []QueryInput
	for i, item := range arr {
		q := QueryInput{}

		// slug (error if missing)
		if v, err := jsonpath.Get(item, "slug"); err == nil {
			if s, ok := v.(string); ok && s != "" {
				q.Slug = s
			}
		}
		if q.Slug == "" {
			return nil, fmt.Errorf("missing slug for item at index %d", i)
		}

		// title
		if v, err := jsonpath.Get(item, "metadata.title"); err == nil {
			if s, ok := v.(string); ok {
				q.Title = s
			}
		}

		// year
		if v, err := jsonpath.Get(item, "metadata.year"); err == nil {
			if s, ok := v.(string); ok {
				q.Year = s
			}
		}

		// alternate title
		if v, err := jsonpath.Get(item, "metadata.alternate_title"); err == nil {
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

func asInt(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	default:
		return 0
	}
}

func parseConcurrency(options map[string]string, fallback int) int {
	if c, ok := options["concurrency"]; ok {
		if v, _ := strconv.Atoi(c); v > 0 {
			return v
		}
	}

	return fallback
}
