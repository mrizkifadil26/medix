package tmdb

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/mrizkifadil26/medix/enricher/tmdb/scorer"
	"github.com/mrizkifadil26/medix/utils"
	"github.com/mrizkifadil26/medix/utils/datawrapper"
	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

const (
	cacheFile      = "cache.json"
	genreCacheFile = "genres.cache.json"
	langCacheFile  = "languages.cache.json"
)

// type Enricher struct {
// 	client   *tmdb.Client
// 	cache    map[string]bool
// 	genreMap map[int]string
// 	langMap  map[string]string
// 	errors   map[string]string

// 	mu     sync.Mutex   // protects cache and errors
// 	jsonMu sync.RWMutex // protects jsonpath.Get/Set
// }

type TMDbEnricher struct {
	client *Client

	genres map[int]string
	langs  map[string]string

	cache  map[string]bool
	errors map[string]string
	mu     sync.Mutex
	jsonMu sync.RWMutex
}

func (t *TMDbEnricher) Name() string { return "tmdb" }

func (t *TMDbEnricher) Enrich(
	data datawrapper.Data,
	params map[string]string,
) (any, error) {
	root := data.Raw()
	client := NewClient("eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI3N2QxNGJiN2JkODYyY2E0ZTE4MzBiZWNiODgxNGU3NyIsIm5iZiI6MTU2MzYzMjEyOC43NzUsInN1YiI6IjVkMzMyMjAwYWU2ZjA5MDAwZTdiNWJlZiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.ZpKTx-psLaWjwS-zRvpDLO7QKmoNJnF_xubbb8vn-48")

	e := &TMDbEnricher{
		client: client,
		cache:  loadCache(cacheFile),
		errors: make(map[string]string),
	}

	genreMap, err := loadGenreCache(client, "movie")
	if err != nil {
		return nil, fmt.Errorf("failed to load genres: %w", err)
	}
	e.genres = genreMap

	langMap, err := loadLangCache(client, langCacheFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load languages: %w", err)
	}
	e.langs = langMap

	e.jsonMu.RLock()
	itemsNode, ok := data.Get("items")
	if !ok {
		itemsNode = datawrapper.WrapData([]any{}) // fallback to empty array
	}
	e.jsonMu.RUnlock()

	// item_count (from JSON)
	itemCount := 0
	e.jsonMu.RLock()
	if vNode, ok := data.Get("item_count"); ok {
		// vNode is Data (could be ValueData)
		raw := vNode.Raw()
		if f, ok := raw.(float64); ok {
			itemCount = int(f)
		}
	}
	e.jsonMu.RUnlock()

	// fallback if item_count is missing or zero
	if itemCount == 0 && itemsNode.Type() == "array" {
		itemCount = len(itemsNode.Keys())
	}

	progress := &Progress{total: int32(itemCount)}

	// --- Parallel enrichment ---
	c := params["concurrency"]
	concurrency, _ := strconv.Atoi(c)
	if concurrency <= 0 {
		concurrency = 5 // default
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)

	for _, idx := range itemsNode.Keys() {
		itemNode, _ := itemsNode.Get(idx)
		wg.Add(1)

		go func(idx any, itemNode datawrapper.Data) {
			defer wg.Done()
			sem <- struct{}{} // acquire slot
			defer func() { <-sem }()

			title, err := e.enrichItem(root, itemNode, idx)
			if title != "" {
				progress.Inc(title, err)
			}
		}(idx, itemNode)
	}

	wg.Wait()
	// --- End parallel ---

	saveCache(cacheFile, e.cache)

	e.jsonMu.Lock()
	jsonpath.Set(root, "errors", e.errors)
	e.jsonMu.Unlock()

	return root, nil
}

type Progress struct {
	total   int32
	current int32
}

func (p *Progress) Inc(title string, errMessage error) {
	newVal := atomic.AddInt32(&p.current, 1)
	status := "✅"
	displayTitle := title
	if errMessage != nil && errMessage.Error() != "" {
		status = "❌"
		displayTitle = fmt.Sprintf("%s (%s)", title, errMessage.Error())
	}

	percent := float64(newVal) / float64(p.total) * 100

	fmt.Printf("[%d/%d %.1f%%] %s %s\n",
		newVal, p.total, percent, status, displayTitle)
}

// func Enrich(
// 	data datawrapper.Data,
// 	config *Config,
// ) (any, error) {

// }

