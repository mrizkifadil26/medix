package enricher

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/mrizkifadil26/medix/enricher/tmdb"
	"github.com/mrizkifadil26/medix/enricher/tmdb/scorer"
	"github.com/mrizkifadil26/medix/utils"
	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

const (
	cacheFile      = "cache.json"
	genreCacheFile = "genres.cache.json"
)

type Enricher struct {
	client   *tmdb.Client
	cache    map[string]bool
	genreMap map[int]string
	errors   map[string]string

	mu sync.Mutex // protects cache and errors
}

type Progress struct {
	total   int32
	current int32
}

func (p *Progress) Inc(title string, isError bool) {
	newVal := atomic.AddInt32(&p.current, 1)
	status := "✅"
	if isError {
		status = "❌"
	}

	percent := float64(newVal) / float64(p.total) * 100

	fmt.Printf("[%d/%d %.1f%%] %s %s\n",
		newVal, p.total, percent, status, title)
}

func Enrich(
	data any,
	config *Config,
) (any, error) {
	client := tmdb.NewClient("eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI3N2QxNGJiN2JkODYyY2E0ZTE4MzBiZWNiODgxNGU3NyIsIm5iZiI6MTU2MzYzMjEyOC43NzUsInN1YiI6IjVkMzMyMjAwYWU2ZjA5MDAwZTdiNWJlZiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.ZpKTx-psLaWjwS-zRvpDLO7QKmoNJnF_xubbb8vn-48")

	e := &Enricher{
		client: client,
		cache:  loadCache(cacheFile),
		errors: make(map[string]string),
	}

	genreMap, err := loadGenreCache(client, "movie")
	if err != nil {
		return nil, fmt.Errorf("failed to load genres: %w", err)
	}
	e.genreMap = genreMap

	// root, ok := data.(map[string]any)
	// if !ok {
	// 	return nil, fmt.Errorf("invalid input: must be map[string]any")
	// }
	// items, _ := root["items"].([]any)
	root, err := AsObject(data)
	if err != nil {
		return nil, err
	}

	itemsVal, _ := root.Get("items")
	items, _ := itemsVal.([]any)

	// item_count (from JSON)
	itemCount := 0
	if v, ok := root.Get("item_count"); ok {
		if f, ok := v.(float64); ok {
			itemCount = int(f)
		}
	}

	if itemCount == 0 {
		itemCount = len(items)
	}

	progress := &Progress{total: int32(itemCount)}

	// --- Parallel enrichment ---
	concurrency := config.Concurrency
	if concurrency <= 0 {
		concurrency = 5 // default
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)

	for idx, item := range items {
		wg.Add(1)

		go func(idx int, item any) {
			defer wg.Done()
			sem <- struct{}{} // acquire slot
			defer func() { <-sem }()

			title, isError := e.enrichItem(root, item, idx)
			if title != "" {
				progress.Inc(title, isError)
			}
		}(idx, item)
	}

	wg.Wait()
	// --- End parallel ---

	saveCache(cacheFile, e.cache)
	jsonpath.Set(data, "errors", e.errors)
	return data, nil
}

func (e *Enricher) enrichItem(root any, item any, idx int) (string, bool) {
	entry, err := AsObject(item)
	if err != nil {
		return "", true
	}

	metadata, ok := entry.Get("metadata")
	if !ok {
		name := getString(entry, "name")
		e.setError(name, "missing metadata")
		return name, true
	}

	title := getString(metadata, "title")
	if _, ok := entry.Get("enriched"); !ok {
		return title, false
	}

	year := getString(metadata, "year") // optional

	// Build query
	query := tmdb.SearchQuery{Query: title}
	if year != "" {
		query.Year = year
	}

	// Search TMDb
	results, err := e.client.Search("movie", query)
	if err != nil {
		e.setError(title, "tmdb search failed: "+err.Error())
		return title, true
	}

	// Fallback with alternate_title
	if len(results) == 0 {
		alt := getString(metadata, "alternate_title")
		if alt != "" {
			altQuery := tmdb.SearchQuery{Query: alt}
			if year != "" {
				altQuery.Year = year
			}

			results, err = e.client.Search("movie", altQuery)
			if err != nil {
				e.setError(title, "tmdb search failed (alt): "+err.Error())
				return title, true
			}
		}
	}

	if len(results) == 0 {
		if year == "" {
			e.setError(title, "missing year, no results found")
		} else {
			e.setError(title, "tmdb no results found")
		}

		return title, true
	}

	var yearInt int
	if year != "" {
		yearInt, _ = strconv.Atoi(year)
	}

	best := scorer.PickBestMatch(results, title, yearInt, e.genreMap)
	if best == nil {
		e.setError(title, "no good match found")
		return title, true
	}

	// Map genre IDs → names
	var genres []string
	for _, id := range best.GenreIDs {
		if g, ok := e.genreMap[id]; ok {
			genres = append(genres, g)
		}
	}

	// Write results back
	jsonpath.Set(root, fmt.Sprintf("items.%d.enriched.title", idx), best.Title)
	jsonpath.Set(root, fmt.Sprintf("items.%d.enriched.release_date", idx), best.ReleaseDate)
	jsonpath.Set(root, fmt.Sprintf("items.%d.enriched.genres", idx), genres)

	e.setCache(title, true)

	return title, false
}

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

func (e *Enricher) setError(key, msg string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.errors[key] = msg
}

func (e *Enricher) setCache(title string, val bool) {
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

func getString(m any, key string) string {
	obj, err := AsObject(m)
	if err != nil {
		return ""
	}

	if v, ok := obj.Get(key); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// func getField(m any, key string) (any, error) {
// 	switch mm := m.(type) {
// 	case map[string]any:
// 		val, ok := mm[key]
// 		if !ok {
// 			return nil, fmt.Errorf("key %q not found", key)
// 		}
// 		return val, nil

// 	case *orderedmap.OrderedMap:
// 		val, ok := mm.Get(key)
// 		if !ok {
// 			return nil, fmt.Errorf("key %q not found", key)
// 		}
// 		return val, nil

// 	default:
// 		return nil, fmt.Errorf("unsupported map type %T, must be map[string]any or *orderedmap.OrderedMap", m)
// 	}
// }

type Object interface {
	Get(key string) (any, bool)
	Set(key string, val any)
	Keys() []string
	Raw() any
}

// --- map[string]any wrapper ---
type MapWrapper map[string]any

func (m MapWrapper) Get(key string) (any, bool) { v, ok := m[key]; return v, ok }
func (m MapWrapper) Set(key string, val any)    { m[key] = val }
func (m MapWrapper) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
func (m MapWrapper) Raw() any { return m }

// --- orderedmap wrapper ---
type OrderedMapWrapper struct{ *utils.OrderedMap[string, any] }

func (o OrderedMapWrapper) Get(key string) (any, bool) { return o.OrderedMap.Get(key) }
func (o OrderedMapWrapper) Set(key string, val any)    { o.OrderedMap.Set(key, val) }
func (o OrderedMapWrapper) Keys() []string             { return o.OrderedMap.Keys() }
func (o OrderedMapWrapper) Raw() any                   { return o.OrderedMap }

// --- normalize function ---
func AsObject(v any) (Object, error) {
	switch root := v.(type) {
	case Object: // already wrapped
		return root, nil
	case map[string]any:
		return MapWrapper(root), nil
	case *utils.OrderedMap[string, any]:
		return OrderedMapWrapper{root}, nil
	default:
		return nil, fmt.Errorf("invalid input type %T, must be map[string]any or *OrderedMap[string, any]", v)
	}
}