func (e *TMDbEnricher) enrichItem(
	root any,
	itemNode datawrapper.Data,
	idx any,
) (string, error) {
	entry, ok := itemNode.(*datawrapper.OrderedMapData)
	if !ok {
		return "", fmt.Errorf("item at index %v is not an object", idx)
	}

	e.jsonMu.RLock()
	metadataNode, ok := entry.Get("metadata")
	e.jsonMu.RUnlock()
	if !ok {
		name := getString(entry, "name")
		e.setError(name, "missing metadata")
		return name, fmt.Errorf("missing metadata for %q", name)
	}

	title := getString(metadataNode, "title")
	if title == "" {
		name := getString(entry, "name")
		e.setError("", "missing title")
		return name, fmt.Errorf("missing title in metadata")
	}

	e.jsonMu.RLock()
	_, alreadyEnriched := entry.Get("enriched")
	e.jsonMu.RUnlock()
	if alreadyEnriched {
		return title, nil
	}

	year := getString(metadataNode, "year") // optional

	// Build query
	query := SearchQuery{Query: title}
	if year != "" {
		query.Year = year
	}

	// Search TMDb
	results, err := e.client.Search("movie", query)
	if err != nil {
		e.setError(title, "tmdb search failed: "+err.Error())
		return title, fmt.Errorf("tmdb search failed for %q: %w", title, err)
	}

	// Fallback with alternate_title
	if len(results) == 0 {
		alt := getString(metadataNode, "alternate_title")
		if alt != "" {
			altQuery := SearchQuery{Query: alt}
			if year != "" {
				altQuery.Year = year
			}

			results, err = e.client.Search("movie", altQuery)
			if err != nil {
				e.setError(title, "tmdb search failed (alt): "+err.Error())
				return title, fmt.Errorf("tmdb alt search failed for %q: %w", title, err)
			}
		}
	}

	if len(results) == 0 {
		if year == "" {
			e.setError(title, "missing year, no results found")
			return title, fmt.Errorf("missing year, no results found for %q", title)

		} else {
			e.setError(title, "tmdb no results found")
			return title, fmt.Errorf("tmdb no results found for %q", title)
		}
	}

	var yearInt int
	if year != "" {
		yearInt, _ = strconv.Atoi(year)
	}

	best := scorer.PickBestMatch(results, title, yearInt, e.genres)
	if best == nil {
		e.setError(title, "no good match found")
		return title, fmt.Errorf("no good match found for %q", title)
	}

	// Map genre IDs → names
	var genres []string
	for _, id := range best.GenreIDs {
		if g, ok := e.genres[id]; ok {
			genres = append(genres, g)
		}
	}

	langName := e.langs[best.OriginalLanguage]

	// Write results back
	e.jsonMu.Lock()
	jsonpath.Set(root, fmt.Sprintf("items.%d.enriched.title", idx), best.Title)
	jsonpath.Set(root, fmt.Sprintf("items.%d.enriched.original_title", idx), best.OriginalTitle)
	jsonpath.Set(root, fmt.Sprintf("items.%d.enriched.release_date", idx), best.ReleaseDate)
	jsonpath.Set(root, fmt.Sprintf("items.%d.enriched.genres", idx), genres)
	jsonpath.Set(root, fmt.Sprintf("items.%d.enriched.language", idx), langName)
	jsonpath.Set(root, fmt.Sprintf("items.%d.enriched.poster_path", idx), best.PosterPath)
	jsonpath.Set(root, fmt.Sprintf("items.%d.enriched.overview", idx), best.Overview)
	e.jsonMu.Unlock()

	e.setCache(title, true)

	return fmt.Sprintf("%v (%v)", title, year), nil
}

// loadLangCache loads the language list and returns a map[iso_code]name
func loadLangCache(client *Client, cachePath string) (map[string]string, error) {
	langMap := make(map[string]string)

	// try reading existing cache
	if _, err := os.Stat(cachePath); err == nil {
		if err := utils.LoadJSON(cachePath, &langMap); err == nil {
			return langMap, nil
		}
	}

	// fallback: fetch from TMDb API
	languages, err := client.GetLanguages()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch languages from TMDb: %w", err)
	}

	for _, l := range languages {
		code := l.ISO639_1
		name := l.EnglishName
		if name == "" {
			name = l.Name
			if name == "" {
				name = "Unknown"
			}
		}
		if code != "" {
			langMap[code] = name
		}
	}

	// save cache for future runs
	if err := utils.WriteJSON(cachePath, langMap); err != nil {
		return nil, fmt.Errorf("failed to save language cache: %w", err)
	}

	return langMap, nil
}

func loadGenreCache(client *Client, kind string) (map[int]string, error) {
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

func (e *TMDbEnricher) setError(key, msg string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.errors[key] = msg
}

func (e *TMDbEnricher) setCache(title string, val bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cache[title] = val
}

// --- helpers ---
func loadCache(path string) map[string]bool {
	cache := make(map[string]bool)
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &cache)
	}

	return cache
}

func saveCache(path string, cache map[string]bool) {
	if data, err := json.Marshal(cache); err == nil {
		_ = os.WriteFile(path, data, 0644)
	}
}

func getString(node any, key string) string {
	dataNode, ok := node.(datawrapper.Data)
	if !ok {
		return ""
	}

	valueNode, ok := dataNode.Get(key)
	if !ok || valueNode == nil {
		return ""
	}

	// unwrap raw value
	raw := valueNode.Raw()
	if s, ok := raw.(string); ok {
		return s
	}

	return ""
}
